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
	user.Nickname = context.UserValue("nickname").(string)

	result, created, err := helpers.CreateNewOrGetExistingUsers(&user)
	if err != nil {
		context.SetStatusCode(fasthttp.StatusInternalServerError)
		errorJSON, _ := easyjson.Marshal(models.Error{Message: err.Error()})
		context.SetBody(errorJSON)
	}

	if responseBody, err := easyjson.Marshal(result); err != nil {
		context.SetStatusCode(fasthttp.StatusInternalServerError)
	} else {
		if created {
			context.SetBody(responseBody[1 : len(responseBody)-1])
			context.SetStatusCode(fasthttp.StatusCreated)
		} else {
			context.SetBody(responseBody)
			context.SetStatusCode(fasthttp.StatusConflict)
		}
	}
}

func GetUser(context *fasthttp.RequestCtx) {
	nickname := context.UserValue("nickname").(string)
	result, err := helpers.GetUserByNickname(nickname)
	if err != nil {
		context.SetStatusCode(fasthttp.StatusNotFound)
		errorJSON, _ := easyjson.Marshal(models.Error{
			Message: fmt.Sprintf("Can't find user with nickname: %s", nickname)})
		context.SetBody(errorJSON)
	} else {
		if user, err := easyjson.Marshal(result); err != nil {
			context.SetStatusCode(fasthttp.StatusInternalServerError)
		} else {
			context.SetStatusCode(fasthttp.StatusOK)
			context.SetBody(user)
		}
	}

}

func UpdateUser(context *fasthttp.RequestCtx) {
	nickname := context.UserValue("nickname").(string)
	var user models.User
	body := context.PostBody()
	if err := easyjson.Unmarshal(body, &user); err != nil {
		context.SetStatusCode(fasthttp.StatusBadRequest)
		context.WriteString(err.Error())
		return
	}
	user.Nickname = nickname
	updated, err := helpers.UpdateUserInfo(&user)
	if err != nil {
		context.SetStatusCode(fasthttp.StatusConflict)
		errorJSON, _ := easyjson.Marshal(models.Error{
			Message: "New user profile data conflicts with existing users."})
		context.SetBody(errorJSON)
		return
	}
	if updated {
		context.SetStatusCode(fasthttp.StatusOK)
		if updatedUser, err := easyjson.Marshal(user); err != nil {
			context.SetStatusCode(fasthttp.StatusInternalServerError)
		} else {
			context.SetBody(updatedUser)
		}
	} else {
		context.SetStatusCode(fasthttp.StatusNotFound)
		errorJSON, _ := easyjson.Marshal(models.Error{
			Message: "User not found"})
		context.SetBody(errorJSON)
	}
}
