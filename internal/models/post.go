package models

import (
	"encoding/json"
	"log"
	"time"
)

type Post struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	DateCreated time.Time `json:"date_created"`
	Title       string    `json:"title"`
	Details     string    `json:"details"`
	CanComment  string    `json:"can_comment"`
}

func (p *Post) ToJSON() []byte {
	b, err := json.Marshal(p)
	if err != nil {
		log.Fatalf("error marshaling (%T): %v\n", p, err)
	}
	return b
}

func (p *Post) FromJSON(b []byte) {
	err := json.Unmarshal(b, &p)
	if err != nil {
		log.Fatalf("error unmarshaling (%T): %v\n", p, err)
	}
}
