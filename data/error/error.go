package error

type Exception struct {
	message string
	code string
}

func New(code string, message string) *Exception {
	exception := Exception{
		code: code,
		message: message,
	}
	return &exception
}

func (ex *Exception) Error() string {
	return ex.message
}

func (ex *Exception) Message() string {
	return ex.message
}

func (ex *Exception) Code() string {
	return ex.code
}