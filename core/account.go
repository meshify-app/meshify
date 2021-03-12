package core

import (
	"errors"
	"reflect"
	"strings"
	"time"

	model "github.com/meshify-app/meshify/model"
	mongo "github.com/meshify-app/meshify/mongo"
	util "github.com/meshify-app/meshify/util"
	log "github.com/sirupsen/logrus"
)

// CreateAccount with all necessary data
func CreateAccount(account *model.Account) (*model.Account, error) {

	var err error

	if account.Id == "" {
		account.Id, err = util.RandomString(16)
		if err != nil {
			return nil, err
		}
	}

	if account.Key == "" {
		account.Key, err = util.RandomString(32)
		if err != nil {
			return nil, err
		}
	}

	if account.Parent == "" {
		account.Parent = account.Id
	}

	account.Created = time.Now()

	errs := account.IsValid()
	if len(errs) != 0 {
		for _, err := range errs {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("account validation error")
		}
		return nil, errors.New("failed to validate account")
	}

	err = mongo.Serialize(account.Id, "id", "accounts", account)

	if err != nil {
		return nil, err
	}

	v, err := mongo.Deserialize(account.Id, "id", "accounts", reflect.TypeOf(model.Account{}))
	if err != nil {
		return nil, err
	}
	account = v.(*model.Account)

	// return current account
	return account, nil
}

// ReadHost host by id
func ReadAllAccountsForUser(email string) ([]*model.Account, error) {

	if strings.Contains(email, "@") {
		return mongo.ReadAllAccounts(email)
	} else {
		return mongo.ReadAllAccountsForID(email)
	}
}

// DeleteHost from disk
func DeleteAccount(id string) error {

	return mongo.Delete(id, "id", "accounts")
}

// DeleteHost from disk
func ActivateAccount(id string) (string, error) {

	var a *model.Account

	v, err := mongo.Deserialize(id, "id", "accounts", reflect.TypeOf(model.Account{}))
	if err != nil {
		return "Error", err
	}
	a = v.(*model.Account)
	a.Status = "Active"

	err = mongo.Serialize(id, "id", "accounts", a)
	if err != nil {
		return "Error", err
	}

	return "Account activated.", nil
}
