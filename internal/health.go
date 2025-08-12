package internal

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

// RegisterHealthRoutes registers health check endpoints
func RegisterHealthRoutes(router *gin.Engine, db *gorm.DB, redisClient *redis.Client) {
	router.GET("/health", func(c *gin.Context) {
		health := HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now(),
			Services:  make(map[string]string),
		}

		// Check database health
		if err := checkDatabaseHealth(db); err != nil {
			health.Status = "unhealthy"
			health.Services["database"] = "unhealthy: " + err.Error()
		} else {
			health.Services["database"] = "healthy"
		}

		// Check Redis health
		if err := checkRedisHealth(redisClient); err != nil {
			health.Status = "unhealthy"
			health.Services["redis"] = "unhealthy: " + err.Error()
		} else {
			health.Services["redis"] = "healthy"
		}

		// Set appropriate HTTP status code
		if health.Status == "healthy" {
			c.JSON(http.StatusOK, health)
		} else {
			c.JSON(http.StatusServiceUnavailable, health)
		}
	})

	// Simple health check endpoint
	router.GET("/health/simple", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

func checkDatabaseHealth(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx)
}

func checkRedisHealth(redisClient *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := redisClient.Ping(ctx).Result()
	return err
}
