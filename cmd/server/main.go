package main

import (
	"context"
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
)

func main() {
	// Load configuration from environment variables
	cfg := config.Load()

	// Setup logger based on configuration
	logger.SetupLogger(cfg)

	// Step 1: Establish a connection to the database using the provided configuration.
	db, err := infrastructure.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Step 2: Establish a connection to Redis
	redisClient, err := infrastructure.ConnectRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer infrastructure.CloseRedis(redisClient)

	// Step 3: Initialize a default HTTP client for making requests to external services.
	httpClient := infrastructure.NewDefaultHTTPClient()

	// Step 4: Initialize email service
	emailService, err := infrastructure.NewEmailService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize email service: %v", err)
	}
	log.Printf("Email service initialized: enabled=%t", emailService.IsEnabled())

	// Step 5: Set up the Gin HTTP router, which will handle all incoming HTTP requests.
	router := infrastructure.SetupRouter(cfg)

	// Step 6: Register all application services and their respective routes with the router.
	// This includes the account service, session service, and workspace service, both of which may depend on
	// the database connection, email service, and the HTTP client for their operations.
	internal.RegisterHealthRoutes(router, db, redisClient)
	account.RegisterAccountService(router, db, httpClient, emailService, cfg)
	session.RegisterSessionService(router, db, cfg)
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
