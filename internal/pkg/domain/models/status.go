package models

//easyjson:json
type Status struct {
	Forum  uint64 `json:"forum"`
	Post   uint64 `json:"post"`
	Thread uint64 `json:"thread"`
	User   uint64 `json:"user"`
}
