package basic

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"encoding/base64"
	"encoding/json"

	mongodb "github.com/meshify-app/meshify/mongo"
	log "github.com/sirupsen/logrus"

	"github.com/coreos/go-oidc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/meshify-app/meshify/core"
	model "github.com/meshify-app/meshify/model"
	"github.com/meshify-app/meshify/shadow"
	util "github.com/meshify-app/meshify/util"
	"golang.org/x/oauth2"
)

type Oauth2Basic struct{}

// Create an OAuth2 shim for BasicAuth
func (o *Oauth2Basic) Setup() error {
	return nil
}

func (o *Oauth2Basic) Logout() error {
	return nil
}

// CodeUrl get url to redirect client for auth
func (o *Oauth2Basic) CodeUrl(state string) string {
	return "_magic_string_fake_auth_no_redirect_"
}

// Exchange exchange code for Oauth2 token
func (o *Oauth2Basic) Exchange(code string) (*oauth2.Token, error) {

	// code contains the username and password base64 encoded
	// base64 decode the string
	userpass, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		return nil, err
	}

	// split the username and password
	parts := strings.SplitN(string(userpass), ":", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid username and password")
	}
	user := parts[0]
	pass := parts[1]

	// validate the username and password
	err = shadow.ShadowAuthPlain(user, pass)
	if err != nil {
		return nil, err
	}

	rand, err := util.GenerateRandomString(32)
	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{
		AccessToken:  rand,
		TokenType:    "Bearer",
		RefreshToken: "",
		Expiry:       time.Now().Add(time.Hour * 24),
	}
	// add the user to the token
	idtoken := &oidc.IDToken{Subject: user, Issuer: "Basic", IssuedAt: time.Now(), Expiry: time.Now().Add(time.Hour * 24)}

	m := make(map[string]interface{})

	raw, err := json.Marshal(idtoken)
	if err != nil {
		return nil, err
	}
	m["id_token"] = string(raw)
	token = token.WithExtra(m)

	return token, nil
}

// UserInfo get token user
func (o *Oauth2Basic) UserInfo(oauth2Token *oauth2.Token) (*model.User, error) {
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("no id_token field in oauth2 token")
	}

	verified := false
	var idToken *oidc.IDToken
	var err error

	// decode the json string into an idToken
	err = json.Unmarshal([]byte(rawIDToken), &idToken)
	if err != nil {
		return nil, err
	}

	verified = true

	if !verified || err != nil {
		return nil, err
	}

	// get some infos about user
	user := &model.User{}
	user.Sub = idToken.Subject
	user.Email = idToken.Subject + "@localhost"
	user.Email = strings.ToLower(user.Email)
	user.Picture = "/meshify-bw.png"
	user.Issuer = idToken.Issuer
	user.IssuedAt = idToken.IssuedAt
	log.Infof("user %v", user)

	// save to mongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_CONNECTION_STRING")))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Error(err)
		}
	}()

	collection := client.Database("meshify").Collection("users")

	data, err := json.Marshal(user)
	//	json := fmt.Sprintf("%v", user)
	var b interface{}
	err = bson.UnmarshalExtJSON([]byte(data), true, &b)

	findstr := fmt.Sprintf("{\"email\":\"%s\"}", user.Email)
	var filter interface{}
	err = bson.UnmarshalExtJSON([]byte(findstr), true, &filter)

	update := bson.M{
		"$set": b,
	}

	opts := options.Update().SetUpsert(true)
	res, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Error(err)
	}

	accounts, err := mongodb.ReadAllAccounts(user.Email)
	if err != nil {
		log.Error(err)
	} else {
		//  If there's no error and no account, create one.
		if len(accounts) == 0 {
			var account model.Account
			account.Email = user.Email
			account.Role = "Owner"
			account.Status = "Active"
			a, err := core.CreateAccount(&account)
			log.Infof("CREATE ACCOUNT = %v", a)
			if err != nil {
				log.Error(err)
			}
		}
	}

	//res, err := collection.InsertOne(ctx, b)

	log.Infof("Res: %v", res)
	return user, nil
}
