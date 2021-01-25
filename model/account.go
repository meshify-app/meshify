package model

import (
	"fmt"
	"time"

	"github.com/meshify-app/meshify/util"
)

type Account struct {
	Id      string    `json:"id"       bson:"id"`
	Email   string    `json:"email"    bson:"email"`
	Key     string    `json:"key"      bson:"key"`
	Created time.Time `json:"created"  bson:"created"`
}

type AccountInfo struct {
	Id   string `json:"id"       bson:"id"`
	Name string `json:"name"     bson:"name"`
}

// IsValid check if model is valid
func (a Account) IsValid() []error {
	errs := make([]error, 0)

	// check if the name empty
	if a.Id == "" {
		errs = append(errs, fmt.Errorf("id is required"))
	}
	// email is required, but if provided must match regex
	if a.Email != "" {
		if !util.RegexpEmail.MatchString(a.Email) {
			errs = append(errs, fmt.Errorf("email %s is invalid", a.Email))
		}
	} else {
		errs = append(errs, fmt.Errorf("email is required."))
	}
	return errs
}
