package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/hinccvi/go-ddd/internal/config"
	// postgres driver required by database/sql.
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

const (
	maxOpenConns   int           = 25
	maxIdleConns   int           = 25
	maxLifetime    time.Duration = 5 * time.Minute
	contextTimeout time.Duration = 5 * time.Second
)

func Connect(ctx context.Context, cfg *config.Config, log *zap.Logger) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(maxLifetime)

	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
