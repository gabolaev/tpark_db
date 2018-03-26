package models

//easyjson:json
type User struct {
	Nickname string
	Email    string
	Fullname string
	About    string
}

//easyjson:json
type Users []*User
