package model

import (
	"fmt"
	"time"

	"github.com/meshify-app/meshify/util"
)

type Account struct {
	Id      string    `json:"id"       bson:"id"`
	Parent  string    `json:"parent"   bson:"parent"`
	Email   string    `json:"email"    bson:"email"`
	Role    string    `json:"role"     bson:"role"`
	Status  string    `json:"status"   bson:"status"`
	Key     string    `json:"key"      bson:"key"`
	Created time.Time `json:"created"  bson:"created"`
}

type Organization struct {
	Id      string    `json:"id"       bson:"id"`
	Parent  string    `json:"parent"   bson:"parent"`
	Name    string    `json:"name"     bson:"name"`
	Created time.Time `json:"created"  bson:"created"`
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

// IsValid check if model is valid
func (a Organization) IsValid() []error {
	errs := make([]error, 0)

	// check if the name empty
	if a.Id == "" {
		errs = append(errs, fmt.Errorf("id is required"))
	}
	// name is required
	if a.Name == "" {
		errs = append(errs, fmt.Errorf("name is required."))
	}

	return errs
}
