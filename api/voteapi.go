package api

import (
	"github.com/gabolaev/tpark_db/errors"
	"github.com/gabolaev/tpark_db/helpers"
	"github.com/gabolaev/tpark_db/models"
	"github.com/valyala/fasthttp"
)

func VoteThread(context *fasthttp.RequestCtx) {
	context.SetContentType("application/json")
	slugOrID := context.UserValue("slug_or_id").(string)
	var vote models.Vote
	if err := vote.UnmarshalJSON(context.PostBody()); err != nil {
		context.SetStatusCode(fasthttp.StatusBadRequest)
		context.WriteString(err.Error())
		return
	}
	thread, err := helpers.VoteThread(&slugOrID, &vote)
	switch err {
	case nil:
		context.SetStatusCode(fasthttp.StatusOK)
		if updatedThread, err := thread.MarshalJSON(); err != nil {
			context.SetStatusCode(fasthttp.StatusInternalServerError)
			context.SetBodyString(err.Error())
		} else {
			context.SetStatusCode(fasthttp.StatusOK)
			context.SetBody(updatedThread)
		}
	case errors.NotFoundError:
		context.SetStatusCode(fasthttp.StatusNotFound)
		err := models.Error{
			Message: "Can't find user or thread"}
		errorJSON, _ := err.MarshalJSON()
		context.SetBody(errorJSON)
	case errors.NothingChangedError:
		context.SetStatusCode(fasthttp.StatusNotModified)
	}

}
