package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/model"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
)

// Repository encapsulates the logic to access users from the data source.
type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (model.GetUserRow, error)
	Count(ctx context.Context) (int64, error)
	Query(ctx context.Context, arg model.ListUserParams) ([]model.ListUserRow, error)
	Create(ctx context.Context, arg model.CreateUserParams) (model.CreateUserRow, error)
	Update(ctx context.Context, arg model.UpdateUserParams) (model.UpdateUserRow, error)
	Delete(ctx context.Context, id uuid.UUID) (model.SoftDeleteUserRow, error)
}

// repository persists albums in database.
type repository struct {
	db     model.DBTX
	logger log.Logger
}

func NewRepository(db model.DBTX, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) Get(ctx context.Context, id uuid.UUID) (model.GetUserRow, error) {
	queries := model.New(r.db)

	user, err := queries.GetUser(ctx, id)
	if err != nil {
		return model.GetUserRow{}, err
	}

	return user, nil
}

func (r repository) Count(ctx context.Context) (int64, error) {
	queries := model.New(r.db)

	count, err := queries.CountUser(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r repository) Query(ctx context.Context, arg model.ListUserParams) ([]model.ListUserRow, error) {
	queries := model.New(r.db)

	users, err := queries.ListUser(ctx, arg)
	if err != nil {
		return make([]model.ListUserRow, 0), err
	}

	return users, nil
}

func (r repository) Create(ctx context.Context, arg model.CreateUserParams) (model.CreateUserRow, error) {
	queries := model.New(r.db)

	user, err := queries.CreateUser(ctx, arg)
	if err != nil {
		return model.CreateUserRow{}, err
	}

	return user, nil
}

func (r repository) Update(ctx context.Context, arg model.UpdateUserParams) (model.UpdateUserRow, error) {
	queries := model.New(r.db)

	user, err := queries.UpdateUser(ctx, arg)
	if err != nil {
		return model.UpdateUserRow{}, err
	}

	return user, nil
}

func (r repository) Delete(ctx context.Context, id uuid.UUID) (model.SoftDeleteUserRow, error) {
	queries := model.New(r.db)

	user, err := queries.SoftDeleteUser(ctx, id)
	if err != nil {
		return model.SoftDeleteUserRow{}, err
	}

	return user, nil
}
