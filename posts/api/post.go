package api

import "time"

type Post struct {
	Id        string    `json:"id,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
