package user

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/constants"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/models"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	tools "github.com/hinccvi/Golang-Project-Structure-Conventional/tools/hash"
	"github.com/jackc/pgx/v4"
)

type (
	// Service encapsulates usecase logic for user.
	Service interface {
		Get(ctx context.Context) (models.GetUserRow, error)
		Query(ctx context.Context) (List, error)
		Create(ctx context.Context) (models.CreateUserRow, error)
		Update(ctx context.Context) (models.UpdateUserRow, error)
		Delete(ctx context.Context) (models.SoftDeleteUserRow, error)
	}

	List struct {
		List  []models.ListUserRow `json:"list"`
		Total int64                `json:"total"`
	}

	service struct {
		rds    redis.Client
		repo   Repository
		logger log.Logger
	}

	key int
)

const (
	ctxID key = iota
	ctxListUser
	ctxCreateUser
	ctxUpdateUser
)

// NewService creates a new user service.
func NewService(rds redis.Client, repo Repository, logger log.Logger) Service {
	return service{rds, repo, logger}
}

func (s service) Get(ctx context.Context) (models.GetUserRow, error) {
	id, ok := ctx.Value(ctxID).(uuid.UUID)
	if !ok {
		return models.GetUserRow{}, constants.ErrSystemError
	}

	item, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.GetUserRow{}, constants.ErrResourceNotFound
		}

		return models.GetUserRow{}, err
	}

	return item, nil
}

func (s service) Query(ctx context.Context) (List, error) {
	args, ok := ctx.Value(ctxListUser).(models.ListUserParams)
	if !ok {
		return List{}, constants.ErrSystemError
	}

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

func (s service) Create(ctx context.Context) (models.CreateUserRow, error) {
	args, ok := ctx.Value(ctxCreateUser).(models.CreateUserParams)
	if !ok {
		return models.CreateUserRow{}, constants.ErrSystemError
	}

	hashedPassword, err := tools.Bcrypt(args.Password, constants.BcryptCost)
	if err != nil {
		return models.CreateUserRow{}, err
	}

	args.Password = hashedPassword

	item, err := s.repo.Create(ctx, args)
	if err != nil {
		return models.CreateUserRow{}, err
	}

	return item, nil
}

func (s service) Update(ctx context.Context) (models.UpdateUserRow, error) {
	args, ok := ctx.Value(ctxUpdateUser).(models.UpdateUserParams)
	if !ok {
		return models.UpdateUserRow{}, constants.ErrSystemError
	}

	if args.Password != "" {
		hashedPassword, err := tools.Bcrypt(args.Password, constants.BcryptCost)
		if err != nil {
			return models.UpdateUserRow{}, err
		}

		args.Password = hashedPassword
	}

	item, err := s.repo.Update(ctx, args)
	if err != nil {
		return models.UpdateUserRow{}, err
	}

	return item, nil
}

func (s service) Delete(ctx context.Context) (models.SoftDeleteUserRow, error) {
	id, ok := ctx.Value(ctxID).(uuid.UUID)
	if !ok {
		return models.SoftDeleteUserRow{}, constants.ErrSystemError
	}

	item, err := s.repo.Delete(ctx, id)
	if err != nil {
		return models.SoftDeleteUserRow{}, err
	}

	return item, nil
}
