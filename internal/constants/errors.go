package constants

import (
	"errors"
	"net/http"
)

var (
	ErrCustomErr          = errors.New("custom error")
	ErrInvalidCredentials = errors.New("incorrect username or password")

	ErrorStatusCodeMaps = map[error]int{
		// Business logic error
		ErrCustomErr:          http.StatusBadRequest,
		ErrInvalidCredentials: http.StatusBadRequest,
	}
)
