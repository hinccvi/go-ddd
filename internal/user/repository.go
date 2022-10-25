package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/entity"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
)

// Repository encapsulates the logic to access users from the data source.
type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (entity.GetUserRow, error)
	Count(ctx context.Context) (int64, error)
	Query(ctx context.Context, arg entity.ListUserParams) ([]entity.ListUserRow, error)
	Create(ctx context.Context, arg entity.CreateUserParams) (entity.CreateUserRow, error)
	Update(ctx context.Context, arg entity.UpdateUserParams) (entity.UpdateUserRow, error)
	Delete(ctx context.Context, id uuid.UUID) (entity.SoftDeleteUserRow, error)
}

// repository persists albums in database.
type repository struct {
	db     entity.DBTX
	logger log.Logger
}

func NewRepository(db entity.DBTX, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) Get(ctx context.Context, id uuid.UUID) (entity.GetUserRow, error) {
	queries := entity.New(r.db)

	user, err := queries.GetUser(ctx, id)
	if err != nil {
		return entity.GetUserRow{}, err
	}

	return user, nil
}

func (r repository) Count(ctx context.Context) (int64, error) {
	queries := entity.New(r.db)

	count, err := queries.CountUser(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r repository) Query(ctx context.Context, arg entity.ListUserParams) ([]entity.ListUserRow, error) {
	queries := entity.New(r.db)

	users, err := queries.ListUser(ctx, arg)
	if err != nil {
		return make([]entity.ListUserRow, 0), err
	}

	return users, nil
}

func (r repository) Create(ctx context.Context, arg entity.CreateUserParams) (entity.CreateUserRow, error) {
	queries := entity.New(r.db)

	user, err := queries.CreateUser(ctx, arg)
	if err != nil {
		return entity.CreateUserRow{}, err
	}

	return user, nil
}

func (r repository) Update(ctx context.Context, arg entity.UpdateUserParams) (entity.UpdateUserRow, error) {
	queries := entity.New(r.db)

	user, err := queries.UpdateUser(ctx, arg)
	if err != nil {
		return entity.UpdateUserRow{}, err
	}

	return user, nil
}

func (r repository) Delete(ctx context.Context, id uuid.UUID) (entity.SoftDeleteUserRow, error) {
	queries := entity.New(r.db)

	user, err := queries.SoftDeleteUser(ctx, id)
	if err != nil {
		return entity.SoftDeleteUserRow{}, err
	}

	return user, nil
}
