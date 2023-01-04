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
		var resp entity.User
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
		var list []entity.User
		var total int64
		list, total, err = s.Query(context.TODO(), 1, 10)
		assert.NoError(t, err)
		assert.Len(t, list, 1)
		assert.Equal(t, total, int64(1))
	})

	t.Run("fail: db error", func(t *testing.T) {
		_, _, err = s.Query(context.TODO(), 1, -1)
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
		u := entity.User{
			Username: "user",
			Password: "secret",
		}
		err = s.Create(context.TODO(), u)
		assert.NoError(t, err)
	})

	t.Run("fail: empty field", func(t *testing.T) {
		err = s.Create(context.TODO(), entity.User{})
		assert.Error(t, err)
		assert.Equal(t, errs.ErrEmptyField, tools.UnwrapRecursive(err))
	})

	t.Run("fail: db error", func(t *testing.T) {
		err = s.Create(context.TODO(), entity.User{Username: "error", Password: "secret"})
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
		u := entity.User{
			ID:       id,
			Username: "newuser",
			Password: "newsecret",
		}
		err = s.Update(context.TODO(), u)
		assert.NoError(t, err)
	})

	t.Run("fail: not found", func(t *testing.T) {
		args := entity.User{
			ID:       uuid.New(),
			Username: "newuser",
			Password: "newsecret",
		}
		err = s.Update(context.TODO(), args)
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
		err = s.Delete(context.TODO(), id)
		assert.NoError(t, err)
	})

	t.Run("fail: not found", func(t *testing.T) {
		err = s.Delete(context.TODO(), uuid.New())
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, tools.UnwrapRecursive(err))
	})
}
