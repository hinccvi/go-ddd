package test

import (
	"context"
	"testing"

	"github.com/hinccvi/go-ddd/internal/config"
	"github.com/hinccvi/go-ddd/pkg/log"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
)

func DB(t *testing.T, cfg *config.Config) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(cfg.Dsn)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	logger, _ := log.NewForTest()
	config.ConnConfig.Logger = zapadapter.NewLogger(logger)

	pgx, err := pgxpool.ConnectConfig(context.TODO(), config)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	return pgx
}

func Reset(t *testing.T, pgx *pgxpool.Pool) {
	sql := `TRUNCATE TABLE "user"`

	_, err := pgx.Exec(context.TODO(), sql)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
