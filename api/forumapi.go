package api

import (
	"fmt"

	"github.com/gabolaev/tpark_db/helpers"

	"github.com/gabolaev/tpark_db/models"
	"github.com/mailru/easyjson"

	"github.com/valyala/fasthttp"
)

func CreateForum(context *fasthttp.RequestCtx) {
	var forum models.Forum
	body := context.PostBody()
	if err := easyjson.Unmarshal(body, &forum); err != nil {
		context.SetStatusCode(fasthttp.StatusBadRequest)
		context.WriteString(err.Error())
		return
	}
	result, created, err := helpers.CreateNewOrGetExistingForum(&forum)
	if err != nil {
		context.SetStatusCode(fasthttp.StatusInternalServerError)
		context.SetBodyString(err.Error())
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
				Message: "Can't find user with that nickname"})
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

func GetForumUsers(context *fasthttp.RequestCtx) {
	fmt.Println(context)
}

func GetForumInfo(context *fasthttp.RequestCtx) {
	fmt.Println(context)
}

func GetForumThreads(context *fasthttp.RequestCtx) {
	fmt.Println(context)
}
