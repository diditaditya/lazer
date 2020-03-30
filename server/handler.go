package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	exception "lazer/error"
	"lazer/laze"
)

func errorHandler(err laze.Exception, c *gin.Context) {
	name := err.Name()
	message := err.Message()

	fmt.Println(err)

	status := http.StatusInternalServerError

	switch name {
	case exception.NOTFOUND:
		status = http.StatusNotFound
	default:
		status = http.StatusInternalServerError
	}

	c.JSON(status, map[string]interface{}{
			"message": message,
		})
}