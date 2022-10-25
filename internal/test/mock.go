package test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/hinccvi/go-ddd/pkg/log"

	errs "github.com/hinccvi/go-ddd/internal/errors"
	m "github.com/hinccvi/go-ddd/internal/middleware"
	"github.com/labstack/echo/v4"
)

type (
	data struct {
		UserName string
	}

	jwtCustomClaims struct {
		data
		jwt.RegisteredClaims
	}
)

// MockRouter creates a echo router for testing APIs.
func MockRouter(logger log.Logger) *echo.Echo {
	e := echo.New()

	e.HTTPErrorHandler = m.NewHTTPErrorHandler(errs.GetStatusCodeMap()).Handler(logger)

	e.Validator = &m.CustomValidator{Validator: validator.New()}

	return e
}

// MockAuthHeader returns an HTTP header that can pass the authentication check by MockAuthHandler.
func MockAuthHeader(id, username string) http.Header {
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(1 * time.Minute)
	signingKey := []byte("secret")

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwtCustomClaims{
			data: data{UserName: username},
			RegisteredClaims: jwt.RegisteredClaims{
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

func MockRefreshToken(id, username string) string {
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(1 * time.Minute)
	signingKey := []byte("secret")

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwtCustomClaims{
			data{UserName: username},
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
