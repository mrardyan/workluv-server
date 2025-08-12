package infrastructure

import (
	"context"
	"fmt"
	"time"

	"go-server/pkg/config"

	"github.com/redis/go-redis/v9"
)

// ConnectRedis establishes a connection to Redis
func ConnectRedis(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
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
