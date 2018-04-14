package api

import (
	"fmt"

	"github.com/gabolaev/tpark_db/database"
	"github.com/gabolaev/tpark_db/errors"
	"github.com/gabolaev/tpark_db/helpers"
	"github.com/gabolaev/tpark_db/models"
	"github.com/valyala/fasthttp"
)

func CreateUser(context *fasthttp.RequestCtx) {
	context.SetContentType("application/json")
	var user models.User
	if err := user.UnmarshalJSON(context.PostBody()); err != nil {
		context.SetStatusCode(fasthttp.StatusBadRequest)
		context.WriteString(err.Error())
		return
	}
	user.Nickname = context.UserValue("nickname").(string)

	result, err := helpers.CreateNewOrGetExistingUsers(&user)
	switch err {
	case errors.ConflictError:
		if responseBody, err := result.MarshalJSON(); err != nil {
			context.SetStatusCode(fasthttp.StatusInternalServerError)
		} else {
			context.SetStatusCode(fasthttp.StatusConflict)
			context.SetBody(responseBody)
		}
	case nil:
		if responseBody, err := result.MarshalJSON(); err != nil {
			context.SetStatusCode(fasthttp.StatusInternalServerError)
		} else {
			defer func() {
				database.Instance.Status.User++
			}()
			context.SetStatusCode(fasthttp.StatusCreated)
			context.SetBody(responseBody[1 : len(responseBody)-1])
		}
	default:
		context.SetStatusCode(fasthttp.StatusInternalServerError)
		context.SetBodyString(err.Error())
	}
}

func GetUser(context *fasthttp.RequestCtx) {
	context.SetContentType("application/json")
	nickname := context.UserValue("nickname").(string)
	result, err := helpers.GetUserByNickname(&nickname)
	switch err {
	case nil:
		if user, err := result.MarshalJSON(); err != nil {
			context.SetStatusCode(fasthttp.StatusInternalServerError)
		} else {
			context.SetStatusCode(fasthttp.StatusOK)
			context.SetBody(user)
		}
	case errors.NotFoundError:
		context.SetStatusCode(fasthttp.StatusNotFound)
		err := models.Error{
			Message: fmt.Sprintf("Can't find user with nickname: %s", nickname)}
		errorJSON, _ := err.MarshalJSON()
		context.SetBody(errorJSON)
	}
}

func UpdateUser(context *fasthttp.RequestCtx) {
	context.SetContentType("application/json")
	nickname := context.UserValue("nickname").(string)
	var user models.User

	if err := user.UnmarshalJSON(context.PostBody()); err != nil {
		context.SetStatusCode(fasthttp.StatusBadRequest)
		context.WriteString(err.Error())
		return
	}
	user.Nickname = nickname
	err := helpers.UpdateUserInfo(&user)

	var errObj models.Error
	var responseStatus int
	switch err {
	case nil:
		if updatedUser, err := user.MarshalJSON(); err != nil {
			context.SetStatusCode(fasthttp.StatusInternalServerError)
		} else {
			context.SetStatusCode(fasthttp.StatusOK)
			context.SetBody(updatedUser)
		}
		return
	case errors.NotFoundError:
		errObj = models.Error{Message: "User not found"}
		responseStatus = fasthttp.StatusNotFound
	case errors.ConflictError:
		errObj = models.Error{Message: "New user profile data conflicts with existing users."}
		responseStatus = fasthttp.StatusConflict
	default:
		context.SetStatusCode(fasthttp.StatusInternalServerError)
		context.SetBodyString(err.Error())
		return
	}
	context.SetStatusCode(responseStatus)
	errorJSON, _ := errObj.MarshalJSON()
	context.SetBody(errorJSON)
}
