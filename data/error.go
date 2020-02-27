package data

import (
	"errors"
	"fmt"
)

type Exception struct {
	message string
	code string
}

func (ex *Exception) Error() (string, string) {
	return ex.code, ex.message
}

func (ex *Exception) Message() string {
	return ex.message
}

func (ex *Exception) Code() string {
	return ex.code
}