package handler

import (
	"github.com/gin-gonic/gin"
)

func (handler *Handler) DeleteByPk(c *gin.Context) {
	tableName := c.Param("name")
	pk := c.Param("pk")

	err := handler.app.DeleteByPk(tableName, pk)

	if err != nil {
		handler.error(err, c)
	} else {
		resp := map[string]interface{}{
			"message": "deleted",
		}
	
		c.JSON(200, resp)
	}
}

func (handler *Handler) Delete(c *gin.Context) {
	tableName := c.Param("name")

	query := c.Request.URL.Query()

	err := handler.app.Delete(tableName, query)
	if err != nil {
		handler.error(err, c)
	} else {
		resp := map[string]interface{}{
			"message": "deleted",
		}

		c.JSON(200, resp)
	}
}