package models

//easyjson:json
type Thread struct {
	ID      int64  `json:"id"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Forum   string `json:"forum"`
	Author  string `json:"author"`
	Created string `json:"created,omitempty"`
	Votes   int64  `json:"votes"`
}

//easyjson:json
type Threads []Thread
