package user

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/models"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
)

// Service encapsulates usecase logic for user.
type Service interface {
	Get(ctx context.Context, id *uuid.UUID) (User, error)
	Query(ctx context.Context, req queryUserRequest) ([]User, error)
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, arg *models.CreateUserParams) (User, error)
	Update(ctx context.Context, req updateUserRequest) (User, error)
	Delete(ctx context.Context, id *uuid.UUID) (User, error)
}

// User represents the data about a user.
type User struct {
	models.User
}

type service struct {
	rds    *redis.Client
	repo   Repository
	logger log.Logger
}

// NewService creates a new user service.
func NewService(rds *redis.Client, repo Repository, logger log.Logger) Service {
	return service{rds, repo, logger}
}

func (s service) Get(ctx context.Context, id *uuid.UUID) (User, error) {
	item, err := s.repo.Get(ctx, id)
	if err != nil {
		return User{}, err
	}

	return User{item}, nil
}

func (s service) Query(ctx context.Context, req queryUserRequest) ([]User, error) {
	items, err := s.repo.Query(ctx, *req.Offset, req.Limit)
	if err != nil {
		return []User{}, err
	}

	users := []User{}
	for _, v := range items {
		users = append(users, User{v})
	}

	return users, nil
}

func (s service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

func (s service) Create(ctx context.Context, arg *models.CreateUserParams) (User, error) {
	user, err := s.repo.Create(ctx, arg)
	if err != nil {
		return User{}, err
	}

	return User{user}, nil
}

func (s service) Update(ctx context.Context, req updateUserRequest) (User, error) {
	return User{models.User{}}, s.repo.Update(ctx)
}

func (s service) Delete(ctx context.Context, id *uuid.UUID) (User, error) {
	user, err := s.repo.Delete(ctx, id)
	if err != nil {
		return User{}, err
	}

	return User{user}, nil
}
