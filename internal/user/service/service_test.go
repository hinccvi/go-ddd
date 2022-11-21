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

func TestGet(t *testing.T) {
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
		var resp entity.GetUserRow
		resp, err = s.Get(context.TODO(), id)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("fail: not found", func(t *testing.T) {
		_, err = s.Get(context.TODO(), uuid.New())
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("fail: db error", func(t *testing.T) {
		_, err = s.Get(context.TODO(), uuid.UUID{})
		assert.Error(t, err)
		assert.Equal(t, mocks.ErrCRUD, tools.UnwrapRecursive(err))
	})
}

func TestQuery(t *testing.T) {
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
		var resp QueryUserResponse
		resp, err = s.Query(context.TODO(), entity.ListUserParams{Limit: 10, Offset: 0})
		assert.NoError(t, err)
		assert.Len(t, resp.List, 1)
		assert.Equal(t, resp.Total, int64(1))
	})

	t.Run("fail: db error", func(t *testing.T) {
		_, err = s.Query(context.TODO(), entity.ListUserParams{Limit: 10, Offset: -1})
		assert.Error(t, err)
	})
}

func TestCount(t *testing.T) {
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

	t.Run("success", func(t *testing.T) {
		s := service{rds, repo, logger, 2 * time.Second}

		var total int64
		total, err = s.Count(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, total, int64(1))
	})
}

func TestCreate(t *testing.T) {
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
		args := entity.CreateUserParams{
			Username: "user",
			Password: "secret",
		}
		var resp entity.CreateUserRow
		resp, err = s.Create(context.TODO(), args)
		assert.NoError(t, err)
		assert.NotEqual(t, resp.ID, uuid.UUID{})
		assert.Equal(t, resp.Username, "user")
	})

	t.Run("fail: empty field", func(t *testing.T) {
		_, err = s.Create(context.TODO(), entity.CreateUserParams{})
		assert.Error(t, err)
		assert.Equal(t, errs.ErrEmptyField, tools.UnwrapRecursive(err))
	})

	t.Run("fail: db error", func(t *testing.T) {
		_, err = s.Create(context.TODO(), entity.CreateUserParams{Username: "error", Password: "secret"})
		assert.Error(t, err)
		assert.Equal(t, mocks.ErrCRUD, tools.UnwrapRecursive(err))
	})
}

func TestUpdate(t *testing.T) {
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
		args := entity.UpdateUserParams{
			ID:       id,
			Username: "newuser",
			Password: "newsecret",
		}
		var resp entity.UpdateUserRow
		resp, err = s.Update(context.TODO(), args)
		assert.NoError(t, err)
		assert.Equal(t, resp.Username, "newuser")
	})

	t.Run("fail: not found", func(t *testing.T) {
		args := entity.UpdateUserParams{
			ID:       uuid.New(),
			Username: "newuser",
			Password: "newsecret",
		}
		_, err = s.Update(context.TODO(), args)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, tools.UnwrapRecursive(err))
	})
}

func TestDelete(t *testing.T) {
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
		var resp entity.SoftDeleteUserRow
		resp, err = s.Delete(context.TODO(), id)
		assert.NoError(t, err)
		assert.Equal(t, resp.Username, "user")
	})

	t.Run("fail: not found", func(t *testing.T) {
		_, err = s.Delete(context.TODO(), uuid.New())
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, tools.UnwrapRecursive(err))
	})
}
