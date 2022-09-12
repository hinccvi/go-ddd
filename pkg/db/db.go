package db

import (
	"context"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/config"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/models"
	"github.com/jackc/pgx/v4"
)

func Connect(mode string, cfg *config.Config) (models.DBTX, error) {
	pgx, err := pgx.Connect(context.Background(), cfg.Dsn)
	if err != nil {
		return nil, err
	}

	return pgx, nil
}
