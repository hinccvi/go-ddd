package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
	"github.com/hinccvi/go-ddd/internal/config"
)

func Connect(ctx context.Context, cfg *config.Config) (redis.Client, error) {
	rds := redis.NewClient(
		&redis.Options{
			Addr: fmt.Sprintf("%s:%d",
				cfg.Redis.Host,
				cfg.Redis.Port,
			),
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
			PoolSize: cfg.Redis.PoolSize,
		})

	_, err := rds.Ping(ctx).Result()

	return *rds, err
}
