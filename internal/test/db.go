package test

import (
	"context"
	"testing"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/config"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
)

func Connect(ctx context.Context, t *testing.T, cfg *config.Config) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(cfg.Dsn)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	logger, _ := log.NewForTest()
	config.ConnConfig.Logger = zapadapter.NewLogger(logger)

	pgx, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	return pgx
}

func Reset(ctx context.Context, t *testing.T, pgx *pgxpool.Pool) {
	sql := "DROP TABLE test"

	_, err := pgx.Exec(ctx, sql)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
