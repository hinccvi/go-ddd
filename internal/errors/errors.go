package errors

import (
	"database/sql"
	"errors"
	"net/http"
)

var (
	ErrMaxAttempt          = errors.New("max attempt reached")
	ErrInvalidCredentials  = errors.New("incorrect username or password")
	ErrConditionNotFulfil  = errors.New("condition not fulfil")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrInvalidJwt          = errors.New("invalid token")
	ErrResourceNotFound    = errors.New("resource fail: not found")
	ErrNoRows              = sql.ErrNoRows
	ErrSystemError         = errors.New("system error")
)

func GetStatusCodeMap() map[error]int {
	return map[error]int{
		ErrInvalidCredentials: http.StatusBadRequest,
		ErrConditionNotFulfil: http.StatusBadRequest,
		ErrResourceNotFound:   http.StatusBadRequest,
		ErrNoRows:             http.StatusBadRequest,

		ErrInvalidRefreshToken: http.StatusForbidden,
		ErrInvalidJwt:          http.StatusForbidden,

		ErrSystemError: http.StatusInternalServerError,

		// Business logic error
		ErrMaxAttempt: http.StatusBadRequest,
	}
}
