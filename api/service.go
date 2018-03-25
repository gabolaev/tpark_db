package api

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func ClearDB(context *fasthttp.RequestCtx) {
	fmt.Println(context)
}

func GetDBInfo(context *fasthttp.RequestCtx) {
	fmt.Println(context)
}
