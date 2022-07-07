package model

import (
	"fmt"
	"time"
)

// Mesh structure
type Service struct {
	Id        string    `json:"id"          bson:"id"`
	AccountId string    `json:"accountid"   bson:"accountid"`
	Email     string    `json:"email"       bson:"email"`
	Created   time.Time `json:"created"     bson:"created"`
	Updated   time.Time `json:"updated"     bson:"updated"`
	CreatedBy string    `json:"createdBy"   bson:"createdBy"`
	UpdatedBy string    `json:"updatedBy"   bson:"updatedBy"`
	MeshName  string    `json:"meshName"    bson:"meshName"`
	MeshId    string    `json:"meshId"      bson:"meshId"`
	RelayHost Host      `json:"relayHost"   bson:"relayHost"`
}

// IsValid check if model is valid
func (s Service) IsValid() []error {
	errs := make([]error, 0)

	if s.Id == "" {
		errs = append(errs, fmt.Errorf("id is required"))
	}

	return errs
}
