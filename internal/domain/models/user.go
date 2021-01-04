package models

//easyjson:json
type User struct {
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname,omitempty"`
	About    string `json:"about,omitempty"`
	Email    string `json:"email"`
}

//easyjson:json
type UserSlice []User
