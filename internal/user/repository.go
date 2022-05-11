package user

import (
	"context"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/entity"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	"gorm.io/gorm"
)

// Repository encapsulates the logic to access users from the data source.
type Repository interface {
	Get(ctx context.Context, id string) (entity.User, error)
	Count(ctx context.Context) (int64, error)
	Query(ctx context.Context, offset, limit int) ([]entity.User, error)
	Create(ctx context.Context, user entity.User) error
	Update(ctx context.Context, user entity.User) error
	Delete(ctx context.Context, user entity.User) error
}

// repository persists albums in database
type repository struct {
	db     *gorm.DB
	logger log.Logger
}

func NewRepository(db *gorm.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) Get(ctx context.Context, id string) (entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).First(&user, "id=?", id).Error
	return user, err
}

func (r repository) Count(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&entity.User{}).Count(&total).Error
	return total, err
}

func (r repository) Query(ctx context.Context, offset, limit int) ([]entity.User, error) {
	var users []entity.User
	err := r.db.WithContext(ctx).
		Order("created_at desc").
		Offset((offset) * limit).
		Limit(limit).
		Find(&users).
		Error

	return users, err
}

func (r repository) Create(ctx context.Context, user entity.User) error {
	return r.db.WithContext(ctx).Create(&user).Error
}

func (r repository) Update(ctx context.Context, user entity.User) error {
	return r.db.WithContext(ctx).Updates(&user).Error
}

func (r repository) Delete(ctx context.Context, user entity.User) error {
	return r.db.WithContext(ctx).Delete(&user).Error
}
