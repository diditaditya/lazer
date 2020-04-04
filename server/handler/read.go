package handler

import (
	"github.com/gin-gonic/gin"
)

func (handler *Handler) GetByPk(c *gin.Context) {
	tableName := c.Param("name")
	pk := c.Param("pk")

	data, err := handler.app.FindByPk(tableName, pk)

	if err != nil {
		handler.error(err, c)
	} else {
		resp := map[string]interface{}{
			"message": "yo",
			"data":    data,
		}

		c.JSON(200, resp)
	}
}

func (handler *Handler) GetAll(c *gin.Context) {
	tableName := c.Param("name")

	query := c.Request.URL.Query()

	data, err := handler.app.FindAll(tableName, query)
	if err != nil {
		handler.error(err, c)
	} else {
		resp := map[string]interface{}{
			"message": "yo",
			"data":    data,
		}

		c.JSON(200, resp)
	}
}
