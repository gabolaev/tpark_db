package api

import (
	"github.com/gabolaev/tpark_db/database"
	"github.com/valyala/fasthttp"
)

func Status(context *fasthttp.RequestCtx) {
	context.SetContentType("application/json")
	if result, err := database.Instance.Status.MarshalJSON(); err != nil {
		context.SetStatusCode(fasthttp.StatusInternalServerError)
	} else {
		context.SetStatusCode(fasthttp.StatusOK)
		context.SetBody(result)
	}
}

func Clear(context *fasthttp.RequestCtx) {
	context.SetContentType("application/json")
	if err := database.Instance.Clear(); err != nil {
		context.SetStatusCode(fasthttp.StatusInternalServerError)
	} else {
		context.SetStatusCode(fasthttp.StatusOK)
	}
}
