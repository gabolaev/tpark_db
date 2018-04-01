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
	Instance.GET("/forum/:slug/details", api.GetForumInfo)
	// Instance.GET("/forum/:slug/threads", api.GetForumThreads)
	// Instance.POST("/post/:id/details", api.UpdatePost)

	// Instance.GET("/service/status", api.GetDBInfo)
	// Instance.POST("/service/clear", api.ClearDB)

	// Instance.GET("/thread/:slug_or_id/details", api.GetThreadInfo)
	// Instance.GET("/thread/:slug_or_id/posts", api.GetThreadPosts)
	// Instance.POST("/thread/:slug_or_id/details", api.UpdateThreadInfo)
	// Instance.POST("/thread/:slug_or_id/create", api.CreateThreadPosts)
	// Instance.POST("/thread/:slug_or_id/vote", api.VoteThread)

	Instance.GET("/user/:nickname/profile", api.GetUser)
	Instance.POST("/user/:nickname/create", api.CreateUser)
	Instance.POST("/user/:nickname/profile", api.UpdateUser)

}
