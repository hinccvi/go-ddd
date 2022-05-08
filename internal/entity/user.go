package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"primaryKey,column:id" json:"id"`
	Name      string         `gorm:"column:name" json:"name"`
	Age       int            `gorm:"column:age" json:"age"`
	Position  string         `gorm:"column:position" json:"position"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (User) TableName() string {
	return "user"
}
