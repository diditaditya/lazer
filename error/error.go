package error

import (
	"github.com/pkg/errors"
	"fmt"
)

type Exception struct {
	message string
	name string
	trace string
}

func (ex *Exception) Error() string {
	return ex.trace
}

func New(name string, message string) Exception {
	cause := errors.New(message)
	err := errors.WithStack(cause)
	trace := fmt.Sprintf("%+v\n", err)
	exception := Exception{
		message: message,
		name: name,
		trace: trace,
	}
	return exception
}

func FromError(cause error, name string) Exception {
	err:= errors.WithStack(cause)
	trace := fmt.Sprintf("%+v\n", err)
	exception := Exception{
		message: fmt.Sprintf("%s", cause),
		name: name,
		trace: trace,
	}
	return exception
}

func (ex *Exception) Message() string {
	return ex.message
}

func (ex *Exception) Name() string {
	return ex.name
}