package comment

import (
	"encoding/json"
	"log"
	"time"
)

type Comment struct {
	ID          int       `json:"id"`
	PostID      int       `json:"post_id"`
	DateCreated time.Time `json:"date_created"`
	Details     string    `json:"details"`
	IsVisible   bool      `json:"is_visible"`
}

func (c *Comment) ToJSON() []byte {
	b, err := json.Marshal(c)
	if err != nil {
		log.Fatalf("error marshaling (%T): %v\n", c, err)
	}
	return b
}

func (c *Comment) FromJSON(b []byte) {
	err := json.Unmarshal(b, &c)
	if err != nil {
		log.Fatalf("error unmarshaling (%T): %v\n", c, err)
	}
}
