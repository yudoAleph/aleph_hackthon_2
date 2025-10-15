package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

func PingRedis(client *redis.Client) error {
	return client.Ping(context.Background()).Err()
}
