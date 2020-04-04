package handler

import (
	"github.com/gin-gonic/gin"

	laze "lazer/laze"
)

type Handler struct {
	app *laze.App
}

func (handler *Handler) Root(c *gin.Context) {
	data := handler.app.GetAllTables()

	resp := map[string]interface{}{
		"message": "welcome",
		"data":    data,
	}

	c.JSON(200, resp)
}

func New(app *laze.App) Handler {
	return Handler{app: app}
}
