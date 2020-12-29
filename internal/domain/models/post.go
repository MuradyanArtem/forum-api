package models

import (
	"github.com/jackc/pgtype"
)

//easyjson:json
type Post struct {
	ID       int64            `json:"id"`
	Author   string           `json:"author"`
	Created  string           `json:"created"`
	Forum    string           `json:"forum"`
	IsEdited bool             `json:"isEdited"`
	Message  string           `json:"message"`
	Parent   int64            `json:"parent"`
	Thread   int64            `json:"thread"`
	Path     pgtype.Int8Array `json:"-"`
}

//easyjson:json
type Posts []Post

//easyjson:json
type PostFull struct {
	Post   *Post   `json:"post"`
	Author *User   `json:"author,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
}
