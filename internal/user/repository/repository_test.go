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
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var errConnectionRefused = errors.New("connection refused")

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbx := sqlx.NewDb(db, "pgx")
	defer db.Close()

	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		username := "user"
		rows := sqlmock.NewRows([]string{"id", "username"}).
			AddRow(id.String(), username)

		mock.ExpectPrepare(regexp.QuoteMeta(getUser)).ExpectQuery().WithArgs(id).WillReturnRows(rows)
		repo := New(dbx, logger)

		var user entity.User
		user, err = repo.GetUser(context.TODO(), id)
		assert.NoError(t, err)
		assert.EqualValues(t, id, user.ID)
		assert.Equal(t, username, user.Username)
	})

	t.Run("fail: not found", func(t *testing.T) {
		var id uuid.UUID
		mock.ExpectPrepare(regexp.QuoteMeta(getUser)).ExpectQuery().WithArgs(id).WillReturnError(sql.ErrNoRows)

		repo := New(dbx, logger)
		_, err = repo.GetUser(context.TODO(), id)
		assert.Error(t, err)
	})

	t.Run("fail: db down", func(t *testing.T) {
		var id uuid.UUID
		mock.ExpectPrepare(regexp.QuoteMeta(getUser)).ExpectQuery().WithArgs(id).WillReturnError(errConnectionRefused)
		repo := New(dbx, logger)
		_, err = repo.GetUser(context.TODO(), id)
		assert.Error(t, err)
	})
}

func TestCountUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbx := sqlx.NewDb(db, "pgx")
	defer db.Close()

	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	t.Run("success", func(t *testing.T) {
		var expectedTotal int64 = 1
		rows := sqlmock.NewRows([]string{"total"}).
			AddRow(expectedTotal)

		mock.ExpectQuery(regexp.QuoteMeta(countUser)).WillReturnRows(rows)

		repo := New(dbx, logger)
		var total int64
		total, err = repo.CountUser(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, expectedTotal, total)
	})

	t.Run("success: empty table", func(t *testing.T) {
		var expectedTotal int64
		rows := sqlmock.NewRows([]string{"total"}).
			AddRow(expectedTotal)

		mock.ExpectQuery(regexp.QuoteMeta(countUser)).WillReturnRows(rows)

		repo := New(dbx, logger)
		var total int64
		total, err = repo.CountUser(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, expectedTotal, total)
	})

	t.Run("fail: db down", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(countUser)).WillReturnError(errConnectionRefused)

		repo := New(dbx, logger)
		_, err = repo.CountUser(context.TODO())
		assert.Error(t, err)
	})
}

func TestQueryUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbx := sqlx.NewDb(db, "pgx")
	defer db.Close()

	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "username"}).
			AddRow(uuid.NewString(), "user1").
			AddRow(uuid.NewString(), "user2").
			AddRow(uuid.NewString(), "user3")

		mock.ExpectPrepare(regexp.QuoteMeta(queryUser)).ExpectQuery().WithArgs(10, 0).WillReturnRows(rows)

		repo := New(dbx, logger)
		var users []entity.User
		users, err = repo.QueryUser(context.TODO(), 1, 10)
		assert.NoError(t, err)
		assert.Len(t, users, 3)
	})

	t.Run("fail: db down", func(t *testing.T) {
		mock.ExpectPrepare(regexp.QuoteMeta(queryUser)).ExpectQuery().WithArgs(1, 10).WillReturnError(errConnectionRefused)

		repo := New(dbx, logger)
		_, err = repo.QueryUser(context.TODO(), 1, 10)
		assert.Error(t, err)
	})
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbx := sqlx.NewDb(db, "pgx")
	defer db.Close()

	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	t.Run("success", func(t *testing.T) {
		u := entity.User{
			Username: "user",
			Password: "secret",
		}
		mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO "user" (username, password)`)).
			ExpectExec().
			WithArgs(u.Username, u.Password).
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo := New(dbx, logger)
		err = repo.CreateUser(context.TODO(), u)
		assert.NoError(t, err)
	})

	t.Run("fail: db down", func(t *testing.T) {
		u := entity.User{
			Username: "user",
			Password: "secret",
		}
		mock.ExpectPrepare(regexp.QuoteMeta(createUser)).
			ExpectExec().WithArgs(u.Username, u.Password).WillReturnError(errConnectionRefused)

		repo := New(dbx, logger)
		err = repo.CreateUser(context.TODO(), u)
		assert.Error(t, err)
	})
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbx := sqlx.NewDb(db, "pgx")
	defer db.Close()

	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	u := entity.User{
		ID:       uuid.New(),
		Username: "user",
		Password: "secret",
	}
	t.Run("success", func(t *testing.T) {
		mock.ExpectPrepare(regexp.QuoteMeta(`UPDATE "user"`)).
			ExpectExec().
			WithArgs(u.Username, u.Password, u.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo := New(dbx, logger)
		err = repo.UpdateUser(context.TODO(), u)
		assert.NoError(t, err)
	})

	t.Run("fail: not found", func(t *testing.T) {
		mock.ExpectPrepare(regexp.QuoteMeta(`UPDATE "user"`)).ExpectExec().
			WithArgs(u.ID, u.Username, u.Password).
			WillReturnError(sql.ErrNoRows)
		repo := New(dbx, logger)
		err = repo.UpdateUser(context.TODO(), u)
		assert.Error(t, err)
	})

	t.Run("fail: db down", func(t *testing.T) {
		mock.ExpectPrepare(regexp.QuoteMeta(`UPDATE "user"`)).ExpectExec().
			WithArgs(u.ID, u.Username, u.Password).
			WillReturnError(errConnectionRefused)
		repo := New(dbx, logger)
		err = repo.UpdateUser(context.TODO(), u)
		assert.Error(t, err)
	})
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbx := sqlx.NewDb(db, "pgx")
	defer db.Close()

	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	t.Run("success", func(t *testing.T) {
		id := uuid.New()

		mock.ExpectPrepare(regexp.QuoteMeta(deleteUser)).
			ExpectExec().
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		repo := New(dbx, logger)
		err = repo.DeleteUser(context.TODO(), id)
		assert.NoError(t, err)
	})

	t.Run("fail: not found", func(t *testing.T) {
		id := uuid.New()
		mock.ExpectPrepare(regexp.QuoteMeta(deleteUser)).ExpectExec().WithArgs(id).WillReturnError(sql.ErrNoRows)
		repo := New(dbx, logger)
		err = repo.DeleteUser(context.TODO(), id)
		assert.Error(t, err)
	})

	t.Run("fail: db down", func(t *testing.T) {
		id := uuid.New()
		mock.ExpectPrepare(regexp.QuoteMeta(deleteUser)).ExpectExec().WithArgs(id).WillReturnError(errConnectionRefused)
		repo := New(dbx, logger)
		err = repo.DeleteUser(context.TODO(), id)
		assert.Error(t, err)
	})
}
