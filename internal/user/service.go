package user

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/entity"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
)

// Service encapsulates usecase logic for user.
type Service interface {
	Get(ctx context.Context, req getOrDeleteUserRequest) (User, error)
	Query(ctx context.Context, req queryUserRequest) ([]User, error)
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, req createUserRequest) (User, error)
	Update(ctx context.Context, req updateUserRequest) (User, error)
	Delete(ctx context.Context, req getOrDeleteUserRequest) (User, error)
}

// User represents the data about a user.
type User struct {
	entity.User
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

func (s service) Get(ctx context.Context, req getOrDeleteUserRequest) (User, error) {
	item, err := s.repo.Get(ctx, req.Id)
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

func (s service) Create(ctx context.Context, req createUserRequest) (User, error) {
	id := uuid.NewString()
	err := s.repo.Create(ctx, entity.User{
		ID:       id,
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		return User{}, err
	}

	return s.Get(ctx, getOrDeleteUserRequest{id})
}

func (s service) Update(ctx context.Context, req updateUserRequest) (User, error) {
	user, err := s.Get(ctx, getOrDeleteUserRequest{req.Id})
	if err != nil {
		return User{}, err
	}

	user.Name = req.Name
	user.Password = req.Password

	if err = s.repo.Update(ctx, user.User); err != nil {
		return user, err
	}

	return user, nil
}

func (s service) Delete(ctx context.Context, req getOrDeleteUserRequest) (User, error) {
	user, err := s.Get(ctx, req)
	if err != nil {
		return User{}, err
	}

	if err = s.repo.Delete(ctx, user.User); err != nil {
		return User{}, err
	}

	return user, nil
}
