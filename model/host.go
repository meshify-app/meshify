package model

import (
	"fmt"
	"time"

	"github.com/grapid/meshify/util"
)

// Host structure
type Host struct {
	Id        string    `json:"id"                        bson:"id"`
	Name      string    `json:"name"                      bson:"name"`
	MeshId    string    `json:"meshid"                    bson:"meshid"`
	MeshName  string    `json:"meshName"                  bson:"meshName"`
	Email     string    `json:"email"                     bson:"email"`
	Enable    bool      `json:"enable"                    bson:"enable"`
	Tags      []string  `json:"tags"                      bson:"tags"`
	CreatedBy string    `json:"createdBy"                 bson:"createdBy"`
	UpdatedBy string    `json:"updatedBy"                 bson:"updatedBy"`
	Created   time.Time `json:"created"                   bson:"created"`
	Updated   time.Time `json:"updated"                   bson:"updated"`
	Current   Settings  `json:"current"                   bson:"current"`
	Default   Settings  `json:"default"                   bson:"default"`
}

// IsValid check if model is valid
func (a Host) IsValid() []error {
	errs := make([]error, 0)

	// check if the name empty
	if a.Name == "" {
		errs = append(errs, fmt.Errorf("name is required"))
	}
	// check the name field is between 3 to 40 chars
	if len(a.Name) < 2 || len(a.Name) > 40 {
		errs = append(errs, fmt.Errorf("name field must be between 2-40 chars"))
	}
	// email is not required, but if provided must match regex
	if a.Email != "" {
		if !util.RegexpEmail.MatchString(a.Email) {
			errs = append(errs, fmt.Errorf("email %s is invalid", a.Email))
		}
	}
	/*	// check if the allowedIPs empty
		if len(a.AllowedIPs) == 0 {
			errs = append(errs, fmt.Errorf("allowedIPs field is required"))
		}
		// check if the allowedIPs are valid
		for _, allowedIP := range a.AllowedIPs {
			if !util.IsValidCidr(allowedIP) {
				errs = append(errs, fmt.Errorf("allowedIP %s is invalid", allowedIP))
			}
		}
	*/ // check if the address empty

	if len(a.Current.Address) == 0 {
		errs = append(errs, fmt.Errorf("address field is required"))
	}
	// check if the address are valid
	for _, address := range a.Current.Address {
		if !util.IsValidCidr(address) {
			errs = append(errs, fmt.Errorf("address %s is invalid", address))
		}
	}

	return errs
}
