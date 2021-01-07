package models

import "time"

//easyjson:json
type Thread struct {
	ID      int       `json:"id"`
	Slug    string    `json:"slug,omitempty"`
	Title   string    `json:"title"`
	Message string    `json:"message"`
	Forum   string    `json:"forum"`
	Author  string    `json:"author"`
	Created time.Time `json:"created,omitempty"`
	Votes   int       `json:"votes"`
}

//easyjson:json
type ThreadSlice []Thread
