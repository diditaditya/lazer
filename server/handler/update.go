package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"

	exception "lazer/error"
)

func (handler *Handler) UpdateByPk(c *gin.Context) {
	tableName := c.Param("name")
	pk := c.Param("pk")
	rawBody, err := c.GetRawData()

	if err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "your fault..",
			"error": err,
		})
		return
	}

	mapped := make(map[string]interface{})

	err = json.Unmarshal(rawBody, &mapped)

	if err != nil {
		ex := exception.FromError(err, exception.INTERNALERROR)
		handler.error(ex, c)
		return
	}

	updateErr := handler.app.UpdateByPk(tableName, pk, mapped)

	if updateErr != nil {
		handler.error(updateErr, c)
	} else {
		resp := map[string]interface{}{
			"message": "updated",
		}
	
		c.JSON(200, resp)
	}
}