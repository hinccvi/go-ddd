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

func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	query := `-- name: GetUser :one
	SELECT id, username FROM "user" WHERE id = $1 AND deleted_at IS NULL LIMIT 1
	`

	t.Run("success", func(t *testing.T) {
		id, _ := uuid.NewRandom()
		username := "user"
		rows := sqlmock.NewRows([]string{"id", "username"}).
			AddRow(id.String(), username)

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnRows(rows)

		repo := New(db, logger)

		var user entity.GetUserRow
		user, err = repo.Get(context.TODO(), id)
		assert.NoError(t, err)
		assert.EqualValues(t, id, user.ID)
		assert.Equal(t, username, user.Username)
	})

	t.Run("fail: not found", func(t *testing.T) {
		var id uuid.UUID
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnError(sql.ErrNoRows)

		repo := New(db, logger)
		_, err = repo.Get(context.TODO(), id)
		assert.Error(t, err)
	})

	t.Run("fail: db down", func(t *testing.T) {
		var id uuid.UUID
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnError(errConnectionRefused)

		repo := New(db, logger)
		_, err = repo.Get(context.TODO(), id)
		assert.Error(t, err)
	})
}

func TestCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	query := `-- name: CountUser :one
	SELECT COUNT(id) FROM "user"
	`

	t.Run("success", func(t *testing.T) {
		var expectedTotal int64 = 1
		rows := sqlmock.NewRows([]string{"total"}).
			AddRow(expectedTotal)

		mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

		repo := New(db, logger)

		var total int64
		total, err = repo.Count(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, expectedTotal, total)
	})

	t.Run("success: empty table", func(t *testing.T) {
		var expectedTotal int64
		rows := sqlmock.NewRows([]string{"total"}).
			AddRow(expectedTotal)

		mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

		repo := New(db, logger)

		var total int64
		total, err = repo.Count(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, expectedTotal, total)
	})

	t.Run("fail: db down", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(errConnectionRefused)

		repo := New(db, logger)
		_, err = repo.Count(context.TODO())
		assert.Error(t, err)
	})
}

func TestQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	query := `-- name: ListUser :many
	SELECT id, username FROM "user" ORDER BY username LIMIT($1) OFFSET($2)
	`

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "username"}).
			AddRow(uuid.NewString(), "user1").
			AddRow(uuid.NewString(), "user2").
			AddRow(uuid.NewString(), "user3")

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(10, 0).WillReturnRows(rows)

		repo := New(db, logger)

		var users []entity.ListUserRow
		users, err = repo.Query(context.TODO(), entity.ListUserParams{Limit: 10, Offset: 0})
		assert.NoError(t, err)
		assert.Len(t, users, 3)
	})

	t.Run("fail: db down", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(10, 0).WillReturnError(errConnectionRefused)
		repo := New(db, logger)

		_, err = repo.Query(context.TODO(), entity.ListUserParams{Limit: 10, Offset: 0})
		assert.Error(t, err)
	})
}

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	query := `-- name: CreateUser :one
	INSERT INTO "user" (username, password) VALUES ($1, $2) RETURNING id, username
	`

	args := entity.CreateUserParams{
		Username: "user",
		Password: "secret",
	}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "username"}).
			AddRow(uuid.NewString(), "user")

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args.Username, args.Password).WillReturnRows(rows)
		repo := New(db, logger)

		var user entity.CreateUserRow
		user, err = repo.Create(context.TODO(), args)
		assert.NoError(t, err)
		assert.Equal(t, args.Username, user.Username)
	})

	t.Run("fail: db down", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args.Username, args.Password).WillReturnError(errConnectionRefused)
		repo := New(db, logger)

		_, err = repo.Create(context.TODO(), args)
		assert.Error(t, err)
	})
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	query := `-- name: UpdateUser :one
	UPDATE "user"
	SET username = CASE WHEN $2::VARCHAR <> ''
				   THEN $2::VARCHAR
				   ELSE username 
				   END,
		password = CASE WHEN $3::VARCHAR <> ''
				   THEN $3::VARCHAR
				   ELSE password 
				   END
	WHERE id = $1
	RETURNING id, username
	`

	args := entity.UpdateUserParams{
		ID:       uuid.New(),
		Username: "user",
	}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "username"}).
			AddRow(args.ID.String(), args.Username)

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args.ID, args.Username, args.Password).WillReturnRows(rows)
		repo := New(db, logger)

		var user entity.UpdateUserRow
		user, err = repo.Update(context.TODO(), args)
		assert.NoError(t, err)
		assert.Equal(t, args.ID, user.ID)
		assert.Equal(t, args.Username, user.Username)
	})

	t.Run("fail: not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(args.ID, args.Username, args.Password).
			WillReturnError(sql.ErrNoRows)
		repo := New(db, logger)

		_, err = repo.Update(context.TODO(), args)
		assert.Error(t, err)
	})

	t.Run("fail: db down", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(args.ID, args.Username, args.Password).
			WillReturnError(errConnectionRefused)
		repo := New(db, logger)

		_, err = repo.Update(context.TODO(), args)
		assert.Error(t, err)
	})
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	query := `-- name: SoftDeleteUser :one
	UPDATE "user" SET deleted_at = (current_timestamp AT TIME ZONE 'UTC') 
	WHERE id = $1 AND deleted_at IS NULL RETURNING id, username
	`

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		username := "user"
		rows := sqlmock.NewRows([]string{"id", "username"}).
			AddRow(id.String(), username)

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnRows(rows)
		repo := New(db, logger)

		var user entity.SoftDeleteUserRow
		user, err = repo.Delete(context.TODO(), id)
		assert.NoError(t, err)
		assert.Equal(t, id, user.ID)
		assert.Equal(t, username, user.Username)
	})

	t.Run("fail: not found", func(t *testing.T) {
		id := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnError(sql.ErrNoRows)
		repo := New(db, logger)

		_, err = repo.Delete(context.TODO(), id)
		assert.Error(t, err)
	})

	t.Run("fail: db down", func(t *testing.T) {
		id := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnError(errConnectionRefused)
		repo := New(db, logger)

		_, err = repo.Delete(context.TODO(), id)
		assert.Error(t, err)
	})
}
