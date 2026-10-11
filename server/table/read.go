package table

import (
	"github.com/gin-gonic/gin"
)

func (handler *Handler) GetAll(c *gin.Context) {
	data := handler.app.GetAllTables()

	resp := map[string]any{
		"message": "welcome",
		"data":    data,
	}

	c.JSON(200, resp)
}
