package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	exception "lazer/error"
)

func errorHandler(err *exception.Exception, c *gin.Context) {
	name := err.Name()
	message := err.Message()

	fmt.Println(err)

	switch name {
	case exception.NOTFOUND:
		c.JSON(404, map[string]interface{}{
			"message": message,
		})
	default:
		c.JSON(500, map[string]interface{}{
			"message": "oops..",
			"error": fmt.Sprintf("%s", err),
		})
	}
}