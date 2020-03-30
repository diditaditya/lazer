package server

import (
	"github.com/gin-gonic/gin"
)

func (handler *Handler) delete(c *gin.Context) {
	tableName := c.Param("name")

	query := c.Request.URL.Query()

	err := handler.app.Delete(tableName, query)
	if err != nil {
		errorHandler(err, c)
	} else {
		resp := map[string]interface{}{
			"message": "deleted",
		}

		c.JSON(200, resp)
	}
}