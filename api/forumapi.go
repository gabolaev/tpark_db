package api

import (
	"fmt"

	"github.com/gabolaev/tpark_db/helpers"

	"github.com/gabolaev/tpark_db/models"
	"github.com/mailru/easyjson"

	"github.com/valyala/fasthttp"
)

func CreateForum(context *fasthttp.RequestCtx) {
	context.SetContentType("application/json")
	var forum models.Forum
	body := context.PostBody()

	if err := forum.UnmarshalJSON(body); err != nil {
		context.SetStatusCode(fasthttp.StatusBadRequest)
		context.WriteString(err.Error())
		return
	}
	result, created, err := helpers.CreateNewOrGetExistingForum(&forum)
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
	fmt.Println(context) // debug
}

func GetForumInfo(context *fasthttp.RequestCtx) {
	fmt.Println(context) // debug
	context.SetContentType("application/json")
	slug := context.UserValue("slug").(string)
	result, err := helpers.GetForumInfoBySlug(&slug)
	if err != nil {
		context.SetStatusCode(fasthttp.StatusNotFound)
		errorJSON, _ := easyjson.Marshal(models.Error{
			Message: fmt.Sprintf("Can't find forum with slug: %s", slug)})
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

func GetForumThreads(context *fasthttp.RequestCtx) {
	context.SetContentType("application/json")
	slug := context.UserValue("slug").(string)
	limit, desc, since := context.QueryArgs().Peek("limit"), context.QueryArgs().Peek("desc"), context.QueryArgs().Peek("since")
	result, emptySearch, err := helpers.GetThreadsByForumSlug(&slug, limit, desc, since)

	switch {
	case err != nil:
		context.SetStatusCode(fasthttp.StatusInternalServerError)
		context.SetBodyString(err.Error())
	case result == nil:
		if emptySearch {
			context.SetStatusCode(fasthttp.StatusOK)
			context.SetBody([]byte{91, 93})
		} else {
			context.SetStatusCode(fasthttp.StatusNotFound)
			errorJSON, _ := easyjson.Marshal(models.Error{
				Message: fmt.Sprintf("Can't find forum with slug: %s", slug)})
			context.SetBody(errorJSON)
		}
	default:

		if thread, err := result.MarshalJSON(); err != nil {
			context.SetStatusCode(fasthttp.StatusInternalServerError)
			context.SetBodyString(err.Error())
		} else {
			context.SetStatusCode(fasthttp.StatusOK)
			context.SetBody(thread)
		}
	}

}
