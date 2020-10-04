package model

import "time"

// User structure
type User struct {
	ID       string    `json:"_id"`
	Sub      string    `json:"sub"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Profile  string    `json:"profile"`
	Issuer   string    `json:"issuer"`
	Plan     string    `json:"plan"`
	IssuedAt time.Time `json:"issuedAt"`
}
