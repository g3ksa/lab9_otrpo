package internal

import (
	"context"

	"github.com/go-redis/redis"
)

func NewRedisClient(ctx context.Context, addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return rdb
}
