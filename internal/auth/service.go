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
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/models"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/tools"
)

// Service encapsulates the authentication logic.
type Service interface {
	// authenticate authenticates a user using username and password.
	// It returns a JWT token if authentication succeeds. Otherwise, an error is returned.
	Login(ctx context.Context, username, password string) (loginResponse, error)
	Refresh(ctx context.Context, at, rt string) (refreshResponse, error)
}

// Identity represents an authenticated user identity.
type Identity interface {
	// GetID returns the user ID.
	GetID() string
	// GetName returns the user name.
	GetName() string
}

type Data struct {
	UserName string
}

type JwtCustomClaims struct {
	Data
	jwt.RegisteredClaims
}

type service struct {
	cfg    *config.Config
	rds    redis.Client
	logger log.Logger
	repo   Repository
}

// NewService creates a new authentication service.
func NewService(cfg *config.Config, rds redis.Client, repo Repository, logger log.Logger) Service {
	return service{cfg, rds, logger, repo}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, username, password string) (loginResponse, error) {
	if user, err := s.authenticate(ctx, username, password); err != nil {
		return loginResponse{}, err
	} else {
		accessToken, err := s.generateJWT(user, "")
		if err != nil {
			return loginResponse{}, err
		}

		refreshToken, err := s.generateJWT(user, "refresh")
		if err != nil {
			return loginResponse{}, err
		}

		if err := s.cacheRefreshToken(ctx, user.ID.String(), refreshToken); err != nil {
			return loginResponse{}, err
		}

		return loginResponse{accessToken, refreshToken}, nil
	}
}

func (s service) Refresh(ctx context.Context, at, rt string) (refreshResponse, error) {
	_, err := s.parseRefreshToken(rt)
	if err != nil {
		return refreshResponse{}, err
	}

	accessClaims, err := s.parseAccessToken(at)
	if err != nil {
		return refreshResponse{}, err
	}

	id := uuid.MustParse(accessClaims.Subject)

	if err := s.validateRefreshToken(ctx, id.String(), rt); err != nil {
		return refreshResponse{}, err
	}

	user := models.GetByUsernameRow{
		ID:       &id,
		Username: accessClaims.Data.UserName,
	}

	accessToken, err := s.generateJWT(user, "")
	if err != nil {
		return refreshResponse{}, err
	}

	return refreshResponse{accessToken}, nil
}

// authenticate authenticates a user using username and password.
// If name and password are correct, an identity is returned. Otherwise, nil is returned.
func (s service) authenticate(ctx context.Context, username, password string) (models.GetByUsernameRow, error) {
	logger := s.logger.With(ctx, "username", username)

	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return models.GetByUsernameRow{}, constants.ErrInvalidCredentials
	}

	if err := tools.BcryptCompare(password, user.Password); err != nil {
		if err := s.cacheIncorrectPassword(ctx, user.ID.String()); err != nil {
			return models.GetByUsernameRow{}, err
		}

		return models.GetByUsernameRow{}, constants.ErrInvalidCredentials
	}

	logger.Info("authentication successful")

	return user, nil
}

// generateJWT generates a JWT that encodes an identity.
func (s service) generateJWT(user models.GetByUsernameRow, jwtType string) (string, error) {
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(time.Duration(s.cfg.Jwt.AccessExpiration) * time.Minute)
	signingKey := []byte(s.cfg.Jwt.AccessSigningKey)

	if jwtType == "refresh" {
		expiresAt = issuedAt.Add(time.Duration(s.cfg.Jwt.RefreshExpiration) * time.Minute)
		signingKey = []byte(s.cfg.Jwt.RefreshSigningKey)
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		JwtCustomClaims{
			Data{user.Username},
			jwt.RegisteredClaims{
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

func (s service) parseRefreshToken(refreshToken string) (JwtCustomClaims, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, constants.ErrInvalidJwt
		}

		return []byte(s.cfg.Jwt.RefreshSigningKey), nil
	})

	if claims, ok := token.Claims.(JwtCustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return JwtCustomClaims{}, err
	}
}

// parseAccessToken extract value from validated token that failed on expired err
func (s service) parseAccessToken(accessToken string) (JwtCustomClaims, error) {
	_, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, constants.ErrInvalidJwt
		}

		return []byte(s.cfg.Jwt.AccessSigningKey), nil
	})

	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return JwtCustomClaims{}, err
	}

	token, _, err := new(jwt.Parser).ParseUnverified(accessToken, &JwtCustomClaims{})
	if err != nil {
		return JwtCustomClaims{}, nil
	}

	claims, ok := token.Claims.(*JwtCustomClaims)
	if !ok {
		return JwtCustomClaims{}, err
	}

	// only allow access token to be refresh 1 min before expire time
	if time.Until(claims.ExpiresAt.Time) > (60 * time.Second) {
		return JwtCustomClaims{}, constants.ErrConditionNotFulfil
	}

	return *claims, nil
}

func (s service) cacheRefreshToken(ctx context.Context, id, refreshToken string) error {
	key := constants.GetRedisKey(constants.RefreshTokenKey) + id

	return s.rds.Set(ctx, key, refreshToken, -1).Err()
}

func (s service) cacheIncorrectPassword(ctx context.Context, id string) error {
	key := constants.GetRedisKey(constants.IncorrectPasswordKey) + id

	val, err := s.rds.Get(ctx, key).Int()
	if err == redis.Nil {
		return s.rds.Set(ctx, key, 1, 24*time.Hour).Err()
	} else if err != nil {
		return err
	} else {
		if val >= 5 {
			return constants.ErrMaxAttempt
		}

		return s.rds.Incr(ctx, key).Err()
	}
}

func (s service) validateRefreshToken(ctx context.Context, id, refreshToken string) error {
	key := constants.GetRedisKey(constants.RefreshTokenKey) + id

	val, err := s.rds.Get(ctx, key).Result()
	if err == redis.Nil {
		return constants.ErrInvalidRefreshToken
	} else if err != nil {
		return err
	} else {
		if refreshToken != val {
			return constants.ErrInvalidRefreshToken
		} else {
			return nil
		}
	}
}
