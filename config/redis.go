package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Global variables for Redis connection
var (
	RedisPool *redis.Client
	ctx       = context.Background()
)

func init() {
	// Initialize the Redis client
	RedisPool = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // No password set
		DB:       0,  // Default database
	})

	// Test the connection
	if _, err := RedisPool.Ping(ctx).Result(); err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	fmt.Println("Successfully connected to Redis")
}
