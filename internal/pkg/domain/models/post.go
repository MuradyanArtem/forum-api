package models

import (
	"time"
)

//easyjson:json
type Post struct {
	Author   string    `json:"author"`
	Created  time.Time `json:"created"`
	Forum    string    `json:"forum"`
	ID       uint64    `json:"id"`
	IsEdited bool      `json:"isEdited"`
	Message  string    `json:"message"`
	Parent   uint64    `json:"parent"`
	Thread   uint64    `json:"thread"`
	Parents  []int32
}
