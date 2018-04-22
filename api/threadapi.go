package api

import (
	"fmt"

	"github.com/gabolaev/tpark_db/database"
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
	if err := thread.UnmarshalJSON(context.PostBody()); err != nil {
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
		defer func() {
			database.Instance.Status.Thread++
		}()
		context.SetStatusCode(responseStatus)
		context.SetBody(existingForumJSON)
	}
}

func GetThreadDetails(context *fasthttp.RequestCtx) {
	context.SetContentType("application/json")
	slugOrID := context.UserValue("slug_or_id").(string)
	var err error
	var thread *models.Thread
	thread, err = helpers.GetThreadDetailsBySlugOrID(&slugOrID)
	switch err {
	case nil:
		if existingForumJSON, err := thread.MarshalJSON(); err != nil {
			context.SetStatusCode(fasthttp.StatusInternalServerError)
			context.SetBodyString(err.Error())
		} else {
			context.SetStatusCode(fasthttp.StatusOK)
			context.SetBody(existingForumJSON)
		}
	case errors.NotFoundError:
		err := models.Error{Message: "Can't find user or forum"}
		errorJSON, _ := err.MarshalJSON()
		context.SetStatusCode(fasthttp.StatusNotFound)
		context.SetBody(errorJSON)
	}

}

func UpdateThreadDetails(context *fasthttp.RequestCtx) {
	context.SetContentType("application/json")
	slugOrID := context.UserValue("slug_or_id").(string)
	body := context.PostBody()
	threadUpdate := models.ThreadUpdate{}
	if err := threadUpdate.UnmarshalJSON(body); err != nil {
		if err.Error() != "EOF" {
			context.SetStatusCode(fasthttp.StatusBadRequest)
			context.WriteString(err.Error())
			return
		}
	}

	thread, err := helpers.UpdateThreadDetails(&slugOrID, &threadUpdate)
	switch err {
	case nil:
		if existingForumJSON, err := thread.MarshalJSON(); err != nil {
			context.SetStatusCode(fasthttp.StatusInternalServerError)
			context.SetBodyString(err.Error())
		} else {
			context.SetStatusCode(fasthttp.StatusOK)
			context.SetBody(existingForumJSON)
		}
	case errors.NotFoundError:
		err := models.Error{Message: "Can't find thread"}
		errorJSON, _ := err.MarshalJSON()
		context.SetStatusCode(fasthttp.StatusNotFound)
		context.SetBody(errorJSON)
	}
	return

}

func GetThreadPosts(context *fasthttp.RequestCtx) {
	context.SetContentType("application/json")
	slugOrID := context.UserValue("slug_or_id").(string)
	limit, desc, since :=
		context.QueryArgs().Peek("limit"),
		context.QueryArgs().Peek("desc"),
		context.QueryArgs().Peek("since")
	var result *models.Posts
	var err error
	switch string(context.QueryArgs().Peek("sort")) {
	case "flat":
		result, err = helpers.GetThreadPostsFlat(&slugOrID, limit, since, desc)
	case "tree":
		fmt.Println("implement")
		return
	case "parent_tree":
		fmt.Println("implement")
		return
	}
	switch err {
	case nil:
		if posts, err := result.MarshalJSON(); err != nil {
			context.SetStatusCode(fasthttp.StatusInternalServerError)
			context.SetBodyString(err.Error())
		} else {
			context.SetStatusCode(fasthttp.StatusOK)
			context.SetBody(posts)
		}
	case errors.NotFoundError:
		err := models.Error{Message: "Can't find thread with that slug or id"}
		errorJSON, _ := err.MarshalJSON()
		context.SetStatusCode(fasthttp.StatusNotFound)
		context.SetBody(errorJSON)
	}
}
