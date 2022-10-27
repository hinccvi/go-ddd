package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/hinccvi/go-ddd/internal/config"
	"github.com/hinccvi/go-ddd/internal/entity"
	errs "github.com/hinccvi/go-ddd/internal/errors"
	"github.com/hinccvi/go-ddd/internal/mocks"
	"github.com/hinccvi/go-ddd/pkg/log"
	"github.com/hinccvi/go-ddd/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	cfg, err := config.Load("local")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	rds, err := mocks.Redis(miniredis.RunT(t).Addr())
	assert.NoError(t, err)

	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	var repo mocks.AuthRepository

	t.Run("success", func(t *testing.T) {
		password, _ := tools.Bcrypt("secret")

		mockGetUserByUsername := entity.GetByUsernameRow{
			ID:       uuid.New(),
			Username: "user",
			Password: password,
		}
		repo.On("GetUserByUsername", mock.Anything, "user").Return(mockGetUserByUsername, nil).Once()
		s := service{&cfg, rds, logger, &repo, 2 * time.Second}

		req := LoginRequest{
			Username: "user",
			Password: "secret",
		}
		var resp loginResponse
		resp, err = s.Login(context.TODO(), req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("fail: incorrect credential", func(t *testing.T) {
		password, _ := tools.Bcrypt("anothersecret")

		mockGetUserByUsername := entity.GetByUsernameRow{
			ID:       uuid.New(),
			Username: "user",
			Password: password,
		}
		repo.On("GetUserByUsername", mock.Anything, "user").Return(mockGetUserByUsername, nil).Once()
		s := service{&cfg, rds, logger, &repo, 2 * time.Second}

		req := LoginRequest{
			Username: "user",
			Password: "secret",
		}
		_, err = s.Login(context.TODO(), req)
		assert.Error(t, err)
	})

	t.Run("fail: max attempt", func(t *testing.T) {
		mockGetUserByUsername := entity.GetByUsernameRow{
			ID:       uuid.New(),
			Username: "user",
			Password: "secret",
		}

		repo.On("GetUserByUsername", mock.Anything, "user").Return(mockGetUserByUsername, nil)
		s := service{&cfg, rds, logger, &repo, 2 * time.Second}

		i := 0
		for i < 6 {
			_, err = s.Login(context.TODO(), LoginRequest{
				Username: "user",
				Password: "secret",
			})

			err = errors.Unwrap(err)

			if assert.Error(t, err) && i < 5 {
				assert.Equal(t, errs.ErrInvalidCredentials, err)
			}

			i++
		}

		assert.Equal(t, errs.ErrMaxAttempt, errors.Unwrap(err))
	})
}

func TestRefresh(t *testing.T) {
	var cfg config.Config
	cfg.App.Name = "test"
	cfg.Jwt.AccessExpiration = 1
	cfg.Jwt.AccessSigningKey = "secret"
	cfg.Jwt.RefreshExpiration = 1
	cfg.Jwt.RefreshSigningKey = "secret"

	rds, err := mocks.Redis(miniredis.RunT(t).Addr())
	assert.NoError(t, err)

	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	var repo mocks.AuthRepository

	password, _ := tools.Bcrypt("secret")
	repo.On("GetUserByUsername", mock.Anything, "user").Return(
		entity.GetByUsernameRow{
			ID:       uuid.New(),
			Username: "user",
			Password: password,
		},
		nil,
	)

	t.Run("success", func(t *testing.T) {
		s := service{&cfg, rds, logger, &repo, 2 * time.Second}

		var loginResp loginResponse
		loginResp, err = s.Login(context.TODO(), LoginRequest{
			Username: "user",
			Password: "secret",
		})

		var refreshResp refreshResponse
		refreshResp, err = s.Refresh(context.TODO(), RefreshTokenRequest{
			AccessToken:  loginResp.AccessToken,
			RefreshToken: loginResp.RefreshToken,
		})
		assert.NoError(t, err)
		assert.NotNil(t, refreshResp)
	})

	t.Run("fail: invalid refresh token", func(t *testing.T) {
		s := service{&cfg, rds, logger, &repo, 2 * time.Second}

		var user entity.GetByUsernameRow
		user, err = s.repo.GetUserByUsername(context.TODO(), "user")
		assert.NoError(t, err)

		var accessJWT, refreshJWT string
		accessJWT, err = s.generateJWT(user, Access)
		assert.NoError(t, err)

		refreshJWT, err = s.generateJWT(user, Refresh)
		assert.NoError(t, err)

		_, err = s.Refresh(context.TODO(), RefreshTokenRequest{
			AccessToken:  accessJWT,
			RefreshToken: refreshJWT,
		})
		assert.Error(t, err)
		assert.Equal(t, errs.ErrInvalidRefreshToken, tools.UnwrapRecursive(err))
	})

	t.Run("fail: access token still valid", func(t *testing.T) {
		cfg.Jwt.AccessExpiration = 5

		s := service{&cfg, rds, logger, &repo, 2 * time.Second}

		var loginResp loginResponse
		loginResp, err = s.Login(context.TODO(), LoginRequest{
			Username: "user",
			Password: "secret",
		})

		_, err = s.Refresh(context.TODO(), RefreshTokenRequest{
			AccessToken:  loginResp.AccessToken,
			RefreshToken: loginResp.RefreshToken,
		})
		assert.Error(t, err)
		assert.Equal(t, errs.ErrConditionNotFulfil, tools.UnwrapRecursive(err))
	})

	t.Run("fail: invalid access token", func(t *testing.T) {
		s := service{&cfg, rds, logger, &repo, 2 * time.Second}

		var loginResp loginResponse
		loginResp, err = s.Login(context.TODO(), LoginRequest{
			Username: "user",
			Password: "secret",
		})

		_, err = s.Refresh(context.TODO(), RefreshTokenRequest{
			AccessToken:  loginResp.AccessToken + "x",
			RefreshToken: loginResp.RefreshToken,
		})
		assert.Error(t, err)
		assert.Equal(t, jwt.ErrSignatureInvalid, tools.UnwrapRecursive(err))
	})

	t.Run("fail: invalid refresh token", func(t *testing.T) {
		s := service{&cfg, rds, logger, &repo, 2 * time.Second}

		var loginResp loginResponse
		loginResp, err = s.Login(context.TODO(), LoginRequest{
			Username: "user",
			Password: "secret",
		})

		_, err = s.Refresh(context.TODO(), RefreshTokenRequest{
			AccessToken:  loginResp.AccessToken,
			RefreshToken: loginResp.RefreshToken + "x",
		})
		assert.Error(t, err)
		assert.Equal(t, jwt.ErrSignatureInvalid, tools.UnwrapRecursive(err))
	})
}
