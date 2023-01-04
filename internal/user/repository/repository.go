package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hinccvi/go-ddd/internal/entity"
	"github.com/hinccvi/go-ddd/pkg/log"
	"github.com/jmoiron/sqlx"
)

type (
	// Repository encapsulates the logic to access users from the data source.
	Repository interface {
		Get(ctx context.Context, id uuid.UUID) (entity.User, error)
		Count(ctx context.Context) (int64, error)
		Query(ctx context.Context, page, size int) ([]entity.User, error)
		Create(ctx context.Context, u entity.User) error
		Update(ctx context.Context, u entity.User) error
		Delete(ctx context.Context, id uuid.UUID) error
	}
	// repository persists albums in database.
	repository struct {
		db     *sqlx.DB
		logger log.Logger
	}
)

//nolint:gosec //false positive
const (
	getUser             string = `SELECT id, username FROM "user" WHERE id = $1 AND deleted_at IS NULL LIMIT 1`
	countUser           string = `SELECT COUNT(id) FROM "user"`
	queryUser           string = `SELECT id, username FROM "user" ORDER BY username LIMIT($1) OFFSET($2)`
	createUser          string = `INSERT INTO "user" (username, password) VALUES (:username, :password)`
	updateUserUsername  string = `UPDATE "user" SET username = VARCHAR(:username)`
	updateUserPassword  string = `, password = VARCHAR(:password)`
	updateUserCondition string = ` WHERE id = UUID(:id)`
	deleteUser          string = `UPDATE "user" 
                       SET deleted_at = (current_timestamp AT TIME ZONE 'UTC') 
                       WHERE id = $1 AND deleted_at IS NULL`
)

func New(db *sqlx.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) Get(ctx context.Context, id uuid.UUID) (entity.User, error) {
	getUserStmt, err := r.db.PreparexContext(ctx, getUser)
	if err != nil {
		return entity.User{}, err
	}
	defer getUserStmt.Close()

	var user entity.User
	if err = getUserStmt.GetContext(ctx, &user, id); err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r repository) Count(ctx context.Context) (int64, error) {
	var total int64
	if err := r.db.GetContext(ctx, &total, countUser); err != nil {
		return 0, err
	}

	return total, nil
}

func (r repository) Query(ctx context.Context, page, size int) ([]entity.User, error) {
	queryUserStmt, err := r.db.PreparexContext(ctx, queryUser)
	if err != nil {
		return []entity.User{}, err
	}
	defer queryUserStmt.Close()

	var users []entity.User
	if err = queryUserStmt.SelectContext(ctx, &users, page, size); err != nil {
		return []entity.User{}, err
	}

	return users, nil
}

func (r repository) Create(ctx context.Context, u entity.User) error {
	createUserStmt, err := r.db.PrepareNamedContext(ctx, createUser)
	if err != nil {
		return err
	}
	defer createUserStmt.Close()

	if _, err = createUserStmt.ExecContext(ctx, u); err != nil {
		return err
	}

	return nil
}

func (r repository) Update(ctx context.Context, u entity.User) error {
	query := updateUserUsername
	if u.Password != "" {
		query += updateUserPassword
	}
	query += updateUserCondition

	updateUserStmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer updateUserStmt.Close()

	if _, err = updateUserStmt.ExecContext(ctx, u); err != nil {
		return err
	}

	return nil
}

func (r repository) Delete(ctx context.Context, id uuid.UUID) error {
	deleteUserStmt, err := r.db.PreparexContext(ctx, deleteUser)
	if err != nil {
		return err
	}
	defer deleteUserStmt.Close()

	if _, err = deleteUserStmt.ExecContext(ctx, &id); err != nil {
		return err
	}

	return nil
}
