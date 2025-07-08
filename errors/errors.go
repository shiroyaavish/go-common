package errors

import "fmt"

// Error represents an error with a status code and a message.
type Error struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

// NewError creates a new Error object with the given statusCode and message. It
func NewError(statusCode int, message string) *Error {
	return &Error{
		StatusCode: statusCode,
		Message:    message,
	}
}

func NewErrorWrapper(err error, statusCode int, message string) *Error {
	return &Error{
		StatusCode: statusCode,
		Message:    fmt.Sprintf("error: %s\nmsg: %s", err.Error(), message),
	}
}

// Error returns a string representation of the Error object. It formats
// the status code and message as "Status Code: %d, Message: %s",
// replacing the placeholders with the actual values of StatusCode
// and Message respectively.
func (a Error) Error() string {
	return fmt.Sprintf("Status Code: %d, Message: %s", a.StatusCode, a.Message)
}

var (
	// ErrIsNil used when data is nil contains Error.StatusCode: 404 and message
	ErrIsNil = NewError(404, "data is nil")
)
