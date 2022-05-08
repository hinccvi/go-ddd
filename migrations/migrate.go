package migrations

import (
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/entity"
	"gorm.io/gorm"
)

func Init(db *gorm.DB) error {
	err := db.AutoMigrate(
		&entity.User{},
	)

	return err
}
