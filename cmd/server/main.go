package main

import (
	"context"
	"database/sql"
	"fmt"
	"go-server/internal"
	"go-server/internal/infrastructure"
	"go-server/internal/service/account"
	"go-server/internal/service/session"
	"go-server/internal/service/workspace"
	"go-server/pkg/config"
	"go-server/pkg/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

// retryConnection attempts to establish a connection with retries
func retryConnection(name string, maxRetries int, delay time.Duration, connectFunc func() error) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := connectFunc(); err != nil {
			lastErr = err
			log.Printf("Failed to connect to %s (attempt %d/%d): %v", name, i+1, maxRetries, err)
			if i < maxRetries-1 {
				log.Printf("Retrying in %v...", delay)
				time.Sleep(delay)
				delay *= 2 // Exponential backoff
			}
		} else {
			log.Printf("Successfully connected to %s", name)
			return nil
		}
	}
	return fmt.Errorf("failed to connect to %s after %d attempts: %v", name, maxRetries, lastErr)
}

func main() {
	// Load configuration from environment variables
	cfg := config.Load()

	// Setup logger based on configuration
	logger.SetupLogger(cfg)

	log.Println("Starting Workluv server...")

	// Step 1: Establish a connection to the database with retries
	log.Println("Connecting to database...")
	var db *sql.DB
	err := retryConnection("database", 5, 5*time.Second, func() error {
		var dbErr error
		db, dbErr = infrastructure.ConnectDB(cfg)
		return dbErr
	})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()
	log.Println("Database connection established")

	// Step 2: Try to establish a connection to Redis (optional)
	log.Println("Attempting to connect to Redis...")
	var redisClient *redis.Client
	redisErr := retryConnection("Redis", 3, 3*time.Second, func() error {
		var redisErr error
		redisClient, redisErr = infrastructure.ConnectRedis(cfg)
		return redisErr
	})
	if redisErr != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", redisErr)
		log.Println("Continuing without Redis...")
		redisClient = nil
	} else {
		defer infrastructure.CloseRedis(redisClient)
		log.Println("Redis connection established")
	}

	// Step 3: Initialize a default HTTP client for making requests to external services.
	httpClient := infrastructure.NewDefaultHTTPClient()

	// Step 4: Initialize email service
	log.Println("Initializing email service...")
	emailService, err := infrastructure.NewEmailService(cfg)
	if err != nil {
		log.Printf("Warning: Failed to initialize email service: %v", err)
		log.Println("Continuing without email service...")
	} else {
		log.Printf("Email service initialized: enabled=%t", emailService.IsEnabled())
	}

	// Step 5: Set up the Gin HTTP router, which will handle all incoming HTTP requests.
	router := infrastructure.SetupRouter(cfg)

	// Step 6: Register all application services and their respective routes with the router.
	// This includes the account service, session service, and workspace service, both of which may depend on
	// the database connection, email service, and the HTTP client for their operations.
	log.Println("Registering health routes...")
	internal.RegisterHealthRoutes(router, db, redisClient)

	log.Println("Registering account service...")
	account.RegisterAccountService(router, db, httpClient, emailService, cfg)

	log.Println("Registering session service...")
	session.RegisterSessionService(router, db, cfg)

	log.Println("Registering workspace service...")
	workspace.RegisterWorkspaceService(router, db, httpClient)

	// Step 7: Create HTTP server with graceful shutdown
	serverAddr := ":" + cfg.Server.Port
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server is starting on port %s", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start the HTTP server: %v", err)
		}
	}()

	// Wait a moment for the server to start
	time.Sleep(2 * time.Second)
	log.Println("Server is ready and listening")

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
