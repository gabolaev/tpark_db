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

func CreateUser(context *fasthttp.RequestCtx) {
	fmt.Println(context)
}

func GetUser(context *fasthttp.RequestCtx) {
	fmt.Println(context)
}

func UpdateUser(context *fasthttp.RequestCtx) {
	fmt.Println(context)
}
