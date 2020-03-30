package server

import (
	"github.com/gin-gonic/gin"

	laze "lazer/laze"
)

type Handler struct {
	app *laze.App
}

func (handler *Handler) root(c *gin.Context) {
	data := handler.app.GetAllTables()

	resp := map[string]interface{}{
		"message": "welcome",
		"data": data,
	}

	c.JSON(200, resp)
}

func Start(app *laze.App) {
	router := gin.Default()

	handler := Handler{app: app}

	router.GET("/", handler.root)
	router.GET("/:name", handler.getAll)
	router.GET("/:name/:pk", handler.getByPk)

	router.POST("/:name", handler.create)

	router.DELETE("/:name", handler.delete)

	router.Run()
}