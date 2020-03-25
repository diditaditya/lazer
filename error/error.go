package error

import (
	"github.com/pkg/errors"
	"fmt"
)

struct Exception {
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
		trace: trace
	}
	return exception
}