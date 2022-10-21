package constants

import (
	"errors"
	"net/http"
)

const (
	MsgSystemError    = "system error"
	MsgBadRequest     = "invalid input"
	MsgRequestTimeout = "request timeout"
)

var (
	ErrMaxAttempt          = errors.New("max attempt reached")
	ErrInvalidCredentials  = errors.New("incorrect username or password")
	ErrConditionNotFulfil  = errors.New("condition not fulfil")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrInvalidJwt          = errors.New("invalid token")
	ErrResourceNotFound    = errors.New("resource not found")
	ErrCRUD                = errors.New("error crud")
	ErrSystemError         = errors.New("system error")

	ErrorStatusCodeMaps = map[error]int{
		ErrInvalidCredentials:  http.StatusBadRequest,
		ErrConditionNotFulfil:  http.StatusBadRequest,
		ErrResourceNotFound:    http.StatusBadRequest,
		ErrInvalidRefreshToken: http.StatusForbidden,
		ErrInvalidJwt:          http.StatusForbidden,
		ErrCRUD:                http.StatusInternalServerError,
		ErrSystemError:         http.StatusInternalServerError,

		// Business logic error
		ErrMaxAttempt: http.StatusBadRequest,
	}
)
