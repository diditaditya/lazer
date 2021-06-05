// The handler package takes all the http requests through the server
package handler

import (
	"github.com/gin-gonic/gin"

	laze "lazer/laze"
)

// The Handler struct
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

// The function will return the instance of Handler struct
func New(app *laze.App) Handler {
	return Handler{app: app}
}
