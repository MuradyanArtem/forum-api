package models

import (
	"time"

	"github.com/jackc/pgtype"
)

//easyjson:json
type Post struct {
	ID       int              `json:"id"`
	Author   string           `json:"author"`
	Created  time.Time        `json:"created"`
	Forum    string           `json:"forum"`
	IsEdited bool             `json:"isEdited"`
	Message  string           `json:"message"`
	Parent   int              `json:"parent"`
	Thread   int              `json:"thread"`
	Path     pgtype.Int8Array `json:"-"`
}

//easyjson:json
type PostSlice []Post

//easyjson:json
type PostDetails struct {
	Post   *Post   `json:"post"`
	Author *User   `json:"author,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
}
