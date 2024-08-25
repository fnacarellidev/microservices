package api

import "time"

type Post struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
