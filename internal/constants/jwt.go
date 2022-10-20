package constants

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type (
	JWTData struct {
		UserName string
	}

	JWTCustomClaims struct {
		JWTData
		jwt.RegisteredClaims
	}
)

const (
	JWTRemainingTime = 60 * time.Second
	JWTpart          = 2
)
