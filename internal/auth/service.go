package auth

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/config"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/constants"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/models"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	"gorm.io/gorm"
)

// Service encapsulates the authentication logic.
type Service interface {
	// authenticate authenticates a user using username and password.
	// It returns a JWT token if authentication succeeds. Otherwise, an error is returned.
	Login(ctx context.Context, username, password string) (string, error)
}

// Identity represents an authenticated user identity.
type Identity interface {
	// GetID returns the user ID.
	GetID() string
	// GetName returns the user name.
	GetName() string
}

type Data struct {
	UserId   string
	UserName string
}

type JwtCustomClaims struct {
	Data
	jwt.StandardClaims
}

type service struct {
	cfg    *config.Config
	logger log.Logger
	repo   Repository
}

// NewService creates a new authentication service.
func NewService(cfg *config.Config, repo Repository, logger log.Logger) Service {
	return service{cfg, logger, repo}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, username, password string) (string, error) {
	if user := s.authenticate(ctx, username, password); !reflect.DeepEqual(user, &models.User{}) {
		return s.generateJWT(user)
	}
	return "", constants.ErrInvalidCredentials
}

// authenticate authenticates a user using username and password.
// If name and password are correct, an identity is returned. Otherwise, nil is returned.
func (s service) authenticate(ctx context.Context, name, password string) models.User {
	logger := s.logger.With(ctx, "user", name)

	arg := &models.GetByUsernameAndPasswordParams{
		Username: name,
		Password: password,
	}

	user, err := s.repo.GetUserByUsernameAndPassword(ctx, arg)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Infof("authentication failed")
		}
		return *new(models.User)
	}

	logger.Infof("authentication successful")
	return user
}

// generateJWT generates a JWT that encodes an identity.
func (s service) generateJWT(user models.User) (string, error) {
	tokenObj := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		JwtCustomClaims{
			Data{user.ID.String(), user.Username},
			jwt.StandardClaims{
				Issuer:    "app",
				Subject:   user.Username,
				Audience:  "all",
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Duration(s.cfg.Jwt.AccessExpiration) * time.Minute).Unix(),
				Id:        uuid.NewString(),
			},
		},
	)

	tokenStr, err := tokenObj.SignedString([]byte(s.cfg.Jwt.AccessSigningKey))
	return tokenStr, err
}
