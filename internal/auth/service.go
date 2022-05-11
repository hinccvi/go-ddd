package auth

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/entity"
	errs "github.com/hinccvi/Golang-Project-Structure-Conventional/internal/errors"
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

type MyCustomClaims struct {
	Data
	jwt.StandardClaims
}

type service struct {
	signingKey      string
	tokenExpiration int
	logger          log.Logger
	repo            Repository
}

// NewService creates a new authentication service.
func NewService(signingKey string, tokenExpiration int, repo Repository, logger log.Logger) Service {
	return service{signingKey, tokenExpiration, logger, repo}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, username, password string) (string, error) {
	if user := s.authenticate(ctx, username, password); !reflect.DeepEqual(user, &entity.User{}) {
		return s.generateJWT(user)
	}
	return "", errs.Unauthorized("incorrect username or password")
}

// authenticate authenticates a user using username and password.
// If name and password are correct, an identity is returned. Otherwise, nil is returned.
func (s service) authenticate(ctx context.Context, name, password string) entity.User {
	logger := s.logger.With(ctx, "user", name)

	user, err := s.repo.GetUserByUsernameAndPassword(ctx, name, password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Infof("authentication failed")
		}
		return *new(entity.User)
	}

	logger.Infof("authentication successful")
	return user
}

// generateJWT generates a JWT that encodes an identity.
func (s service) generateJWT(user entity.User) (string, error) {
	tokenObj := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		MyCustomClaims{
			Data{user.ID, user.Name},
			jwt.StandardClaims{
				Issuer:    "app",
				Subject:   user.Name,
				Audience:  "all",
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Duration(s.tokenExpiration) * time.Hour).Unix(),
				Id:        uuid.NewString(),
			},
		},
	)

	fmt.Println(s.signingKey)
	tokenStr, err := tokenObj.SignedString([]byte(s.signingKey))
	return tokenStr, err

}
