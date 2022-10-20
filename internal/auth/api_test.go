package auth

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/config"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/constants"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/models"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/test"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	hTools "github.com/hinccvi/Golang-Project-Structure-Conventional/tools/hash"
	uTools "github.com/hinccvi/Golang-Project-Structure-Conventional/tools/uuid"
)

func TestAPI(t *testing.T) {
	id1, err := uTools.GenerateUUIDv4()
	for err != nil {
		id1, err = uTools.GenerateUUIDv4()
	}

	id2, err := uTools.GenerateUUIDv4()
	for err != nil {
		id2, err = uTools.GenerateUUIDv4()
	}

	hashedPassword, err := hTools.Bcrypt("secret", constants.BcryptCost)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	repo := &mockRepository{items: []models.User{
		{
			ID:        id1,
			Username:  "user",
			Password:  hashedPassword,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: sql.NullTime{},
		},
		{
			ID:        id2,
			Username:  "user2",
			Password:  hashedPassword,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: sql.NullTime{},
		},
	}}

	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	router := test.MockRouter(logger)

	var cfg config.Config
	cfg.App.Name = "test"
	cfg.Jwt.AccessSigningKey = "secret"
	cfg.Jwt.RefreshSigningKey = "secret"
	cfg.Jwt.AccessExpiration = 1
	cfg.Jwt.RefreshExpiration = 1

	refreshToken := test.MockRefreshToken(id2.String(), "user2")

	s := miniredis.RunT(t)

	rds, err := test.Redis(s.Addr())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	key := string(constants.GetRedisKey(constants.RefreshTokenKey)) + id2.String()

	rds.Set(context.TODO(), key, refreshToken, -1)

	RegisterHandlers(router.Group(""), NewService(&cfg, rds, repo, logger), logger)

	header := test.MockAuthHeader(id2.String(), "user2")

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
