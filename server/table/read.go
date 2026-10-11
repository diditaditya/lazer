package table

import (
	"github.com/gin-gonic/gin"
)

func (handler *Handler) GetAll(c *gin.Context) {
	data := handler.app.GetAllTables()

	resp := map[string]interface{}{
		"message": "welcome",
		"data":    data,
	}

	c.JSON(200, resp)
}