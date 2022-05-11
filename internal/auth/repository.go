package auth

import (
	"context"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/entity"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	"gorm.io/gorm"
)

type Repository interface {
	GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entity.User, error)
}

type repository struct {
	db     *gorm.DB
	logger log.Logger
}

func NewRepository(db *gorm.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Select("id, username").Where(&user).First(&user).Error
	return user, err
}
