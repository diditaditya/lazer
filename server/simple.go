package server

import (
	"github.com/gin-gonic/gin"

	laze "lazer/laze"

	"lazer/server/handler"
)

func Start(app *laze.App) {
	router := gin.Default()

	handle := handler.New(app)

	router.GET("/", handle.Root)
	router.GET("/:name", handle.GetAll)
	router.GET("/:name/:pk", handle.GetByPk)

	router.POST("/:name", handle.Create)

	router.DELETE("/:name", handle.Delete)
	router.DELETE("/:name/:pk", handle.DeleteByPk)

	router.Run()
}