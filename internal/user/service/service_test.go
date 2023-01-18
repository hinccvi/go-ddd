package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/hinccvi/go-ddd/internal/config"
	"github.com/hinccvi/go-ddd/internal/entity"
	errs "github.com/hinccvi/go-ddd/internal/errors"
	"github.com/hinccvi/go-ddd/internal/mocks"
	"github.com/hinccvi/go-ddd/pkg/log"
	"github.com/hinccvi/go-ddd/tools"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	cfg, err := config.Load("local")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	rds, err := mocks.Redis(miniredis.RunT(t).Addr())
	assert.NoError(t, err)

	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	id := uuid.New()
	repo := &mocks.UserRepository{Items: []entity.User{
		{
			ID:        id,
			Username:  "user",
			Password:  "secret",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: sql.NullTime{}},
	}}
	s := service{rds, repo, logger, 2 * time.Second}

	t.Run("success", func(t *testing.T) {
		var resp entity.User
		resp, err = s.GetUser(context.TODO(), GetUserRequest{id})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("fail: not found", func(t *testing.T) {
		_, err = s.GetUser(context.TODO(), GetUserRequest{uuid.New()})
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("fail: db error", func(t *testing.T) {
		_, err = s.GetUser(context.TODO(), GetUserRequest{uuid.UUID{}})
		assert.Error(t, err)
		assert.Equal(t, mocks.ErrCRUD, tools.UnwrapRecursive(err))
	})
}

func TestQueryUser(t *testing.T) {
	cfg, err := config.Load("local")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	rds, err := mocks.Redis(miniredis.RunT(t).Addr())
	assert.NoError(t, err)

	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	id := uuid.New()
	repo := &mocks.UserRepository{Items: []entity.User{
		{
			ID:        id,
			Username:  "user",
			Password:  "secret",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: sql.NullTime{}},
	}}
	s := service{rds, repo, logger, 2 * time.Second}

	t.Run("success", func(t *testing.T) {
		var list []entity.User
		var total int64
		list, total, err = s.QueryUser(context.TODO(), QueryUserRequest{Size: 10, Page: 1})
		assert.NoError(t, err)
		assert.Len(t, list, 1)
		assert.Equal(t, total, int64(1))
	})

	t.Run("fail: db error", func(t *testing.T) {
		_, _, err = s.QueryUser(context.TODO(), QueryUserRequest{Size: -1, Page: -1})
		assert.Error(t, err)
	})
}

func TestCreateUser(t *testing.T) {
	cfg, err := config.Load("local")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	rds, err := mocks.Redis(miniredis.RunT(t).Addr())
	assert.NoError(t, err)

	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	repo := &mocks.UserRepository{}
	s := service{rds, repo, logger, 2 * time.Second}

	t.Run("success", func(t *testing.T) {
		err = s.CreateUser(context.TODO(), CreateUserRequest{Username: "user", Password: "secret"})
		assert.NoError(t, err)
	})

	t.Run("fail: empty field", func(t *testing.T) {
		err = s.CreateUser(context.TODO(), CreateUserRequest{})
		assert.Error(t, err)
		assert.Equal(t, errs.ErrEmptyField, tools.UnwrapRecursive(err))
	})

	t.Run("fail: db error", func(t *testing.T) {
		err = s.CreateUser(context.TODO(), CreateUserRequest{Username: "error", Password: "secret"})
		assert.Error(t, err)
		assert.Equal(t, mocks.ErrCRUD, tools.UnwrapRecursive(err))
	})
}

func TestUpdateUser(t *testing.T) {
	cfg, err := config.Load("local")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	rds, err := mocks.Redis(miniredis.RunT(t).Addr())
	assert.NoError(t, err)

	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	id := uuid.New()
	repo := &mocks.UserRepository{Items: []entity.User{
		{
			ID:        id,
			Username:  "user",
			Password:  "secret",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: sql.NullTime{}},
	}}
	s := service{rds, repo, logger, 2 * time.Second}

	t.Run("success", func(t *testing.T) {
		err = s.UpdateUser(context.TODO(), UpdateUserRequest{
			ID:       id,
			Username: "newuser",
			Password: "newsecret",
		})
		assert.NoError(t, err)
	})

	t.Run("fail: not found", func(t *testing.T) {
		err = s.UpdateUser(context.TODO(), UpdateUserRequest{
			ID:       uuid.New(),
			Username: "newuser",
			Password: "newsecret",
		})
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, tools.UnwrapRecursive(err))
	})
}

func TestDeleteUser(t *testing.T) {
	cfg, err := config.Load("local")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	rds, err := mocks.Redis(miniredis.RunT(t).Addr())
	assert.NoError(t, err)

	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	id := uuid.New()
	repo := &mocks.UserRepository{Items: []entity.User{
		{
			ID:        id,
			Username:  "user",
			Password:  "secret",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: sql.NullTime{}},
	}}
	s := service{rds, repo, logger, 2 * time.Second}

	t.Run("success", func(t *testing.T) {
		err = s.DeleteUser(context.TODO(), DeleteUserRequest{id})
		assert.NoError(t, err)
	})

	t.Run("fail: not found", func(t *testing.T) {
		err = s.DeleteUser(context.TODO(), DeleteUserRequest{uuid.New()})
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, tools.UnwrapRecursive(err))
	})
}
