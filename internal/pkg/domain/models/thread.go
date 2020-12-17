package models

import "time"

//easyjson:json
type Thread struct {
	Id      uint64    `json:"id"`
	Slug    string    `json:"slug"`
	Title   string    `json:"title"`
	Message string    `json:"message"`
	Forum   string    `json:"forum"`
	Author  string    `json:"author"`
	Created time.Time `json:"created,omitempty"`
	Votes   uint64    `json:"votes"`
}

//easyjson:json
type ThreadUpdate struct {
	Message string `json:"message"`
	Title   string `json:"title"`
}
