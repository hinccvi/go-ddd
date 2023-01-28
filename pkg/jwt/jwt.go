package jwt

import "github.com/golang-jwt/jwt/v4"

type CustomClaims struct {
	UserName string `json:"username"`
	jwt.RegisteredClaims
}
