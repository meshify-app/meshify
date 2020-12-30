package model

import (
	"fmt"
	"time"
)

// Mesh structure
type Mesh struct {
	Id        string    `json:"id"          bson:"id"`
	MeshName  string    `json:"meshName"    bson:"meshName"`
	CreatedBy string    `json:"createdBy"   bson:"createdBy"`
	UpdatedBy string    `json:"updatedBy"   bson:"updatedBy"`
	Created   time.Time `json:"created"     bson:"created"`
	Updated   time.Time `json:"updated"     bson:"updated"`
	Default   Settings  `json:"default"     bson:"default"`
}

// IsValid check if model is valid
func (a Mesh) IsValid() []error {
	errs := make([]error, 0)

	if a.Id == "" {
		errs = append(errs, fmt.Errorf("Id is required"))
	}

	// check if the name empty
	if a.MeshName == "" {
		errs = append(errs, fmt.Errorf("name is required"))
	}
	// check the name field is between 3 to 40 chars
	if len(a.MeshName) < 2 || len(a.MeshName) > 12 {
		errs = append(errs, fmt.Errorf("name field must be between 2-12 chars"))
	}

	return errs
}
