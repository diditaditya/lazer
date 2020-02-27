package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	laze "lazer/laze"
)

type Handler struct {
	app *laze.App
}

func (handler *Handler) defaultHandler(c *gin.Context) {
	tableName := c.Param("name")

	data, err := handler.app.FindAll(tableName)

	if err != nil {
		fmt.Println(err)
		c.JSON(404, map[string]interface{}{
			"message": "oops..",
			"error": fmt.Sprintf("%s", err),
		})
	} else {
		resp := map[string]interface{}{
			"message": "yo",
			"data": data,
		}
	
		c.JSON(200, resp)
	}

	
}

func Start(app *laze.App) {
	r := gin.Default()

	handler := Handler{app: app}

	r.GET("/:name", handler.defaultHandler)

	r.Run()
}