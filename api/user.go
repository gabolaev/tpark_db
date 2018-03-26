package api

import (
	"fmt"

	"github.com/mailru/easyjson"
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
	var user User
	if err := easyjson.Unmarshal(context.PostBody(), &user); err != nil {
		context.SetStatusCode(fasthttp.StatusBadRequest)
		context.WriteString(err.Error())
		return
	}

	if responseJSON, err := easyjson.Marshal(user); err != nil {
		context.SetStatusCode(fasthttp.StatusInternalServerError)
		context.WriteString(err.Error())
	} else {
		context.SetStatusCode(fasthttp.StatusCreated)
		context.SetBody(responseJSON)
	}
}

func GetUser(context *fasthttp.RequestCtx) {
	fmt.Println(context)
}

func UpdateUser(context *fasthttp.RequestCtx) {
	fmt.Println(context)
}
