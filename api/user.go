package api

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

//easyjson:json
type User struct {
	Nickname string
	Email    string
	Fullname string
	About    string
}

//easyjson:json
type Users []User

func Create(context *fasthttp.RequestCtx) (*User, error) {
	fmt.Println("Create user")
	return nil, nil
}

func Get(nickname string) (*User, error) {
	fmt.Println("Get user")
	return nil, nil
}

func Update(context *fasthttp.RequestCtx) (*User, error) {
	fmt.Println("Update user")
	return nil, nil
}
