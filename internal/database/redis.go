package database

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(addr, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Println("Redis connection established")
	return client, nil
}
