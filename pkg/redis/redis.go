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
				cfg.RedisConfig.Host,
				cfg.RedisConfig.Port,
			),
			Password: cfg.RedisConfig.Password,
			DB:       cfg.RedisConfig.DB,
			PoolSize: cfg.RedisConfig.PoolSize,
		})

	_, err := rds.Ping(context.Background()).Result()

	return rds, err
}
