package repository

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hinccvi/go-ddd/internal/entity"
	"github.com/hinccvi/go-ddd/pkg/log"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var errConnectionRefused = errors.New("connection refused")

func TestGetUserByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	query := `-- name: GetByUsername :one
	SELECT id, username, password FROM "user" WHERE username = $1 AND deleted_at IS NULL LIMIT 1`

	t.Run("success", func(t *testing.T) {
		id := uuid.NewString()
		username := "user"
		password := "secret"

		rows := sqlmock.NewRows([]string{"id", "username", "password"}).
			AddRow(id, username, password)

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(username).WillReturnRows(rows)

		repo := New(db, logger)

		var user entity.GetByUsernameRow
		user, err = repo.GetUserByUsername(context.TODO(), username)
		assert.NoError(t, err)
		assert.Equal(t, id, user.ID.String())
		assert.Equal(t, username, user.Username)
		assert.Equal(t, password, user.Password)
	})

	t.Run("fail: not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("xxx").WillReturnError(sql.ErrNoRows)

		repo := New(db, logger)
		_, err = repo.GetUserByUsername(context.TODO(), "xxx")
		assert.Error(t, err)
	})

	t.Run("fail: db down", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("xxx").WillReturnError(errConnectionRefused)

		repo := New(db, logger)
		_, err = repo.GetUserByUsername(context.TODO(), "xxx")
		assert.Error(t, err)
	})
}
