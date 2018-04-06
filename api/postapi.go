package api

import (
	"github.com/gabolaev/tpark_db/errors"
	"github.com/gabolaev/tpark_db/helpers"
	"github.com/gabolaev/tpark_db/models"

	"github.com/valyala/fasthttp"
)

func CreateThreadPosts(context *fasthttp.RequestCtx) {
	context.SetContentType("application/json")
	var posts models.Posts
	if err := posts.UnmarshalJSON(context.PostBody()); err != nil {
		context.SetStatusCode(fasthttp.StatusBadRequest)
		context.WriteString(err.Error())
		return
	}

	slugOrID := context.UserValue("slug_or_id").(string)
	result, err := helpers.CreatePostsByThreadSlugOrID(&posts, &slugOrID)

	var errorObj models.Error
	var responseStatus int
	switch err {
	case nil:
		if posts, err := result.MarshalJSON(); err != nil {
			context.SetStatusCode(fasthttp.StatusInternalServerError)
			context.SetBodyString(err.Error())
		} else {
			context.SetStatusCode(fasthttp.StatusCreated)
			context.SetBody(posts)
		}
		return
	case errors.NotFoundError:
		responseStatus = fasthttp.StatusNotFound
		errorObj = models.Error{Message: "Forum not found"}
	case errors.ConflictError:
		responseStatus = fasthttp.StatusConflict
		errorObj = models.Error{Message: "Conflict or wrong parents"}
	default:
		context.SetStatusCode(fasthttp.StatusInternalServerError)
		context.SetBodyString(err.Error())
		return
	}
	context.SetStatusCode(responseStatus)
	errorJSON, _ := errorObj.MarshalJSON()
	context.SetBody(errorJSON)
}
