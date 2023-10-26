package errorlib

type CustomError struct {
	StatusCode int
	Message    string
}

func New(statusCode int, message string) *CustomError {
	return &CustomError{StatusCode: statusCode, Message: message}
}
