package models

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

//easyjson:json
type Post struct {
	ID      int64
	Author  string
	Created string
	Forum   string
	Edited  bool
	Message string
	Parent  int64
	Thread  int
}

//easyjson:json
type PostFull struct {
	Author *User
	Forum  *Forum
	Post   *Post
	Thread *Thread
}

//easyjson:json
type PostUpdate struct {
	Message string
}

func GetPostFull(context *fasthttp.RequestCtx) {
	fmt.Println(context) // debug
}

func UpdatePost(context *fasthttp.RequestCtx) {
	fmt.Println(context) // debug
}
