package redis

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisClient interface {
	SetKey(ctx context.Context, key, value string, expiration time.Duration) error
}

func NewClient(redisAddr string) RedisClient {
	log.Printf("Connecting to redis at: %s", redisAddr)
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	return &redisWrapper{client}
}

type redisWrapper struct {
	client *redis.Client
}

func (r *redisWrapper) SetKey(ctx context.Context, key, value string, expiration time.Duration) error {
	log.Printf("Storing key: %s, value: %s", key, value)
	err := r.client.Set(ctx, key, value, expiration).Err()
	return err
}
