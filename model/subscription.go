package model

import (
	"fmt"
	"time"
)

// Mesh structure
type Subscription struct {
	Id          string    `json:"id"          bson:"id"`
	AccountId   string    `json:"accountid"   bson:"accountid"`
	Email       string    `json:"email"       bson:"email"`
	Name        string    `json:"name"        bson:"name"`
	Description string    `json:"description" bson:"description"`
	Created     time.Time `json:"created"     bson:"created"`
	Updated     time.Time `json:"updated"     bson:"updated"`
	CreatedBy   string    `json:"createdBy"   bson:"createdBy"`
	UpdatedBy   string    `json:"updatedBy"   bson:"updatedBy"`
}

// IsValid check if model is valid
func (s Subscription) IsValid() []error {
	errs := make([]error, 0)

	if s.Id == "" {
		errs = append(errs, fmt.Errorf("id is required"))
	}

	return errs
}
