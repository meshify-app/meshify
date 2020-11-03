package model

import (
	"fmt"
	"time"

	"github.com/alan-grapid/meshify/util"
)

// Host structure
type Host struct {
	Id                        string    `json:"id"                        bson:"id"`
	MeshID                    string    `json:"meshid"                    bson:"meshid"`
	MeshName                  string    `json:"meshName"                  bson:"meshName"`
	PrivateKey                string    `json:"privateKey"                bson:"privateKey"`
	PublicKey                 string    `json:"publicKey"                 bson:"publicKey"`
	Name                      string    `json:"name"                      bson:"name"`
	Email                     string    `json:"email"                     bson:"email"`
	Enable                    bool      `json:"enable"         			  bson:"enable"`
	IgnorePersistentKeepalive bool      `json:"ignorePersistentKeepalive" bson:"ignorePersistentKeepalive"`
	PresharedKey              string    `json:"presharedKey"              bson:"presharedKey"`
	AllowedIPs                []string  `json:"allowedIPs"                bson:"allowedIPs"`
	Address                   []string  `json:"address" 				  bson:"address"`
	Tags                      []string  `json:"tags"                      bson:"tags"`
	Dns                       []string  `json:"dns"                       bson:"dns"`
	PersistentKeepalive       int       `json:"persistentKeepalive"       bson:"persistentKeepAlive"`
	ListenPort                int       `json:"listenPort"                bson:"listenPort"`
	Endpoint                  string    `json:"endpoint"                  bson:"endpoint"`
	Mtu                       int       `json:"mtu"                       bson:"mtu"`
	PreUp                     string    `json:"preUp"                     bson:"preUp"`
	PostUp                    string    `json:"postUp"                    bson:"postUp"`
	PreDown                   string    `json:"preDown"                   bson:"preDown"`
	PostDown                  string    `json:"postDown"                  bson:"postDown"`
	CreatedBy                 string    `json:"createdBy"                 bson:"createdBy"`
	UpdatedBy                 string    `json:"updatedBy"                 bson:"updatedBy"`
	Created                   time.Time `json:"created"                   bson:"created"`
	Updated                   time.Time `json:"updated"                   bson:"updated"`
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
