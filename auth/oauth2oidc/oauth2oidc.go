package oauth2oidc

import (
	"context"
	"fmt"
	"time"

	"encoding/json"

	"github.com/coreos/go-oidc"
	"github.com/meshify-app/meshify/core"
	model "github.com/meshify-app/meshify/model"
	mongodb "github.com/meshify-app/meshify/mongo"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"
	"gopkg.in/auth0.v4/management"

	//	"gopkg.in/auth0.v4"
	"os"
)

// Oauth2idc in order to implement interface, struct is required
type Oauth2idc struct{}

var (
	oauth2Config        *oauth2.Config
	oidcProvider        *oidc.Provider
	oidcIDTokenVerifier *oidc.IDTokenVerifier
)

// Setup validate provider
func (o *Oauth2idc) Setup() error {
	var err error

	oidcProvider, err = oidc.NewProvider(context.TODO(), os.Getenv("OAUTH2_PROVIDER"))
	if err != nil {
		return err
	}

	oidcIDTokenVerifier = oidcProvider.Verifier(&oidc.Config{
		ClientID: os.Getenv("OAUTH2_CLIENT_ID"),
	})

	oauth2Config = &oauth2.Config{
		ClientID:     os.Getenv("OAUTH2_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH2_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("OAUTH2_REDIRECT_URL"),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		Endpoint:     oidcProvider.Endpoint(),
	}

	return nil
}

// CodeUrl get url to redirect client for auth
func (o *Oauth2idc) CodeUrl(state string) string {
	return oauth2Config.AuthCodeURL(state)
}

// Exchange exchange code for Oauth2 token
func (o *Oauth2idc) Exchange(code string) (*oauth2.Token, error) {
	oauth2Token, err := oauth2Config.Exchange(context.TODO(), code)
	if err != nil {
		return nil, err
	}

	return oauth2Token, nil
}

// UserInfo get token user
func (o *Oauth2idc) UserInfo(oauth2Token *oauth2.Token) (*model.User, error) {
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("no id_token field in oauth2 token")
	}

	iDToken, err := oidcIDTokenVerifier.Verify(context.TODO(), rawIDToken)
	if err != nil {
		return nil, err
	}

	userInfo, err := oidcProvider.UserInfo(context.TODO(), oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		return nil, err
	}

	// ID Token payload is just JSON
	var claims map[string]interface{}
	if err := userInfo.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to get id token claims: %s", err)
	}

	// get some infos about user
	user := &model.User{}
	user.Sub = userInfo.Subject
	user.Email = userInfo.Email
	user.Profile = userInfo.Profile

	//	for k, v :=  range claims {
	//		user.Claims = user.Claims + "<br>" + k + ":" + fmt.Sprintf("%v", v)
	//	}

	log.Infof("user.Sub: %s", user.Sub)

	if v, found := claims["name"]; found && v != nil {
		user.Name = v.(string)
	} else {
		log.Error("name not found in user info claims")
	}

	user.Issuer = iDToken.Issuer
	user.IssuedAt = iDToken.IssuedAt

	domain := os.Getenv("OAUTH2_PROVIDER_URL")
	id := os.Getenv("OAUTH2_CLIENT_ID")
	secret := os.Getenv("OAUTH2_CLIENT_SECRET")
	m, err := management.New(domain, id, secret)
	if err != nil {
		log.Errorf("Error talking to auth0: %v", err)
		// handle err
	}
	u, err := m.User.Read(user.Sub)
	if err != nil {
		log.Errorf("Error reading user %s %v", user.Sub, err)
	}

	if u != nil {
		log.Infof("User: %v", u)
		if u.UserMetadata["Plan"] != nil {
			user.Plan = u.UserMetadata["Plan"].(string)
		}
	}

	user.Picture = *u.Picture

	log.Infof("user.Plan: %s", user.Plan)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_CONNECTION_STRING")))

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

	accounts := mongodb.ReadAllAccounts(user.Email)
	if err != nil {
		log.Error(err)
	}

	if len(accounts) == 0 {

		var account model.Account
		account.Email = user.Email
		a, err := core.CreateAccount(&account)
		log.Infof("account = %v", a)
		if err != nil {
			log.Error(err)
		}

	}

	//res, err := collection.InsertOne(ctx, b)

	log.Infof("Res: %v", res)

	return user, nil
}
