package models

import (
	"encoding/json"
	"log"
	"time"
)

type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DOB         time.Time `json:"dob"`
	LastLogin   time.Time `json:"last_login"`
	DateCreated time.Time `json:"date_created"`
	IsActive    bool      `json:"is_active"`
}

func (u *User) ToJSON() []byte {
	b, err := json.Marshal(u)
	if err != nil {
		log.Fatalf("error marshaling (%T): %v\n", u, err)
	}
	return b
}

func (u *User) FromJSON(b []byte) {
	err := json.Unmarshal(b, &u)
	if err != nil {
		log.Fatalf("error unmarshaling (%T): %v\n", u, err)
	}
}
