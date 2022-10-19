package test

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/constants"
	m "github.com/hinccvi/Golang-Project-Structure-Conventional/internal/middleware"
	"github.com/labstack/echo/v4"
)

// MockRouter creates a echo router for testing APIs.
func MockRouter(logger log.Logger) *echo.Echo {
	e := echo.New()

	e.HTTPErrorHandler = m.NewHTTPErrorHandler(constants.ErrorStatusCodeMaps).Handler(logger)

	e.Validator = &m.CustomValidator{Validator: validator.New()}

	return e
}

// MockAuthHeader returns an HTTP header that can pass the authentication check by MockAuthHandler.
func MockAuthHeader() http.Header {
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{})
	token, _ := jwt.SignedString([]byte("secret"))

	header := http.Header{}
	header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return header
}

// func MockGenerateJWT() (string, error) {

// }

// func MockParseJWT(token string) (jwt.MapClaims, error) {
// 	j, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, constants.ErrInvalidJwt
// 		}

// 		return []byte("secret"), nil
// 	})

// 	if claims, ok := j.Claims.(jwt.MapClaims); ok && j.Valid {
// 		return claims, nil
// 	}

// 	return jwt.MapClaims{}, err
// }
