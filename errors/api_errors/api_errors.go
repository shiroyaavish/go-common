package api_errors

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type Error struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func NewError(statusCode int, message string) *Error {
	return &Error{
		StatusCode: statusCode,
		Message:    message,
	}
}

func (a Error) Error() string {
	return fmt.Sprintf("Status Code: %d, Message: %s", a.StatusCode, a.Message)
}

var (
	ErrInvalidParams      = NewError(fiber.StatusBadRequest, "invalid params")
	ErrNotFound           = NewError(fiber.StatusNotFound, "not found")
	ErrUnauthorized       = NewError(fiber.StatusUnauthorized, "unauthorized")
	ErrForbidden          = NewError(fiber.StatusForbidden, "forbidden")
	ErrConflict           = NewError(fiber.StatusConflict, "conflict")
	ErrSomethingWentWrong = NewError(fiber.StatusInternalServerError, "something went wrong")
)
