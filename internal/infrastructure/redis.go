package infrastructure

import (
	"context"
	"fmt"
	"time"

	"workluv/pkg/config"

	"github.com/redis/go-redis/v9"
)

// ConnectRedis establishes a connection to Redis
func ConnectRedis(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.Database.Redis.Password,
		DB:       cfg.Database.Redis.Database,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return rdb, nil
}

// CloseRedis closes the Redis connection
func CloseRedis(rdb *redis.Client) error {
	return rdb.Close()
}
