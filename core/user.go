package core

import (
	"errors"
	"os"
	"reflect"
	"sort"
	"strconv"
	"time"

	model "github.com/grapid/meshify/model"
	mongo "github.com/grapid/meshify/mongo"
	template "github.com/grapid/meshify/template"
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
	return user, UpdateServerConfigWg()
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
	return user, UpdateServerConfigWg()
}

// DeleteUser from disk
func DeleteUser(id string) error {

	err := mongo.DeleteHost(id, "users")
	//	path := filepath.Join(os.Getenv("WG_CONF_DIR"), id)
	//	err := os.Remove(path)
	if err != nil {
		return err
	}

	// data modified, dump new config
	return UpdateServerConfigWg()
}

// ReadUsers all users
func ReadUsers() ([]*model.User, error) {
	users := make([]*model.User, 0)
	/*
		files, err := ioutil.ReadDir(filepath.Join(os.Getenv("WG_CONF_DIR")))
		if err != nil {
			return nil, err
		}

		for _, f := range files {
			// users file name is an uuid
			_, err := uuid.FromString(f.Name())
			if err == nil {
				c, err := mongo.Deserialize(f.Name())
				if err != nil {
					log.WithFields(log.Fields{
						"err":  err,
						"path": f.Name(),
					}).Error("failed to deserialize user")
				} else {
					users = append(users, c.(*model.User))
				}
			}
		}
	*/
	users = mongo.ReadAllUsers()

	sort.Slice(users, func(i, j int) bool {
		return users[i].Created.After(users[j].Created)
	})

	return users, nil
}

// ReadUserConfig in wg format
func ReadUserConfig(id string) ([]byte, error) {
	user, err := ReadHost(id)
	if err != nil {
		return nil, err
	}

	server, err := ReadServer()
	if err != nil {
		return nil, err
	}

	configDataWg, err := template.DumpClientWg(user, server)
	if err != nil {
		return nil, err
	}

	return configDataWg, nil
}

// EmailHost send email to host
func EmailUser(id string) error {
	// get email body
	emailBody, err := template.DumpUserEmail()
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
