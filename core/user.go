package core

import (
	"errors"
	"os"
	"reflect"
	"sort"
	"strconv"
	"time"

	model "github.com/meshify-app/meshify/model"
	mongo "github.com/meshify-app/meshify/mongo"
	template "github.com/meshify-app/meshify/template"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

// CreateUser user with all necessary data
func CreateUser(user *model.User) (*model.User, error) {

	user.Created = time.Now().UTC()
	user.Updated = user.Created

	// check if user is valid
	errs := user.IsValid()
	if len(errs) != 0 {
		for _, err := range errs {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("user validation error")
		}
		return nil, errors.New("failed to validate user")
	}

	err := mongo.Serialize(user.Email, "email", "users", user)
	if err != nil {
		return nil, err
	}

	v, err := mongo.Deserialize(user.Email, "email", "users", reflect.TypeOf(model.User{}))
	if err != nil {
		return nil, err
	}
	user = v.(*model.User)

	// data modified, dump new config
	return user, nil
}

// ReadUser user by id
func ReadUser(id string) (*model.User, error) {
	v, err := mongo.Deserialize(id, "email", "users", reflect.TypeOf(model.User{}))
	if err != nil {
		return nil, err
	}
	user := v.(*model.User)

	return user, nil
}

// UpdateUser preserve keys
func UpdateUser(Id string, user *model.User) (*model.User, error) {
	v, err := mongo.Deserialize(Id, "email", "users", reflect.TypeOf(model.User{}))
	if err != nil {
		return nil, err
	}
	current := v.(*model.User)

	if current != nil && user != nil &&
		current.Email != user.Email {
		return nil, errors.New("records Id mismatch")
	}

	// check if user is valid
	errs := user.IsValid()
	if len(errs) != 0 {
		for _, err := range errs {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("user validation error")
		}
		return nil, errors.New("failed to validate user")
	}

	// keep keys
	user.Updated = time.Now().UTC()

	err = mongo.Serialize(user.Email, "email", "users", user)
	if err != nil {
		return nil, err
	}

	v, err = mongo.Deserialize(Id, "email", "users", reflect.TypeOf(model.User{}))
	if err != nil {
		return nil, err
	}
	user = v.(*model.User)

	// data modified, dump new config
	return user, nil
}

// DeleteUser from database
func DeleteUser(id string) error {

	return mongo.Delete(id, "id", "users")
}

// ReadUsers all users
func ReadUsers() ([]*model.User, error) {
	users := make([]*model.User, 0)
	users = mongo.ReadAllUsers()

	sort.Slice(users, func(i, j int) bool {
		return users[i].Created.After(users[j].Created)
	})

	return users, nil
}

// EmailHost send email to host
func EmailUser(id string, account string, mesh string) error {
	// get email body
	emailBody, err := template.DumpUserEmail(account, mesh)
	if err != nil {
		return err
	}

	// port to int
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return err
	}

	d := gomail.NewDialer(os.Getenv("SMTP_HOST"), port, os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"))
	s, err := d.Dial()
	if err != nil {
		return err
	}
	m := gomail.NewMessage()

	m.SetHeader("From", os.Getenv("SMTP_FROM"))
	m.SetAddressHeader("To", id, id)
	m.SetHeader("Subject", "Meshify.app Invitation")
	m.SetBody("text/html", string(emailBody))

	err = gomail.Send(s, m)
	if err != nil {
		return err
	}

	return nil
}
