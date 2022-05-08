package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func Connect(mode string, cfg config.Config) (*gorm.DB, error) {
	// Disable gorm logging if local environment
	var gLogger gormlogger.Interface
	if mode == "local" {
		gLogger = gormlogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			gormlogger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  gormlogger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		)
	} else {
		gLogger = gormlogger.Default.LogMode(gormlogger.Silent)
	}

	// connect to the database
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Bangkok",
			cfg.DBConfig.Host,
			cfg.DBConfig.User,
			cfg.DBConfig.Password,
			cfg.DBConfig.DBName,
			cfg.DBConfig.Port),
	}), &gorm.Config{
		Logger: gLogger,
	})

	if err != nil {
		return &gorm.DB{}, err
	}

	return db, nil
}
