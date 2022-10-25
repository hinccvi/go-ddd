package auth

import (
	"context"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/model"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
)

type Repository interface {
	GetUserByUsername(ctx context.Context, username string) (model.GetByUsernameRow, error)
}

type repository struct {
	db     model.DBTX
	logger log.Logger
}

func NewRepository(db model.DBTX, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) GetUserByUsername(ctx context.Context, username string) (model.GetByUsernameRow, error) {
	queries := model.New(r.db)

	user, err := queries.GetByUsername(ctx, username)
	if err != nil {
		return model.GetByUsernameRow{}, err
	}

	return user, nil
}
