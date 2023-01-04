package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/hinccvi/go-ddd/internal/auth/repository"
	"github.com/hinccvi/go-ddd/internal/config"
	"github.com/hinccvi/go-ddd/internal/entity"
	errs "github.com/hinccvi/go-ddd/internal/errors"
	"github.com/hinccvi/go-ddd/pkg/log"
	"github.com/hinccvi/go-ddd/tools"
)

type (
	// Service encapsulates the authentication logic.
	Service interface {
		// authenticate authenticates a user using username and password.
		// It returns a JWT token if authentication succeeds. Otherwise, an error is returned.
		Login(ctx context.Context, req LoginRequest) (loginResponse, error)
		Refresh(ctx context.Context, req RefreshTokenRequest) (refreshResponse, error)
	}

	service struct {
		cfg     *config.Config
		rds     redis.Client
		logger  log.Logger
		repo    repository.Repository
		timeout time.Duration
	}

	JWTCustomClaims struct {
		UserName string `json:"username"`
		jwt.RegisteredClaims
	}

	// http request struct.
	LoginRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
		AccessToken  string
	}

	// http response struct.
	loginResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	refreshResponse struct {
		RefreshToken string `json:"refresh_token"`
	}

	jwtType string

	RedisKey string
)

const (
	JWTBearerFormat = 2

	Access  jwtType = "access"
	Refresh jwtType = "refresh"

	jwtRemainingTime                     = 60 * time.Second
	maxLoginAttempt                      = 5
	incorrectPasswordExpiration          = 24 * time.Hour
	prefix                      RedisKey = "app"
	refreshToken                RedisKey = "refresh_token"
	incorrectPassword           RedisKey = "incorrect_password"
	smsCooldown                 RedisKey = "sms_cooldown"
	smsCode                     RedisKey = "sms_code"
	smsLimit                    RedisKey = "sms_limit"
	smsAttempt                  RedisKey = "sms_attempt"
)

// New creates a new authentication service.
func New(
	cfg *config.Config,
	rds redis.Client,
	repo repository.Repository,
	logger log.Logger,
	timeout time.Duration,
) Service {
	return service{cfg, rds, logger, repo, timeout}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, req LoginRequest) (loginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	user, err := s.authenticate(ctx, req.Username, req.Password)
	if err != nil {
		return loginResponse{}, fmt.Errorf("[Login] internal error: %w", err)
	}

	accessToken, err := s.generateJWT(user.ID, user.Username, Access)
	if err != nil {
		return loginResponse{}, fmt.Errorf("[Login] internal error: %w", err)
	}

	refreshToken, err := s.generateJWT(user.ID, user.Username, Refresh)
	if err != nil {
		return loginResponse{}, fmt.Errorf("[Login] internal error: %w", err)
	}

	if err = s.cacheRefreshToken(ctx, user.ID.String(), refreshToken); err != nil {
		return loginResponse{}, fmt.Errorf("[Login] internal error: %w", err)
	}

	return loginResponse{accessToken, refreshToken}, nil
}

func (s service) Refresh(ctx context.Context, req RefreshTokenRequest) (refreshResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	_, err := s.parseRefreshToken(req.RefreshToken)
	if err != nil {
		return refreshResponse{}, fmt.Errorf("[Refresh] internal error: %w", err)
	}

	accessClaims, err := s.parseAccessToken(req.AccessToken)
	if err != nil {
		return refreshResponse{}, fmt.Errorf("[Refresh] internal error: %w", err)
	}

	id, err := uuid.Parse(accessClaims.Subject)
	if err != nil {
		return refreshResponse{}, fmt.Errorf("[Refresh] internal error: %w", err)
	}

	if err = s.validateRefreshToken(ctx, id.String(), req.RefreshToken); err != nil {
		return refreshResponse{}, fmt.Errorf("[Refresh] internal error: %w", err)
	}

	accessToken, err := s.generateJWT(id, accessClaims.UserName, Access)
	if err != nil {
		return refreshResponse{}, fmt.Errorf("[Refresh] internal error: %w", err)
	}

	return refreshResponse{accessToken}, nil
}

// authenticate authenticates a user using username and password.
// If name and password are correct, an identity is returned. Otherwise, nil is returned.
func (s service) authenticate(ctx context.Context, username, password string) (entity.User, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.User{}, errs.ErrInvalidCredentials
	} else if err != nil {
		return entity.User{}, err
	}

	if err = tools.BcryptCompare(password, user.Password); err != nil {
		if err = s.cacheIncorrectPassword(ctx, user.ID.String()); err != nil {
			return entity.User{}, fmt.Errorf("[authenticate] internal error: %w", err)
		}

		return entity.User{}, errs.ErrInvalidCredentials
	}

	return user, nil
}

// generateJWT generates a JWT that encodes an identity.
func (s service) generateJWT(id uuid.UUID, username string, t jwtType) (string, error) {
	issuedAt := time.Now()
	var expiresAt time.Time
	var signingKey []byte

	if t == Refresh {
		expiresAt = issuedAt.Add(time.Duration(s.cfg.Jwt.RefreshExpiration) * time.Minute)
		signingKey = []byte(s.cfg.Jwt.RefreshSigningKey)
	} else {
		expiresAt = issuedAt.Add(time.Duration(s.cfg.Jwt.AccessExpiration) * time.Minute)
		signingKey = []byte(s.cfg.Jwt.AccessSigningKey)
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&JWTCustomClaims{
			username,
			jwt.RegisteredClaims{
				Issuer:    s.cfg.App.Name,
				Subject:   id.String(),
				Audience:  jwt.ClaimStrings{"all"},
				IssuedAt:  jwt.NewNumericDate(issuedAt),
				ExpiresAt: jwt.NewNumericDate(expiresAt),
				ID:        uuid.NewString(),
			},
		},
	)

	return token.SignedString(signingKey)
}

func (s service) parseRefreshToken(refreshToken string) (JWTCustomClaims, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.ErrInvalidJwt
		}

		return []byte(s.cfg.Jwt.RefreshSigningKey), nil
	})

	if token == nil {
		return JWTCustomClaims{}, errs.ErrInvalidJwt
	}

	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return *claims, nil
	}

	return JWTCustomClaims{}, fmt.Errorf("[parseRefreshToken] internal error: %w", err)
}

// parseAccessToken extract value from validated token that failed on expired err.
func (s service) parseAccessToken(accessToken string) (JWTCustomClaims, error) {
	_, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.ErrInvalidJwt
		}

		return []byte(s.cfg.Jwt.AccessSigningKey), nil
	})

	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return JWTCustomClaims{}, fmt.Errorf("[parseAccessToken] internal error: %w", err)
	}

	token, _, err := new(jwt.Parser).ParseUnverified(accessToken, &JWTCustomClaims{})
	if err != nil {
		return JWTCustomClaims{}, fmt.Errorf("[parseAccessToken] internal error: %w", err)
	}

	claims, ok := token.Claims.(*JWTCustomClaims)
	if !ok {
		return JWTCustomClaims{}, fmt.Errorf("[parseAccessToken] internal error: %w", err)
	}

	// only allow access token to be refresh just before expire time
	if time.Until(claims.ExpiresAt.Time) > jwtRemainingTime {
		return JWTCustomClaims{}, errs.ErrConditionNotFulfil
	}

	return *claims, nil
}

func (s service) cacheRefreshToken(ctx context.Context, id, token string) error {
	key := s.getRedisKey(refreshToken, id)

	return s.rds.Set(ctx, key, token, -1).Err()
}

func (s service) cacheIncorrectPassword(ctx context.Context, id string) error {
	key := s.getRedisKey(incorrectPassword, id)

	val, err := s.rds.Get(ctx, key).Int()
	switch {
	case errors.Is(err, redis.Nil):
		return s.rds.Set(ctx, key, 1, incorrectPasswordExpiration).Err()
	case err != nil:
		return fmt.Errorf("[cacheIncorrectPassword] internal error: %w", err)
	default:
		if val >= maxLoginAttempt {
			return errs.ErrMaxAttempt
		}

		return s.rds.Incr(ctx, key).Err()
	}
}

func (s service) validateRefreshToken(ctx context.Context, id, token string) error {
	key := s.getRedisKey(refreshToken, id)

	val, err := s.rds.Get(ctx, key).Result()
	switch {
	case errors.Is(err, redis.Nil):
		return errs.ErrInvalidRefreshToken
	case err != nil:
		return fmt.Errorf("[validateRefreshToken] internal error: %w", err)
	default:
		if token != val {
			return errs.ErrInvalidRefreshToken
		}

		return nil
	}
}

func (s service) getRedisKey(key RedisKey, field string) string {
	return fmt.Sprintf("%s:%s:%s", s.cfg.App.Name, string(key), field)
}
