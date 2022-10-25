package user

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/constants"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/model"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	tools "github.com/hinccvi/Golang-Project-Structure-Conventional/tools/hash"
	"github.com/jackc/pgx/v4"
)

type (
	// Service encapsulates usecase logic for user.
	Service interface {
		Get(ctx context.Context, id uuid.UUID) (model.GetUserRow, error)
		Query(ctx context.Context, args model.ListUserParams) (List, error)
		Create(ctx context.Context, args model.CreateUserParams) (model.CreateUserRow, error)
		Update(ctx context.Context, args model.UpdateUserParams) (model.UpdateUserRow, error)
		Delete(ctx context.Context, id uuid.UUID) (model.SoftDeleteUserRow, error)
	}

	List struct {
		List  []model.ListUserRow `json:"list"`
		Total int64               `json:"total"`
	}

	service struct {
		rds    redis.Client
		repo   Repository
		logger log.Logger
	}
)

// NewService creates a new user service.
func NewService(rds redis.Client, repo Repository, logger log.Logger) Service {
	return service{rds, repo, logger}
}

func (s service) Get(ctx context.Context, id uuid.UUID) (model.GetUserRow, error) {
	item, err := s.repo.Get(ctx, id)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return model.GetUserRow{}, err
	case err != nil:
		return model.GetUserRow{}, constants.ErrResourceNotFound
	}

	return item, nil
}

func (s service) Query(ctx context.Context, args model.ListUserParams) (List, error) {
	items, err := s.repo.Query(ctx, args)
	if err != nil {
		return List{}, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return List{}, err
	}

	return List{items, total}, nil
}

func (s service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

func (s service) Create(ctx context.Context, args model.CreateUserParams) (model.CreateUserRow, error) {
	hashedPassword, err := tools.Bcrypt(args.Password, constants.BcryptCost)
	if err != nil {
		return model.CreateUserRow{}, err
	}

	args.Password = hashedPassword

	item, err := s.repo.Create(ctx, args)
	if err != nil {
		return model.CreateUserRow{}, err
	}

	return item, nil
}

func (s service) Update(ctx context.Context, args model.UpdateUserParams) (model.UpdateUserRow, error) {
	if args.Password != "" {
		hashedPassword, err := tools.Bcrypt(args.Password, constants.BcryptCost)
		if err != nil {
			return model.UpdateUserRow{}, err
		}

		args.Password = hashedPassword
	}

	item, err := s.repo.Update(ctx, args)
	if err != nil {
		return model.UpdateUserRow{}, err
	}

	return item, nil
}

func (s service) Delete(ctx context.Context, id uuid.UUID) (model.SoftDeleteUserRow, error) {
	item, err := s.repo.Delete(ctx, id)
	if err != nil {
		return model.SoftDeleteUserRow{}, err
	}

	return item, nil
}
