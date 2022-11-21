package mocks

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

func Redis(host string) (redis.Client, error) {
	rds := redis.NewClient(
		&redis.Options{
			Addr: host,
		})

	_, err := rds.Ping(context.TODO()).Result()

	return *rds, err
}

func RefreshTokenKey(field string) string {
	return fmt.Sprintf("%s:%s:%s", "test", "refresh_token", field)
}
