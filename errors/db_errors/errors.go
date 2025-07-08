package db_errors

import "github.com/shiroyaavish/go-common/errors"

var (
	// ErrAlreadyExists is a variable of type *Error that represents an error indicating that data already exists.
	// It is created using the errors.NewError function, with a status code of 403 and a message of "data already exists".
	// The *Error struct has two fields: StatusCode (int) and Message (string).
	// The Error() method of *Error returns a formatted string containing the status code and message.
	// Use ErrAlreadyExists to handle cases where duplicate data is detected.
	ErrAlreadyExists = errors.NewError(403, "data already exists")
)
