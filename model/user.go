package model

import (
	"fmt"
	"time"

	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/util"
)

// User structure
type User struct {
	Sub       string    `json:"sub"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Profile   string    `json:"profile"`
	Issuer    string    `json:"issuer"`
	Plan      string    `json:"plan"`
	IssuedAt  time.Time `json:"issuedAt"`
	AccountID string    `json:"accountId"`
	CreatedBy string    `json:"createdBy"`
	Created   time.Time `json:"created_at"`
	UpdatedBy string    `json:"updatedBy"`
	Updated   time.Time `json:"updated_at"`
}

// IsValid check if model is valid
func (a User) IsValid() []error {
	errs := make([]error, 0)

	// check if the name empty
	if a.Name == "" {
		errs = append(errs, fmt.Errorf("name is required"))
	}

	if a.Email == "" {
		errs = append(errs, fmt.Errorf("email is required"))
	}
	// email is not required, but if provided must match regex
	if a.Email != "" {
		if !util.RegexpEmail.MatchString(a.Email) {
			errs = append(errs, fmt.Errorf("email %s is invalid", a.Email))
		}
	}

	return errs
}
