package table

import (
	"github.com/gin-gonic/gin"
)

func (handler *Handler) Create(c *gin.Context) {
	c.JSON(201, map[string]interface{}{
		"message": "table created",
	})
}
