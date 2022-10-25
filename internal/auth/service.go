package auth

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/config"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/constants"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/model"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	hTools "github.com/hinccvi/Golang-Project-Structure-Conventional/tools/hash"
)

type (
	// Service encapsulates the authentication logic.
	Service interface {
		// authenticate authenticates a user using username and password.
		// It returns a JWT token if authentication succeeds. Otherwise, an error is returned.
		Login(ctx context.Context) (loginResponse, error)
		Refresh(ctx context.Context) (refreshResponse, error)
	}

	// Identity represents an authenticated user identity.
	Identity interface {
		// GetID returns the user ID.
		GetID() string
		// GetName returns the user name.
		GetName() string
	}

	service struct {
		cfg    *config.Config
		rds    redis.Client
		logger log.Logger
		repo   Repository
	}

	key int

	loginResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	refreshResponse struct {
		RefreshToken string `json:"refresh_token"`
	}

	jwtCustomClaims constants.JWTCustomClaims
)

const (
	ctxUsername key = iota
	ctxPassword
	ctxAccessToken
	ctxRefreshToken
)

// NewService creates a new authentication service.
func NewService(cfg *config.Config, rds redis.Client, repo Repository, logger log.Logger) Service {
	return service{cfg, rds, logger, repo}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context) (loginResponse, error) {
	username, ok := ctx.Value(ctxUsername).(string)
	if !ok {
		return loginResponse{}, constants.ErrSystemError
	}

	password, ok := ctx.Value(ctxPassword).(string)
	if !ok {
		return loginResponse{}, constants.ErrSystemError
	}

	user, err := s.authenticate(ctx, username, password)
	if err != nil {
		return loginResponse{}, err
	}

	accessToken, err := s.generateJWT(user, "")
	if err != nil {
		return loginResponse{}, err
	}

	refreshToken, err := s.generateJWT(user, "refresh")
	if err != nil {
		return loginResponse{}, err
	}

	if err = s.cacheRefreshToken(ctx, user.ID.String(), refreshToken); err != nil {
		return loginResponse{}, err
	}

	return loginResponse{accessToken, refreshToken}, nil
}

func (s service) Refresh(ctx context.Context) (refreshResponse, error) {
	at, ok := ctx.Value(ctxAccessToken).(string)
	if !ok {
		return refreshResponse{}, constants.ErrSystemError
	}

	rt, ok := ctx.Value(ctxRefreshToken).(string)
	if !ok {
		return refreshResponse{}, constants.ErrSystemError
	}

	_, err := s.parseRefreshToken(rt)
	if err != nil {
		return refreshResponse{}, err
	}

	accessClaims, err := s.parseAccessToken(at)
	if err != nil {
		return refreshResponse{}, err
	}

	id, err := uuid.Parse(accessClaims.Subject)
	if err != nil {
		return refreshResponse{}, err
	}

	if err = s.validateRefreshToken(ctx, id.String(), rt); err != nil {
		return refreshResponse{}, err
	}

	user := model.GetByUsernameRow{
		ID:       id,
		Username: accessClaims.JWTData.UserName,
	}

	accessToken, err := s.generateJWT(user, "")
	if err != nil {
		return refreshResponse{}, err
	}

	return refreshResponse{accessToken}, nil
}

// authenticate authenticates a user using username and password.
// If name and password are correct, an identity is returned. Otherwise, nil is returned.
func (s service) authenticate(ctx context.Context, username, password string) (model.GetByUsernameRow, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return model.GetByUsernameRow{}, constants.ErrInvalidCredentials
	}

	if err = hTools.BcryptCompare(password, user.Password); err != nil {
		if err = s.cacheIncorrectPassword(ctx, user.ID.String()); err != nil {
			return model.GetByUsernameRow{}, err
		}

		return model.GetByUsernameRow{}, constants.ErrInvalidCredentials
	}

	return user, nil
}

// generateJWT generates a JWT that encodes an identity.
func (s service) generateJWT(user model.GetByUsernameRow, jwtType string) (string, error) {
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(time.Duration(s.cfg.Jwt.AccessExpiration) * time.Minute)
	signingKey := []byte(s.cfg.Jwt.AccessSigningKey)

	if jwtType == "refresh" {
		expiresAt = issuedAt.Add(time.Duration(s.cfg.Jwt.RefreshExpiration) * time.Minute)
		signingKey = []byte(s.cfg.Jwt.RefreshSigningKey)
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		constants.JWTCustomClaims{
			JWTData: constants.JWTData{UserName: user.Username},
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    s.cfg.App.Name,
				Subject:   user.ID.String(),
				Audience:  jwt.ClaimStrings{"all"},
				IssuedAt:  jwt.NewNumericDate(issuedAt),
				ExpiresAt: jwt.NewNumericDate(expiresAt),
				ID:        uuid.NewString(),
			},
		},
	)

	return token.SignedString(signingKey)
}

func (s service) parseRefreshToken(refreshToken string) (jwtCustomClaims, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, constants.ErrInvalidJwt
		}

		return []byte(s.cfg.Jwt.RefreshSigningKey), nil
	})

	if token == nil {
		return jwtCustomClaims{}, constants.ErrInvalidJwt
	}

	if claims, ok := token.Claims.(jwtCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return jwtCustomClaims{}, err
}

// parseAccessToken extract value from validated token that failed on expired err.
func (s service) parseAccessToken(accessToken string) (jwtCustomClaims, error) {
	_, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, constants.ErrInvalidJwt
		}

		return []byte(s.cfg.Jwt.AccessSigningKey), nil
	})

	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return jwtCustomClaims{}, err
	}

	token, _, err := new(jwt.Parser).ParseUnverified(accessToken, &jwtCustomClaims{})
	if err != nil {
		return jwtCustomClaims{}, err
	}

	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok {
		return jwtCustomClaims{}, err
	}

	// only allow access token to be refresh just before expire time
	if time.Until(claims.ExpiresAt.Time) > constants.JWTRemainingTime {
		return jwtCustomClaims{}, constants.ErrConditionNotFulfil
	}

	return *claims, nil
}

func (s service) cacheRefreshToken(ctx context.Context, id, refreshToken string) error {
	key := string(constants.GetRedisKey(constants.RefreshTokenKey)) + id

	return s.rds.Set(ctx, key, refreshToken, -1).Err()
}

func (s service) cacheIncorrectPassword(ctx context.Context, id string) error {
	key := string(constants.GetRedisKey(constants.IncorrectPasswordKey)) + id

	val, err := s.rds.Get(ctx, key).Int()
	switch {
	case errors.Is(err, redis.Nil):
		return s.rds.Set(ctx, key, 1, constants.IncorrectPasswordExpiration).Err()
	case err != nil:
		return err
	default:
		if val >= constants.MaxLoginAttempt {
			return constants.ErrMaxAttempt
		}

		return s.rds.Incr(ctx, key).Err()
	}
}

func (s service) validateRefreshToken(ctx context.Context, id, refreshToken string) error {
	key := string(constants.GetRedisKey(constants.RefreshTokenKey)) + id

	val, err := s.rds.Get(ctx, key).Result()
	switch {
	case errors.Is(err, redis.Nil):
		return constants.ErrInvalidRefreshToken
	case err != nil:
		return err
	default:
		if refreshToken != val {
			return constants.ErrInvalidRefreshToken
		}

		return nil
	}
}
