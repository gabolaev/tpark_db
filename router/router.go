package router

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/gabolaev/tpark_db/api"
)

// Instance of router
var Instance = fasthttprouter.New()

func init() {

	Instance.POST("/forum/*catch-all-param", api.CreateThreadOrForum)
	Instance.GET("/forum/:slug/users", api.GetForumUsers)
	Instance.GET("/forum/:slug/details", api.GetForumDetails)
	Instance.GET("/forum/:slug/threads", api.GetForumThreads)
	Instance.GET("/post/:id/details", api.GetPostDetails)
	Instance.POST("/post/:id/details", api.UpdatePost)
	Instance.GET("/service/status", api.Status)
	Instance.POST("/service/clear", api.Clear)

	Instance.GET("/thread/:slug_or_id/details", api.GetThreadDetails)
	// Instance.GET("/thread/:slug_or_id/posts", api.GetThreadPosts)
	Instance.POST("/thread/:slug_or_id/details", api.UpdateThreadDetails)
	Instance.POST("/thread/:slug_or_id/create", api.CreateThreadPosts)
	Instance.POST("/thread/:slug_or_id/vote", api.VoteThread)

	Instance.GET("/user/:nickname/profile", api.GetUser)
	Instance.POST("/user/:nickname/create", api.CreateUser)
	Instance.POST("/user/:nickname/profile", api.UpdateUser)

}
