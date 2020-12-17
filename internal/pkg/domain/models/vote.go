package models

//easyjson:json
type Vote struct {
	Vote     uint32  `json:"voice"`
	Nickname string `json:"nickname"`
	Thread   uint64  `json:"thread,omitempty"`
}
