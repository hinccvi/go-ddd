package auth

import (
	"context"
	"time"

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
	logger log.Logger
	repo   Repository
}

// NewService creates a new authentication service.
func NewService(cfg *config.Config, repo Repository, logger log.Logger) Service {
	return service{cfg, logger, repo}
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

		return loginResponse{accessToken, refreshToken}, nil
	}
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
