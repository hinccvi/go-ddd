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
	"github.com/hinccvi/go-ddd/internal/user/repository"
	"github.com/hinccvi/go-ddd/pkg/log"
	"github.com/hinccvi/go-ddd/tools"
)

type (
	// Service encapsulates usecase logic for user.
	Service interface {
		GetUser(ctx context.Context, req GetUserRequest) (entity.User, error)
		QueryUser(ctx context.Context, req QueryUserRequest) ([]entity.User, int64, error)
		CreateUser(ctx context.Context, req CreateUserRequest) error
		UpdateUser(ctx context.Context, req UpdateUserRequest) error
		DeleteUser(ctx context.Context, req DeleteUserRequest) error
	}

	service struct {
		rds     redis.Client
		repo    repository.Repository
		logger  log.Logger
		timeout time.Duration
	}

	GetUserRequest struct {
		ID uuid.UUID `param:"id" validate:"required"`
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
		ID uuid.UUID `param:"id" validate:"required"`
	}
)

// NewService creates a new user service.
func New(rds redis.Client, repo repository.Repository, logger log.Logger, timeout time.Duration) Service {
	return service{rds, repo, logger, timeout}
}

func (s service) GetUser(ctx context.Context, req GetUserRequest) (entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	item, err := s.repo.GetUser(ctx, req.ID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.User{}, sql.ErrNoRows
	case err != nil:
		return entity.User{}, fmt.Errorf("[GetUser] internal error: %w", err)
	}

	return item, nil
}

func (s service) QueryUser(ctx context.Context, req QueryUserRequest) ([]entity.User, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	items, err := s.repo.QueryUser(ctx, req.Page, req.Size)
	if err != nil {
		return []entity.User{}, 0, fmt.Errorf("[QueryUser] internal error: %w", err)
	}

	total, err := s.repo.CountUser(ctx)
	if err != nil {
		return []entity.User{}, 0, fmt.Errorf("[QueryUser] internal error: %w", err)
	}

	return items, total, nil
}

func (s service) CreateUser(ctx context.Context, req CreateUserRequest) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	hashedPassword, err := tools.Bcrypt(req.Password)
	if err != nil {
		return fmt.Errorf("[CreateUser] internal error: %w", err)
	}
	req.Password = hashedPassword

	u := entity.User{
		Username: req.Username,
		Password: req.Password,
	}
	if err := s.repo.CreateUser(ctx, u); err != nil {
		return fmt.Errorf("[CreateUser] internal error: %w", err)
	}

	return nil
}

func (s service) UpdateUser(ctx context.Context, req UpdateUserRequest) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if req.Password != "" {
		hashedPassword, err := tools.Bcrypt(req.Password)
		if err != nil {
			return fmt.Errorf("[UpdateUser] internal error: %w", err)
		}

		req.Password = hashedPassword
	}

	u := entity.User{
		ID:       req.ID,
		Username: req.Username,
		Password: req.Password,
	}
	if err := s.repo.UpdateUser(ctx, u); err != nil {
		return fmt.Errorf("[UpdateUser] internal error: %w", err)
	}

	return nil
}

func (s service) DeleteUser(ctx context.Context, req DeleteUserRequest) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if err := s.repo.DeleteUser(ctx, req.ID); err != nil {
		return fmt.Errorf("[DeleteUser] internal error: %w", err)
	}

	return nil
}
