package models

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

//easyjson:json
type Thread struct {
	ID      int
	Slug    string
	Author  string
	Created string
	Forum   string
	Message string
	Title   string
	Votes   int
}

//easyjson:json
type ThreadUpdate struct {
	Message string
	Title   string
}

//easyjson:json
type Threads []Thread

func GetThreadInfo(context *fasthttp.RequestCtx) {
	fmt.Println(context) // debug
}

func GetThreadPosts(context *fasthttp.RequestCtx) {
	fmt.Println(context) // debug
}

func UpdateThreadInfo(context *fasthttp.RequestCtx) {
	fmt.Println(context) // debug
}

func CreateThreadPosts(context *fasthttp.RequestCtx) {
	fmt.Println(context) // debug
}

func VoteThread(context *fasthttp.RequestCtx) {
	fmt.Println(context) // debug
}
