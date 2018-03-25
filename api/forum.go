package api

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

//easyjson:json
type Forum struct {
	Slug    string
	Posts   int64
	Threads int
	Title   string
	Creator string
}

func GetForumUsers(context *fasthttp.RequestCtx) {
	fmt.Println(context)
}

func GetForumInfo(context *fasthttp.RequestCtx) {
	fmt.Println(context)
}

func GetForumThreads(context *fasthttp.RequestCtx) {
	fmt.Println(context)
}

func CreateForum(context *fasthttp.RequestCtx) {
	fmt.Println(context)
}
