package db

import (
	"context"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/config"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/models"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

func Connect(cfg *config.Config, log *zap.Logger) (models.DBTX, error) {
	config, err := pgxpool.ParseConfig(cfg.Dsn)
	if err != nil {
		return nil, err
	}

	config.ConnConfig.Logger = zapadapter.NewLogger(log)

	pgx, err := pgxpool.ConnectConfig(context.TODO(), config)
	if err != nil {
		return nil, err
	}

	return pgx, nil
}
