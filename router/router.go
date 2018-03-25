package router

import (
	"github.com/buaazp/fasthttprouter"
)

func Router() *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.POST("/user/:nickname/create", user.Create)

	return router
}
