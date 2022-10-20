package auth

import (
	"context"
	"testing"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/config"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/test"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	cfg, _ := config.Load("local")

	db := test.DB(t, &cfg)
	repo := NewRepository(db, logger)

	ctx := context.TODO()

	// get user
	user, err := repo.GetUserByUsername(ctx, "testuser1")
	assert.Nil(t, err)
	assert.NotNil(t, user)

	_, err = repo.GetUserByUsername(ctx, "xxx")
	assert.NotNil(t, err)
}
