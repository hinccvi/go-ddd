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

	ErrorStatusCodeMaps = map[error]int{
		// Business logic error
		ErrMaxAttempt: http.StatusBadRequest,

		ErrInvalidCredentials:  http.StatusBadRequest,
		ErrConditionNotFulfil:  http.StatusBadRequest,
		ErrInvalidRefreshToken: http.StatusBadRequest,
		ErrInvalidJwt:          http.StatusBadRequest,
	}
)
