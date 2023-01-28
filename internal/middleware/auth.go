package middleware

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	j "github.com/hinccvi/go-ddd/pkg/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	contextKey int
	JWTKey     string
)

const (
	tokenInfoKey contextKey = iota
)

func (k JWTKey) Auth(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	tokenInfo, err := k.parseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	grpc_ctxtags.Extract(ctx).Set("auth.sub", tokenInfo.Subject)

	newCtx := context.WithValue(ctx, tokenInfoKey, tokenInfo)

	return newCtx, nil
}

func (k JWTKey) parseToken(t string) (j.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(t, &j.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(k), nil
	})

	if token == nil {
		return j.CustomClaims{}, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(*j.CustomClaims); ok && token.Valid {
		return *claims, nil
	}

	return j.CustomClaims{}, fmt.Errorf("[parseRefreshToken] internal error: %w", err)
}
