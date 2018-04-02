package api

import (
	"fmt"

	"github.com/gabolaev/tpark_db/helpers"
	"github.com/gabolaev/tpark_db/models"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func CreateThreadPosts(context *fasthttp.RequestCtx) {
	context.SetContentType("application/json")
	slugOrID := context.UserValue("slug_or_id").(string)
	var posts models.Posts
	body := context.PostBody()
	if err := easyjson.Unmarshal(body, &posts); err != nil {
		context.SetStatusCode(fasthttp.StatusBadRequest)
		context.WriteString(err.Error())
		return
	}

	result, code, err := helpers.CreatePostsByThreadSlugOrID(&posts, &slugOrID)
	switch code {
	case 201:
		if posts, err := easyjson.Marshal(*result); err != nil {
			context.SetStatusCode(fasthttp.StatusInternalServerError)
		} else {
			context.SetStatusCode(code)
			context.SetBody(posts)
		}
	case 404:
		context.SetStatusCode(code)
		errorJSON, _ := easyjson.Marshal(models.Error{
			Message: "Forum not found"})
		context.SetBody(errorJSON)
	case 409:
		context.SetStatusCode(code)
		errorJSON, _ := easyjson.Marshal(models.Error{
			Message: "Conflict parents"})
		context.SetBody(errorJSON)
	default:
		fmt.Println(err)
	}
	// if err != nil {
	// 	context.SetStatusCode(fasthttp.StatusInternalServerError)
	// 	context.SetBodyString(err.Error())
	// 	return
	// }
	// if created {
	// 	if createdForumJSON, err := easyjson.Marshal(result); err != nil {
	// 		context.SetStatusCode(fasthttp.StatusInternalServerError)
	// 		context.SetBodyString(err.Error())
	// 	} else {
	// 		context.SetStatusCode(fasthttp.StatusCreated)
	// 		context.SetBody(createdForumJSON)
	// 	}
	// } else {
	// 	if result == nil {
	// 		errorJSON, _ := easyjson.Marshal(models.Error{
	// 			Message: "Can't find user or forum"})
	// 		context.SetStatusCode(fasthttp.StatusNotFound)
	// 		context.SetBody(errorJSON)
	// 	} else {
	// 		if existingForumJSON, err := easyjson.Marshal(result); err != nil {
	// 			context.SetStatusCode(fasthttp.StatusInternalServerError)
	// 			context.SetBodyString(err.Error())
	// 		} else {
	// 			context.SetStatusCode(fasthttp.StatusConflict)
	// 			context.SetBody(existingForumJSON)
	// 		}
	// 	}
	// }

}
