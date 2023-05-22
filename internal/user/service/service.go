package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/hinccvi/go-ddd/internal/entity"
	errs "github.com/hinccvi/go-ddd/internal/errors"
	"github.com/hinccvi/go-ddd/internal/user/repository"
	"github.com/hinccvi/go-ddd/pkg/log"
	"github.com/hinccvi/go-ddd/tools"
)

type (
	// Service encapsulates usecase logic for user.
	Service interface {
		Get(ctx context.Context, id uuid.UUID) (entity.User, error)
		Query(ctx context.Context, page, size int) ([]entity.User, int64, error)
		Create(ctx context.Context, u entity.User) error
		Update(ctx context.Context, u entity.User) error
		Delete(ctx context.Context, id uuid.UUID) error
	}

	service struct {
		rds     redis.Client
		repo    repository.Repository
		logger  log.Logger
		timeout time.Duration
	}

	GetUserRequest struct {
		ID *uuid.UUID `param:"id" validate:"required"`
	}

	QueryUserRequest struct {
		Page int `query:"page"`
		Size int `query:"size"`
	}

	CreateUserRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	UpdateUserRequest struct {
		ID       uuid.UUID `json:"id" validate:"required"`
		Username string    `json:"username"`
		Password string    `json:"password"`
	}

	DeleteUserRequest struct {
		ID *uuid.UUID `param:"id" validate:"required"`
	}
)

// NewService creates a new user service.
func New(rds redis.Client, repo repository.Repository, logger log.Logger, timeout time.Duration) Service {
	return service{rds, repo, logger, timeout}
}

func (s service) Get(ctx context.Context, id uuid.UUID) (entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	item, err := s.repo.Get(ctx, id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.User{}, sql.ErrNoRows
	case err != nil:
		return entity.User{}, fmt.Errorf("[Get] internal error: %w", err)
	}

	return item, nil
}

func (s service) Query(ctx context.Context, page, size int) ([]entity.User, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	items, err := s.repo.Query(ctx, page, size)
	if err != nil {
		return []entity.User{}, 0, fmt.Errorf("[Query] internal error: %w", err)
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return []entity.User{}, 0, fmt.Errorf("[Query] internal error: %w", err)
	}

	return items, total, nil
}

func (s service) Count(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.repo.Count(ctx)
}

func (s service) Create(ctx context.Context, u entity.User) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if u.Username == "" || u.Password == "" {
		return fmt.Errorf("[Create] internal error: %w", errs.ErrEmptyField)
	}

	hashedPassword, err := tools.Bcrypt(u.Password)
	if err != nil {
		return fmt.Errorf("[Create] internal error: %w", err)
	}
	u.Password = hashedPassword

	if err = s.repo.Create(ctx, u); err != nil {
		return fmt.Errorf("[Create] internal error: %w", err)
	}

	return nil
}

func (s service) Update(ctx context.Context, u entity.User) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if u.Password != "" {
		hashedPassword, err := tools.Bcrypt(u.Password)
		if err != nil {
			return fmt.Errorf("[Update] internal error: %w", err)
		}

		u.Password = hashedPassword
	}

	if err := s.repo.Update(ctx, u); err != nil {
		return fmt.Errorf("[Update] internal error: %w", err)
	}

	return nil
}

func (s service) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("[Delete] internal error: %w", err)
	}

	return nil
}
