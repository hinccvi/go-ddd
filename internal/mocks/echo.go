package mocks

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type (
	jwtCustomClaims struct {
		UserName string `json:"username"`
		jwt.RegisteredClaims
	}
)

// // Router creates a echo router for testing APIs.
//
//	func Router(logger log.Logger) *echo.Echo {
//		e := echo.New()
//
//		e.HTTPErrorHandler = m.NewHTTPErrorHandler(errs.GetStatusCodeMap()).Handler(logger)
//
//		e.Validator = &m.CustomValidator{Validator: validator.New()}
//
//		return e
//	}
//
// // AuthHeader returns an HTTP header that can pass the authentication check by MockAuthHandler.
func AuthHeader(id, username string) http.Header {
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(1 * time.Minute)
	signingKey := []byte("secret")

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwtCustomClaims{
			username,
			jwt.RegisteredClaims{
				Issuer:    "test",
				Subject:   id,
				Audience:  jwt.ClaimStrings{"all"},
				IssuedAt:  jwt.NewNumericDate(issuedAt),
				ExpiresAt: jwt.NewNumericDate(expiresAt),
				ID:        uuid.NewString(),
			},
		},
	)

	jwt, _ := token.SignedString(signingKey)

	header := http.Header{}
	header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))
	return header
}

func Token(id, username string) string {
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(1 * time.Minute)
	signingKey := []byte("secret")

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwtCustomClaims{
			username,
			jwt.RegisteredClaims{
				Issuer:    "test",
				Subject:   id,
				Audience:  jwt.ClaimStrings{"all"},
				IssuedAt:  jwt.NewNumericDate(issuedAt),
				ExpiresAt: jwt.NewNumericDate(expiresAt),
				ID:        uuid.NewString(),
			},
		},
	)

	jwt, _ := token.SignedString(signingKey)

	return jwt
}
