package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/models"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
)

// Repository encapsulates the logic to access users from the data source.
type Repository interface {
	Get(ctx context.Context, id *uuid.UUID) (models.User, error)
	Count(ctx context.Context) (int64, error)
	Query(ctx context.Context, arg *models.ListUserParams) ([]models.User, error)
	Create(ctx context.Context, arg *models.CreateUserParams) (models.CreateUserRow, error)
	Update(ctx context.Context, arg *models.UpdateUserParams) (models.UpdateUserRow, error)
	Delete(ctx context.Context, id *uuid.UUID) (models.User, error)
}

// repository persists albums in database
type repository struct {
	db     *models.DBTX
	logger log.Logger
}

func NewRepository(db *models.DBTX, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) Get(ctx context.Context, id *uuid.UUID) (models.User, error) {
	queries := models.New(*r.db)

	user, err := queries.GetUser(ctx, *id)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r repository) Count(ctx context.Context) (int64, error) {
	queries := models.New(*r.db)

	count, err := queries.CountUser(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r repository) Query(ctx context.Context, arg *models.ListUserParams) ([]models.User, error) {
	queries := models.New(*r.db)

	users, err := queries.ListUser(ctx, arg)
	if err != nil {
		return make([]models.User, 0), err
	}

	return users, nil
}

func (r repository) Create(ctx context.Context, arg *models.CreateUserParams) (models.CreateUserRow, error) {
	queries := models.New(*r.db)

	user, err := queries.CreateUser(ctx, arg)
	if err != nil {
		return models.CreateUserRow{}, err
	}

	return user, nil
}

func (r repository) Update(ctx context.Context, arg *models.UpdateUserParams) (models.UpdateUserRow, error) {
	queries := models.New(*r.db)

	user, err := queries.UpdateUser(ctx, arg)
	if err != nil {
		return models.UpdateUserRow{}, err
	}

	return user, nil
}

func (r repository) Delete(ctx context.Context, id *uuid.UUID) (models.User, error) {
	queries := models.New(*r.db)

	user, err := queries.DeleteUser(ctx, *id)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
