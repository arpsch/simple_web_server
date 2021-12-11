package model

import (
	"time"
)

// User object to store the user details
type User struct {
	ID string `json:"id,required"`

	Name string `json:"name,required"`

	SignupTime time.Time `json:"signupTime,required"`
}
