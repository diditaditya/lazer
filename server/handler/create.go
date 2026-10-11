package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"

	exception "lazer/error"
)

func (handler *Handler) Create(c *gin.Context) {
	tableName := c.Param("name")
	rawBody, err := c.GetRawData()

	if err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "your fault..",
			"error":   err,
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

	_, err = handler.app.Create(tableName, mapped)

	if err != nil {
		ex, ok := err.(*exception.Exception)
		if ok {
			handler.error(ex, c)
		} else {
			ex = exception.FromError(err, exception.INTERNALERROR)
			handler.error(ex, c)
		}
		return
	}

	c.JSON(201, map[string]interface{}{
		"message": "created..",
	})
}
