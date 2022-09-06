package errors

import (
	"errors"
	"net/http"
)

// ErrorResponse is the response that represents an error.
type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func New(msg string) error {
	return errors.New(msg)
}

// Error is required by the error interface.
func (e ErrorResponse) Error() string {
	return e.Message
}

// InternalServerError creates a new error response representing an internal server error (HTTP 500)
func InternalServerError() ErrorResponse {
	return ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: "We encountered an error while processing your request.",
	}
}

// NotFound creates a new error response representing a resource-not-found error (HTTP 404)
func NotFound(msg string) ErrorResponse {
	if msg == "" {
		msg = "The requested resource was not found."
	}
	return ErrorResponse{
		Code:    http.StatusNotFound,
		Message: msg,
	}
}

// Unauthorized creates a new error response representing an authentication/authorization failure (HTTP 401)
func Unauthorized(msg string) ErrorResponse {
	if msg == "" {
		msg = "You are not authenticated to perform the requested action."
	}
	return ErrorResponse{
		Code:    http.StatusUnauthorized,
		Message: msg,
	}
}

// Forbidden creates a new error response representing an authorization failure (HTTP 403)
func Forbidden(msg string) ErrorResponse {
	if msg == "" {
		msg = "You are not authorized to perform the requested action."
	}
	return ErrorResponse{
		Code:    http.StatusForbidden,
		Message: msg,
	}
}

// BadRequest creates a new error response representing a bad request (HTTP 400)
func BadRequest(msg string) ErrorResponse {
	if msg == "" {
		msg = "Your request is in a bad format."
	}
	return ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: msg,
	}
}

type invalidField struct {
	Error string `json:"error"`
}

// InvalidInput creates a new error response representing a data validation error (HTTP 400).
func InvalidInput(status int, errs string) ErrorResponse {
	field := invalidField{
		Error: errs,
	}

	return ErrorResponse{
		Code:    status,
		Message: "There is some problem with the data you submitted.",
		Data:    field,
	}
}
