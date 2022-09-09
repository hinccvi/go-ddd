package auth

import (
	"context"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/models"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
)

type Repository interface {
	GetUserByUsernameAndPassword(ctx context.Context, arg *models.GetByUsernameAndPasswordParams) (models.User, error)
}

type repository struct {
	db     *models.DBTX
	logger log.Logger
}

func NewRepository(db *models.DBTX, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) GetUserByUsernameAndPassword(ctx context.Context, arg *models.GetByUsernameAndPasswordParams) (models.User, error) {
	queries := models.New(*r.db)

	user, err := queries.GetByUsernameAndPassword(ctx, arg)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
