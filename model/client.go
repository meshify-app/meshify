package model

import (
	"fmt"
	"time"

	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/util"
)

// Host structure
type Host struct {
	Id                        string    `json:"id"`
	MeshID                    string    `json:"meshid"`
	MeshName                  string    `json:"meshName"`
	PrivateKey                string    `json:"privateKey"`
	PublicKey                 string    `json:"publicKey"`
	Name                      string    `json:"name"`
	Email                     string    `json:"email"`
	Enable                    bool      `json:"enable"`
	IgnorePersistentKeepalive bool      `json:"ignorePersistentKeepalive"`
	PresharedKey              string    `json:"presharedKey"`
	AllowedIPs                []string  `json:"allowedIPs"`
	Address                   []string  `json:"address"`
	Tags                      []string  `json:"tags"`
	Dns                       []string  `json:"dns"`
	PersistentKeepalive       int       `json:"persistentKeepalive"`
	ListenPort                int       `json:"listenPort"`
	Endpoint                  string    `json:"endpoint"`
	Mtu                       int       `json:"mtu"`
	PreUp                     string    `json:"preUp"`
	PostUp                    string    `json:"postUp"`
	PreDown                   string    `json:"preDown"`
	PostDown                  string    `json:"postDown"`
	CreatedBy                 string    `json:"createdBy"`
	UpdatedBy                 string    `json:"updatedBy"`
	Created                   time.Time `json:"created"`
	Updated                   time.Time `json:"updated"`
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
	if len(a.Address) == 0 {
		errs = append(errs, fmt.Errorf("address field is required"))
	}
	// check if the address are valid
	for _, address := range a.Address {
		if !util.IsValidCidr(address) {
			errs = append(errs, fmt.Errorf("address %s is invalid", address))
		}
	}

	return errs
}
