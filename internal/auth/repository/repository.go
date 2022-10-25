package repository

import (
	"context"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/entity"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
)

type Repository interface {
	GetUserByUsername(ctx context.Context, username string) (entity.GetByUsernameRow, error)
}

type repository struct {
	db     entity.DBTX
	logger log.Logger
}

func New(db entity.DBTX, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) GetUserByUsername(ctx context.Context, username string) (entity.GetByUsernameRow, error) {
	queries := entity.New(r.db)

	user, err := queries.GetByUsername(ctx, username)
	if err != nil {
		return entity.GetByUsernameRow{}, err
	}

	return user, nil
}
