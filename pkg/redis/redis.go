package redis

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func NewClient(redisAddr string) *redis.Client {
	log.Printf("Connecting to redis at: %s", redisAddr)
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	return client
}

func StorePayload(client *redis.Client, key string, value string) error {
	log.Printf("Storing key: %s, value: %s", key, value)
	err := client.Set(ctx, key, value, 0).Err()
	return err
}
