package user

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/config"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/constants"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/models"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/test"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	tools "github.com/hinccvi/Golang-Project-Structure-Conventional/tools/hash"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	cfg, _ := config.Load("local")

	db := test.DB(t, &cfg)
	test.Reset(t, db)
	repo := NewRepository(db, logger)

	ctx := context.TODO()

	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	// create
	createdUser, err := repo.Create(ctx, models.CreateUserParams{
		Username: "user",
		Password: "secret",
	})
	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	count2, _ := repo.Count(ctx)
	assert.Equal(t, int64(1), count2-count)

	// get
	user, err := repo.Get(ctx, createdUser.ID)
	assert.Nil(t, err)
	assert.Equal(t, "user", user.Username)
	_, err = repo.Get(ctx, uuid.UUID{})
	assert.Equal(t, pgx.ErrNoRows, err)

	// update
	updatedUser, err := repo.Update(ctx, models.UpdateUserParams{
		ID:       user.ID,
		Username: "testuser",
		Password: "newsecret",
	})
	assert.Nil(t, err)
	user, _ = repo.Get(ctx, updatedUser.ID)
	assert.Equal(t, "testuser", user.Username)

	// query
	albums, err := repo.Query(ctx, models.ListUserParams{
		Limit:  10,
		Offset: 0,
	})
	assert.Nil(t, err)
	assert.Equal(t, count2, int64(len(albums)))

	// delete
	deletedUser, err := repo.Delete(ctx, user.ID)
	assert.Nil(t, err)
	_, err = repo.Get(ctx, deletedUser.ID)
	assert.Equal(t, pgx.ErrNoRows, err)
	_, err = repo.Delete(ctx, deletedUser.ID)
	assert.Equal(t, pgx.ErrNoRows, err)

	// test user
	hashedPassword, err := tools.Bcrypt("secret", constants.BcryptCost)
	assert.Nil(t, err)

	testUser, err := repo.Create(ctx, models.CreateUserParams{
		Username: "testuser1",
		Password: hashedPassword,
	})
	assert.Nil(t, err)
	assert.NotNil(t, testUser)
}
