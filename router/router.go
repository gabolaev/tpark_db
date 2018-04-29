package router

import (
	"go.uber.org/zap"

	"github.com/buaazp/fasthttprouter"
	"github.com/gabolaev/tpark_db/api"
	"github.com/gabolaev/tpark_db/logger"
	"github.com/valyala/fasthttp"
)

// Instance of router
var Instance = fasthttprouter.New()

func LogHandledRequests(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(context *fasthttp.RequestCtx) {
		logger.Instance.Info("[REQ]",
			zap.String("FROM", context.RemoteAddr().String()),
			zap.ByteString("URI", context.Path()),
			zap.ByteString("METHOD", context.Method()),
		)
		logger.Instance.Sync()
		handler(context)
	})
}

func init() {
	Instance.POST("/api/forum/*catch-all-param", api.CreateThreadOrForum)
	Instance.GET("/api/forum/:slug/users", api.GetForumUsers)
	Instance.GET("/api/forum/:slug/details", api.GetForumDetails)
	Instance.GET("/api/forum/:slug/threads", api.GetForumThreads)
	Instance.GET("/api/post/:id/details", api.GetPostDetails)
	Instance.POST("/api/post/:id/details", api.UpdatePost)
	Instance.GET("/api/service/status", api.Status)
	Instance.POST("/api/service/clear", api.Clear)

	Instance.GET("/api/thread/:slug_or_id/details", api.GetThreadDetails)
	Instance.GET("/api/thread/:slug_or_id/posts", api.GetThreadPosts)
	Instance.POST("/api/thread/:slug_or_id/details", api.UpdateThreadDetails)
	Instance.POST("/api/thread/:slug_or_id/create", api.CreateThreadPosts)
	Instance.POST("/api/thread/:slug_or_id/vote", api.VoteThread)

	Instance.GET("/api/user/:nickname/profile", api.GetUser)
	Instance.POST("/api/user/:nickname/create", api.CreateUser)
	Instance.POST("/api/user/:nickname/profile", api.UpdateUser)

}
