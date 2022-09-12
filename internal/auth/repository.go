package auth

import (
	"context"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/models"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
)

type Repository interface {
	GetUserByUsername(ctx context.Context, username string) (models.GetByUsernameRow, error)
}

type repository struct {
	db     models.DBTX
	logger log.Logger
}

func NewRepository(db models.DBTX, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) GetUserByUsername(ctx context.Context, username string) (models.GetByUsernameRow, error) {
	queries := models.New(r.db)

	user, err := queries.GetByUsername(ctx, username)
	if err != nil {
		return models.GetByUsernameRow{}, err
	}

	return user, nil
}
