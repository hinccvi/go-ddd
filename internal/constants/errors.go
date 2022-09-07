package constants

import (
	"errors"
)

var (
	ErrCustomErr = errors.New("custom error")

	ErrorStatusCodeMaps = map[error]int{
		// Business logic error
		ErrCustomErr: 123,
	}
)
