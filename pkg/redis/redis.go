package redis

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisClient interface {
	SetKey(ctx context.Context, key, value string, expiration time.Duration) error
}

func NewClient() RedisClient {
	redisAddr := ""

	redisAddr = os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	log.Print("REDIS_HOST: ", os.Getenv("REDIS_HOST"))
	log.Print("REDIS_PORT: ", os.Getenv("REDIS_PORT"))
	log.Print("REDIS_ADDR: ", redisAddr)

	if redisAddr == ":" {
		redisAddr = "localhost:6379" // Default to localhost if environment variables are not set
	} else {
		log.Printf("REDIS_HOST and REDIS_PORT are set, as: %s", redisAddr)
	}

	log.Printf("Connecting to redis at: %s", redisAddr)
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
		return nil
	}
	log.Printf("Connected to redis at: %s", redisAddr)

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
