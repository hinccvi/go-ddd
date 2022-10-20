package user

import (
	"database/sql"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/models"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/test"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	tools "github.com/hinccvi/Golang-Project-Structure-Conventional/tools/uuid"
	"github.com/labstack/echo/v4/middleware"
)

func TestAPI(t *testing.T) {
	uuid, err := tools.GenerateUUIDv4()
	for err != nil {
		uuid, err = tools.GenerateUUIDv4()
	}

	repo := &mockRepository{items: []models.User{
		{
			ID:        uuid,
			Username:  "user",
			Password:  "secret",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: sql.NullTime{}},
	}}

	authHandler := middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &jwt.MapClaims{},
		SigningKey: []byte("secret"),
	})

	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	router := test.MockRouter(logger)

	s := miniredis.RunT(t)

	rds, err := test.Redis(s.Addr())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	RegisterHandlers(router.Group(""), NewService(rds, repo, logger), logger, authHandler)
	header := test.MockAuthHeader(uuid.String(), "user")

	tests := []test.APITestCase{
		{
			Name:         "get all",
			Method:       http.MethodGet,
			URL:          "/v1/user/list",
			WantStatus:   http.StatusOK,
			WantResponse: fmt.Sprintf(`*{"list":[{"id":"%s","username":"user"}],"total":1}*`, uuid.String()),
		},
		{
			Name:         "get init user",
			Method:       http.MethodGet,
			URL:          fmt.Sprintf("/v1/user/%s", uuid.String()),
			Param:        uuid.String(),
			WantStatus:   http.StatusOK,
			WantResponse: fmt.Sprintf(`*{"id":"%s","username":"user"}*`, uuid.String()),
		},
		{
			Name:       "get unknown",
			Method:     http.MethodGet,
			URL:        "/v1/user/xxx",
			WantStatus: http.StatusBadRequest,
		},
		{
			Name:         "create ok",
			Method:       http.MethodPost,
			URL:          "/v1/user",
			Body:         `{"username": "user","password": "secret"}`,
			Header:       header,
			WantStatus:   http.StatusOK,
			WantResponse: `*"message":"created"*`,
		},
		{
			Name:         "create ok count",
			Method:       http.MethodGet,
			URL:          "/v1/user/list",
			Header:       nil,
			WantStatus:   http.StatusOK,
			WantResponse: `*"total":2*`,
		},
		{
			Name:       "create input error",
			Method:     http.MethodPost,
			URL:        "/v1/user",
			Body:       `"name":"test"}`,
			Header:     header,
			WantStatus: http.StatusBadRequest,
		},
		{
			Name:         "update ok",
			Method:       http.MethodPatch,
			URL:          "/v1/user",
			Body:         fmt.Sprintf(`{"id":"%s","username": "newuser","password": "newsecret"}`, uuid.String()),
			Header:       header,
			WantStatus:   http.StatusOK,
			WantResponse: `*"message":"updated"*`,
		},
		{
			Name:       "update verify",
			Method:     http.MethodPatch,
			URL:        "/v1/user",
			Header:     header,
			WantStatus: http.StatusBadRequest,
		},
		{
			Name:       "update auth error",
			Method:     http.MethodPatch,
			URL:        "/v1/user",
			Body:       `{"id":"xxx","username": "user","password": "secret"}`,
			WantStatus: http.StatusBadRequest,
		},
		{
			Name:         "update input error",
			Method:       http.MethodPatch,
			URL:          "/v1/user",
			Body:         `"name":"albumxyz"}`,
			Header:       header,
			WantStatus:   http.StatusBadRequest,
			WantResponse: "",
		},
		{
			Name:         "delete ok",
			Method:       http.MethodDelete,
			URL:          fmt.Sprintf("/v1/user/%s", uuid.String()),
			Header:       header,
			WantStatus:   http.StatusOK,
			WantResponse: `*"message":"deleted"*`,
		},
		{
			Name:       "delete verify",
			Method:     http.MethodDelete,
			URL:        "/v1/user/xxx",
			Header:     header,
			WantStatus: http.StatusBadRequest,
		},
		{
			Name:       "delete auth error",
			Method:     http.MethodDelete,
			URL:        fmt.Sprintf("/v1/user/%s", uuid.String()),
			WantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		test.Endpoint(t, router, tc)
	}
}
