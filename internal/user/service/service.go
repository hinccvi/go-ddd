package service

import (
	"context"
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
	"github.com/jackc/pgx/v4"
)

type (
	// Service encapsulates usecase logic for user.
	Service interface {
		Get(ctx context.Context, id uuid.UUID) (entity.GetUserRow, error)
		Query(ctx context.Context, args entity.ListUserParams) (QueryUserResponse, error)
		Create(ctx context.Context, args entity.CreateUserParams) (entity.CreateUserRow, error)
		Update(ctx context.Context, args entity.UpdateUserParams) (entity.UpdateUserRow, error)
		Delete(ctx context.Context, id uuid.UUID) (entity.SoftDeleteUserRow, error)
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
		Limit  int32  `query:"limit"`
		Offset *int32 `query:"offset"`
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

	QueryUserResponse struct {
		List  []entity.ListUserRow `json:"list"`
		Total int64                `json:"total"`
	}
)

// NewService creates a new user service.
func New(rds redis.Client, repo repository.Repository, logger log.Logger, timeout time.Duration) Service {
	return service{rds, repo, logger, timeout}
}

func (s service) Get(ctx context.Context, id uuid.UUID) (entity.GetUserRow, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	item, err := s.repo.Get(ctx, id)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return entity.GetUserRow{}, errs.ErrResourceNotFound
	case err != nil:
		return entity.GetUserRow{}, fmt.Errorf("[Get] internal error: %w", err)
	}

	return item, nil
}

func (s service) Query(ctx context.Context, args entity.ListUserParams) (QueryUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	items, err := s.repo.Query(ctx, args)
	if err != nil {
		return QueryUserResponse{}, fmt.Errorf("[Query] internal error: %w", err)
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return QueryUserResponse{}, fmt.Errorf("[Query] internal error: %w", err)
	}

	return QueryUserResponse{items, total}, nil
}

func (s service) Count(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.repo.Count(ctx)
}

func (s service) Create(ctx context.Context, args entity.CreateUserParams) (entity.CreateUserRow, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	hashedPassword, err := tools.Bcrypt(args.Password)
	if err != nil {
		return entity.CreateUserRow{}, fmt.Errorf("[Create] internal error: %w", err)
	}

	args.Password = hashedPassword

	item, err := s.repo.Create(ctx, args)
	if err != nil {
		return entity.CreateUserRow{}, fmt.Errorf("[Create] internal error: %w", err)
	}

	return item, nil
}

func (s service) Update(ctx context.Context, args entity.UpdateUserParams) (entity.UpdateUserRow, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if args.Password != "" {
		hashedPassword, err := tools.Bcrypt(args.Password)
		if err != nil {
			return entity.UpdateUserRow{}, fmt.Errorf("[Update] internal error: %w", err)
		}

		args.Password = hashedPassword
	}

	item, err := s.repo.Update(ctx, args)
	if err != nil {
		return entity.UpdateUserRow{}, fmt.Errorf("[Update] internal error: %w", err)
	}

	return item, nil
}

func (s service) Delete(ctx context.Context, id uuid.UUID) (entity.SoftDeleteUserRow, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	item, err := s.repo.Delete(ctx, id)
	if err != nil {
		return entity.SoftDeleteUserRow{}, fmt.Errorf("[Delete] internal error: %w", err)
	}

	return item, nil
}
