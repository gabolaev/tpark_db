package models

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func ClearDB(context *fasthttp.RequestCtx) {
	fmt.Println(context) // debug
}

func GetDBInfo(context *fasthttp.RequestCtx) {
	fmt.Println(context) // debug
}
