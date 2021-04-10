package error

type Error struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	HttpCode int
}

func BadError(code, message string) *Error {
	return &Error{
		Code:     code,
		Message:  message,
		HttpCode: 400,
	}
}

func NotFoundError(code, message string) *Error {
	return &Error{
		Code:     code,
		Message:  message,
		HttpCode: 404,
	}
}

func InternalServerError(code, message string) *Error {
	return &Error{
		Code:     code,
		Message:  message,
		HttpCode: 500,
	}
}