package internal

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

// RegisterHealthRoutes registers health check endpoints
func RegisterHealthRoutes(router *gin.Engine, db *sql.DB, redisClient *redis.Client) {
	// Startup health check endpoint (no external dependencies)
	router.GET("/health/startup", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "starting",
			"timestamp": time.Now(),
			"message":   "Server is starting up",
		})
	})

	// Main health check endpoint (includes all service checks)
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

		// Check Redis health (if available)
		if redisClient != nil {
			if err := checkRedisHealth(redisClient); err != nil {
				health.Status = "unhealthy"
				health.Services["redis"] = "unhealthy: " + err.Error()
			} else {
				health.Services["redis"] = "healthy"
			}
		} else {
			health.Services["redis"] = "not_configured"
		}

		// Set appropriate HTTP status code
		if health.Status == "healthy" {
			c.JSON(http.StatusOK, health)
		} else {
			c.JSON(http.StatusServiceUnavailable, health)
		}
	})

	// Simple health check endpoint (for basic connectivity)
	router.GET("/health/simple", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Liveness probe endpoint (for Kubernetes/DO App Platform)
	router.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	})

	// Readiness probe endpoint (for Kubernetes/DO App Platform)
	router.GET("/health/ready", func(c *gin.Context) {
		// Check if the application is ready to serve traffic
		health := HealthResponse{
			Status:    "ready",
			Timestamp: time.Now(),
			Services:  make(map[string]string),
		}

		// Check database readiness
		if err := checkDatabaseHealth(db); err != nil {
			health.Status = "not_ready"
			health.Services["database"] = "not_ready: " + err.Error()
		} else {
			health.Services["database"] = "ready"
		}

		// Check Redis readiness (if available)
		if redisClient != nil {
			if err := checkRedisHealth(redisClient); err != nil {
				health.Status = "not_ready"
				health.Services["redis"] = "not_ready: " + err.Error()
			} else {
				health.Services["redis"] = "ready"
			}
		} else {
			health.Services["redis"] = "not_configured"
		}

		// Set appropriate HTTP status code
		if health.Status == "ready" {
			c.JSON(http.StatusOK, health)
		} else {
			c.JSON(http.StatusServiceUnavailable, health)
		}
	})
}

func checkDatabaseHealth(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return db.PingContext(ctx)
}

func checkRedisHealth(redisClient *redis.Client) error {
	if redisClient == nil {
		return fmt.Errorf("Redis client not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := redisClient.Ping(ctx).Result()
	return err
}
