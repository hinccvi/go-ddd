package constants

import (
	"errors"
	"net/http"
)

const (
	MsgSystemError    = "system error"
	MsgBadRequest     = "invalid input"
	MsgRequestTimeout = "Request timeout"
)

var (
	ErrCustomErr           = errors.New("custom error")
	ErrInvalidCredentials  = errors.New("incorrect username or password")
	ErrConditionNotFulfil  = errors.New("condition not fulfil")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrInvalidJwt          = errors.New("invalid token")
	ErrMaxAttempt          = errors.New("max attempt reached")

	ErrorStatusCodeMaps = map[error]int{
		// Business logic error
		ErrCustomErr:           http.StatusBadRequest,
		ErrInvalidCredentials:  http.StatusBadRequest,
		ErrConditionNotFulfil:  http.StatusBadRequest,
		ErrInvalidRefreshToken: http.StatusBadRequest,
		ErrInvalidJwt:          http.StatusBadRequest,
		ErrMaxAttempt:          http.StatusBadRequest,
	}
)
