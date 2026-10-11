package handler

import (
	"fmt"
	"net/http"

	exception "lazer/error"
	"lazer/laze"

	"github.com/gin-gonic/gin"
)

func (handler *Handler) error(err laze.Exception, c *gin.Context) {
	name := err.Name()
	message := err.Message()

	fmt.Println(err)

	status := http.StatusInternalServerError

	switch name {
	case exception.NOTFOUND:
		status = http.StatusNotFound
	case exception.BADREQUEST:
		status = http.StatusBadRequest
	case exception.UNPROCESSABLE:
		status = http.StatusUnprocessableEntity
	default:
		status = http.StatusInternalServerError
	}

	c.JSON(status, map[string]any{
		"message": message,
	})
}
