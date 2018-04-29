package api

import (
	"strings"

	"github.com/gabolaev/tpark_db/database"
	"github.com/gabolaev/tpark_db/errors"
	"github.com/gabolaev/tpark_db/helpers"
	"github.com/gabolaev/tpark_db/logger"
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
			defer func() {
				database.Instance.Status.Post++
			}()
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

func GetPostDetails(context *fasthttp.RequestCtx) {
	context.SetContentType("application/json")
	id := context.UserValue("id").(string)
	related := string(context.QueryArgs().Peek("related"))
	var result *models.PostFull
	var err error
	fullParams := []string{"post"}
	if len(related) != 0 {
		fullParams = append(fullParams, strings.Split(related, ",")...)
	}
	result, err = helpers.GetPostFullDetails(&id, fullParams)
	switch err {
	case nil:
		if posts, err := result.MarshalJSON(); err != nil {
			context.SetStatusCode(fasthttp.StatusInternalServerError)
			errStr := err.Error()
			logger.Instance.Error(errStr)
			context.SetBodyString(errStr)
		} else {
			context.SetStatusCode(fasthttp.StatusOK)
			context.SetBody(posts)
		}
	case errors.NotFoundError:
		err := models.Error{Message: "Can't find post"}
		errorJSON, _ := err.MarshalJSON()
		context.SetStatusCode(fasthttp.StatusNotFound)
		context.SetBody(errorJSON)
	}
}

func UpdatePost(context *fasthttp.RequestCtx) {
	context.SetContentType("application/json")
	id := context.UserValue("id").(string)
	var postUpdate models.PostUpdate
	if err := postUpdate.UnmarshalJSON(context.PostBody()); err != nil {
		context.SetStatusCode(fasthttp.StatusBadRequest)
		context.WriteString(err.Error())
		return
	}
	res, err := helpers.UpdatePostDetails(&id, &postUpdate)
	switch err {
	case nil:
		if post, err := res.MarshalJSON(); err != nil {
			context.SetStatusCode(fasthttp.StatusInternalServerError)
			errStr := err.Error()
			logger.Instance.Error(errStr)
			context.SetBodyString(errStr)
		} else {
			context.SetStatusCode(fasthttp.StatusOK)
			context.SetBody(post)
		}
	case errors.NotFoundError:
		err := models.Error{Message: "Can't find post"}
		errorJSON, _ := err.MarshalJSON()
		context.SetStatusCode(fasthttp.StatusNotFound)
		context.SetBody(errorJSON)
	}
}
