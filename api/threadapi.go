package api

import (
	"github.com/gabolaev/tpark_db/errors"
	"github.com/gabolaev/tpark_db/helpers"

	"github.com/gabolaev/tpark_db/models"

	"github.com/valyala/fasthttp"
)

func CreateThreadOrForum(context *fasthttp.RequestCtx) {
	context.SetContentType("application/json")
	param := context.UserValue("catch-all-param").(string)
	if param == "/create" {
		CreateForum(context)
		return
	}
	var thread models.Thread
	body := context.PostBody()
	if err := thread.UnmarshalJSON(body); err != nil {
		context.SetStatusCode(fasthttp.StatusBadRequest)
		context.WriteString(err.Error())
		return
	}

	thread.Forum = param[1 : len(param)-7]
	result, err := helpers.CreateNewOrGetExistingThread(&thread)
	var responseStatus int
	switch err {
	case nil:
		responseStatus = fasthttp.StatusCreated
	case errors.ConflictError:
		responseStatus = fasthttp.StatusConflict
	case errors.NotFoundError:
		err := models.Error{Message: "Can't find user or forum"}
		errorJSON, _ := err.MarshalJSON()
		context.SetStatusCode(fasthttp.StatusNotFound)
		context.SetBody(errorJSON)
		return
	default:
		context.SetStatusCode(fasthttp.StatusInternalServerError)
		context.SetBodyString(err.Error())
		return
	}
	if existingForumJSON, err := result.MarshalJSON(); err != nil {
		context.SetStatusCode(fasthttp.StatusInternalServerError)
		context.SetBodyString(err.Error())
	} else {
		context.SetStatusCode(responseStatus)
		context.SetBody(existingForumJSON)
	}
}
