package handler

import (
	"github.com/gin-gonic/gin"
)

func (handler *Handler) GetAssociations(c *gin.Context) {
	tableName := c.Param("name")

	data := handler.app.GetAssociations(tableName)
	
	resp := map[string]interface{}{
		"message": "it's complicated",
		"data":    data,
	}

	c.JSON(200, resp)
}
