package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/config"
)

func Connect(cfg config.Config) (*redis.Client, error) {
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

	_, err := rds.Ping(context.Background()).Result()

	return rds, err
}
