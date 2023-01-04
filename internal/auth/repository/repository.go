package repository

import (
	"context"

	"github.com/hinccvi/go-ddd/internal/entity"
	"github.com/hinccvi/go-ddd/pkg/log"
	"github.com/jmoiron/sqlx"
)

type (
	Repository interface {
		GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	}

	repository struct {
		db     *sqlx.DB
		logger log.Logger
	}
)

const (
	getUserByUsername string = `SELECT id, username, password 
                              FROM "user"
                              WHERE username = $1 AND deleted_at IS NULL 
                              LIMIT 1`
)

func New(db *sqlx.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	getUserStmt, err := r.db.PreparexContext(ctx, getUserByUsername)
	if err != nil {
		return entity.User{}, err
	}
	defer getUserStmt.Close()

	var user entity.User
	if err = getUserStmt.GetContext(ctx, &user, username); err != nil {
		return entity.User{}, err
	}

	return user, nil
}
