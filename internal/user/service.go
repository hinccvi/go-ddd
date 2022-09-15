package user

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/models"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/tools"
)

// Service encapsulates usecase logic for user.
type Service interface {
	Get(ctx context.Context, id uuid.UUID) (User, error)
	Query(ctx context.Context, arg models.ListUserParams) ([]User, error)
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, arg models.CreateUserParams) (User, error)
	Update(ctx context.Context, arg models.UpdateUserParams) (User, error)
	Delete(ctx context.Context, id uuid.UUID) (User, error)
}

// User represents the data about a user.
type User struct {
	models.User
}

type service struct {
	rds    redis.Client
	repo   Repository
	logger log.Logger
}

// NewService creates a new user service.
func NewService(rds redis.Client, repo Repository, logger log.Logger) Service {
	return service{rds, repo, logger}
}

func (s service) Get(ctx context.Context, id uuid.UUID) (User, error) {
	item, err := s.repo.Get(ctx, id)
	if err != nil {
		return User{}, err
	}

	return User{item}, nil
}

func (s service) Query(ctx context.Context, arg models.ListUserParams) ([]User, error) {
	items, err := s.repo.Query(ctx, arg)
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

func (s service) Create(ctx context.Context, arg models.CreateUserParams) (u User, err error) {
	arg.Password, err = tools.Bcrypt(arg.Password)
	if err != nil {
		return User{}, err
	}

	ur, err := s.repo.Create(ctx, arg)
	if err != nil {
		return User{}, err
	}

	user := models.User{
		ID:        ur.ID,
		Username:  ur.Username,
		CreatedAt: ur.CreatedAt,
		UpdatedAt: ur.UpdatedAt,
	}

	return User{user}, nil
}

func (s service) Update(ctx context.Context, arg models.UpdateUserParams) (u User, err error) {
	if arg.Password != "" {
		arg.Password, err = tools.Bcrypt(arg.Password)
		if err != nil {
			return User{}, err
		}
	}

	ur, err := s.repo.Update(ctx, arg)
	if err != nil {
		return User{}, err
	}

	user := models.User{
		ID:        ur.ID,
		Username:  ur.Username,
		CreatedAt: ur.CreatedAt,
		UpdatedAt: ur.UpdatedAt,
	}

	return User{user}, nil
}

func (s service) Delete(ctx context.Context, id uuid.UUID) (u User, err error) {
	ur, err := s.repo.Delete(ctx, id)
	if err != nil {
		return User{}, err
	}

	user := models.User{
		ID:        ur.ID,
		Username:  ur.Username,
		CreatedAt: ur.CreatedAt,
		UpdatedAt: ur.UpdatedAt,
		DeletedAt: ur.DeletedAt,
	}

	return User{user}, nil
}
