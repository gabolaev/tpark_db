package api

import (
	"github.com/gabolaev/tpark_db/helpers"

	"github.com/gabolaev/tpark_db/models"
	"github.com/mailru/easyjson"

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
	if err := easyjson.Unmarshal(body, &thread); err != nil {
		context.SetStatusCode(fasthttp.StatusBadRequest)
		context.WriteString(err.Error())
		return
	}

	thread.Forum = param[1 : len(param)-7]
	result, created, err := helpers.CreateNewOrGetExistingThread(&thread)
	if err != nil {
		context.SetStatusCode(fasthttp.StatusInternalServerError)
		context.SetBodyString(err.Error())
		return
	}
	if created {
		if createdForumJSON, err := easyjson.Marshal(result); err != nil {
			context.SetStatusCode(fasthttp.StatusInternalServerError)
			context.SetBodyString(err.Error())
		} else {
			context.SetStatusCode(fasthttp.StatusCreated)
			context.SetBody(createdForumJSON)
		}
	} else {
		if result == nil {
			errorJSON, _ := easyjson.Marshal(models.Error{
				Message: "Can't find user or forum"})
			context.SetStatusCode(fasthttp.StatusNotFound)
			context.SetBody(errorJSON)
		} else {
			if existingForumJSON, err := easyjson.Marshal(result); err != nil {
				context.SetStatusCode(fasthttp.StatusInternalServerError)
				context.SetBodyString(err.Error())
			} else {
				context.SetStatusCode(fasthttp.StatusConflict)
				context.SetBody(existingForumJSON)
			}
		}
	}
}