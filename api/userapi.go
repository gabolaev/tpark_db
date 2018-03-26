package api

import (
	"fmt"

	"github.com/gabolaev/tpark_db/helpers"
	"github.com/gabolaev/tpark_db/models"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func CreateUser(context *fasthttp.RequestCtx) {
	var user models.User
	body := context.PostBody()
	if err := easyjson.Unmarshal(body, &user); err != nil {
		context.SetStatusCode(fasthttp.StatusBadRequest)
		context.WriteString(err.Error())
		return
	}

	result, created, err := helpers.CreateNewOrGetExistingUsers(&user)
	if err != nil {
		context.SetStatusCode(fasthttp.StatusInternalServerError)
		errorJSON, _ := easyjson.Marshal(models.Error{Message: err.Error()})
		context.SetBody(errorJSON)
	}
	if created == false {
		context.SetStatusCode(fasthttp.StatusConflict)
		if existingUsersJSON, err := easyjson.Marshal(result); err != nil {
			context.SetStatusCode(fasthttp.StatusInternalServerError)
		} else {
			context.SetBody(existingUsersJSON)
			return
		}
	}

	context.SetStatusCode(fasthttp.StatusCreated)
	context.SetBody(body)

}

func GetUser(context *fasthttp.RequestCtx) {
	fmt.Println(context)
}

func UpdateUser(context *fasthttp.RequestCtx) {
	fmt.Println(context)
}
