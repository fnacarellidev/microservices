package api

import "time"

type Record struct {
	Id    string    `json:"id,omitempty"`
	Title string    `json:"title"`
	Text  string    `json:"content"`
	Date  time.Time `json:"created_at,omitempty"`
}
