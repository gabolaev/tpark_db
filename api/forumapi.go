package api

import (
	"fmt"

	"github.com/gabolaev/tpark_db/database"

	"github.com/gabolaev/tpark_db/errors"
	"github.com/gabolaev/tpark_db/helpers"

	"github.com/gabolaev/tpark_db/models"
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
	result, err := helpers.CreateNewOrGetExistingForum(&forum)

	var responseStatus int
	switch err {
	case nil:
		responseStatus = fasthttp.StatusCreated
	case errors.ConflictError:
		responseStatus = fasthttp.StatusConflict
	case errors.NotFoundError:
		context.SetStatusCode(fasthttp.StatusNotFound)
		err := models.Error{Message: "Can't find user with that nickname"}
		errorJSON, _ := err.MarshalJSON()
		context.SetBody(errorJSON)
		return
	default:
		context.SetStatusCode(fasthttp.StatusInternalServerError)
		context.SetBodyString(err.Error())
		return
	}

	if responseJSON, err := result.MarshalJSON(); err != nil {
		context.SetStatusCode(fasthttp.StatusInternalServerError)
	} else {
		defer func() {
			database.Instance.Status.Forum++
		}()
		context.SetStatusCode(responseStatus)
		context.SetBody(responseJSON)
	}

}

func GetForumUsers(context *fasthttp.RequestCtx) {
	// todo
}

func GetForumDetails(context *fasthttp.RequestCtx) {

	context.SetContentType("application/json")
	slug := context.UserValue("slug").(string)
	result, err := helpers.GetForumDetailsBySlug(&slug)
	switch err {
	case nil:
		if user, err := result.MarshalJSON(); err != nil {
			context.SetStatusCode(fasthttp.StatusInternalServerError)
		} else {
			context.SetStatusCode(fasthttp.StatusOK)
			context.SetBody(user)
		}
	default:
		context.SetStatusCode(fasthttp.StatusNotFound)
		err := models.Error{
			Message: fmt.Sprintf("Can't find forum with slug: %s", slug)}
		errorJSON, _ := err.MarshalJSON()
		context.SetBody(errorJSON)
	}
}

func GetForumThreads(context *fasthttp.RequestCtx) {
	context.SetContentType("application/json")
	slug := context.UserValue("slug").(string)
	limit, desc, since := context.QueryArgs().Peek("limit"), context.QueryArgs().Peek("desc"), context.QueryArgs().Peek("since")
	result, err := helpers.GetThreadsByForumSlug(&slug, limit, desc, since)
	switch err {
	case nil:
		if thread, err := result.MarshalJSON(); err != nil {
			context.SetStatusCode(fasthttp.StatusInternalServerError)
			context.SetBodyString(err.Error())
		} else {
			context.SetStatusCode(fasthttp.StatusOK)
			context.SetBody(thread)
		}
	case errors.NotFoundError:
		context.SetStatusCode(fasthttp.StatusNotFound)
		err := models.Error{
			Message: fmt.Sprintf("Can't find forum with slug: %s", slug)}
		errorJSON, _ := err.MarshalJSON()
		context.SetBody(errorJSON)
	case errors.EmptySearchError:
		context.SetStatusCode(fasthttp.StatusOK)
		context.SetBody([]byte{91, 93})
	default:
		context.SetStatusCode(fasthttp.StatusInternalServerError)
		context.SetBodyString(err.Error())
	}
}
