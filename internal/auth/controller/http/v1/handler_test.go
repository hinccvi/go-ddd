package v1

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/hinccvi/go-ddd/internal/auth/service"
	"github.com/hinccvi/go-ddd/internal/config"
	"github.com/hinccvi/go-ddd/internal/entity"
	"github.com/hinccvi/go-ddd/internal/mocks"
	"github.com/hinccvi/go-ddd/internal/test"
	"github.com/hinccvi/go-ddd/pkg/log"
	"github.com/hinccvi/go-ddd/tools"
	"github.com/stretchr/testify/mock"
)

func TestHandler(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()

	hashedPassword, err := tools.Bcrypt("secret")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	mockGetUserByUsername := []entity.User{
		{
			ID:       id1,
			Username: "user",
			Password: hashedPassword,
		},
		{
			ID:       id2,
			Username: "user1",
			Password: hashedPassword,
		},
	}

	var repo mocks.AuthRepository
	repo.On("GetUserByUsername", mock.Anything, "user").Return(mockGetUserByUsername[0], nil)
	repo.On("GetUserByUsername", mock.Anything, "user1").Return(mockGetUserByUsername[1], nil)

	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	router := mocks.Router(logger)

	var cfg config.Config
	cfg.App.Name = "test"
	cfg.Jwt.AccessSigningKey = "secret"
	cfg.Jwt.RefreshSigningKey = "secret"
	cfg.Jwt.AccessExpiration = 1
	cfg.Jwt.RefreshExpiration = 1

	refreshToken := mocks.Token(id2.String(), "user1")

	header := mocks.AuthHeader(id2.String(), "user1")

	rds, err := mocks.Redis(miniredis.RunT(t).Addr())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	rds.Set(context.TODO(), mocks.RefreshTokenKey(id2.String()), refreshToken, -1)

	RegisterHandlers(router.Group("v1"), service.New(&cfg, rds, &repo, logger, 2*time.Second), logger)

	tests := []test.APITestCase{
		{
			Name:         "login ok",
			Method:       http.MethodPost,
			URL:          "/v1/auth/login",
			Body:         `{"username":"user","password":"secret"}`,
			WantStatus:   http.StatusOK,
			WantResponse: `*"message":"success"*`,
		},
		{
			Name:       "login fail",
			Method:     http.MethodPost,
			URL:        "/v1/auth/login",
			Body:       `{"username":"user","password":"xxx"}`,
			WantStatus: http.StatusBadRequest,
		},
		{
			Name:       "login bind fail",
			Method:     http.MethodPost,
			URL:        "/v1/auth/login",
			Body:       `"username":"user","password":"secret"}`,
			WantStatus: http.StatusBadRequest,
		},
		{
			Name:       "login bind fail",
			Method:     http.MethodPost,
			URL:        "/v1/auth/login",
			Body:       `{"username":255,"password":255}`,
			WantStatus: http.StatusBadRequest,
		},
		{
			Name:       "login validate fail",
			Method:     http.MethodPost,
			URL:        "/v1/auth/login",
			Body:       `{"username":"","password":"secret"}`,
			WantStatus: http.StatusBadRequest,
		},
		{
			Name:         "refresh token ok",
			Method:       http.MethodPost,
			URL:          "/v1/auth/refresh",
			Body:         fmt.Sprintf(`{"refresh_token":"%s"}`, refreshToken),
			Header:       header,
			WantStatus:   http.StatusOK,
			WantResponse: `*"message":"success"*`,
		},
		{
			Name:       "refresh token fail",
			Method:     http.MethodPost,
			URL:        "/v1/auth/refresh",
			Body:       `{"refresh_token":"xxx"}`,
			WantStatus: http.StatusForbidden,
		},
		{
			Name:       "refresh token bind fail",
			Method:     http.MethodPost,
			URL:        "/v1/auth/refresh",
			Body:       `"refresh_token":"xxx"}`,
			WantStatus: http.StatusBadRequest,
		},
		{
			Name:       "refresh token bind fail",
			Method:     http.MethodPost,
			URL:        "/v1/auth/refresh",
			Body:       `{"refresh_token":1}`,
			WantStatus: http.StatusBadRequest,
		},
		{
			Name:       "refresh token validate fail",
			Method:     http.MethodPost,
			URL:        "/v1/auth/refresh",
			Body:       `{"refresh_token":""}`,
			WantStatus: http.StatusBadRequest,
		},
		{
			Name:       "refresh token validate fail",
			Method:     http.MethodPost,
			URL:        "/v1/auth/refresh",
			Body:       `{"refresh_token":"xxx"}`,
			Header:     header,
			WantStatus: http.StatusForbidden,
		},
	}

	for _, tc := range tests {
		test.Endpoint(t, router, tc)
	}
}
