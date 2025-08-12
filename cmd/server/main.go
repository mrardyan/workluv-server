package main

import (
	"go-server/internal"
	"go-server/internal/infrastructure"
	"go-server/internal/service/user"
	"go-server/internal/service/workspace"
	"go-server/pkg/config"
	"go-server/pkg/logger"
	"log"
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

	// Step 4: Set up the Gin HTTP router, which will handle all incoming HTTP requests.
	router := infrastructure.SetupRouter()

	// Step 5: Register all application services and their respective routes with the router.
	// This includes the user service and the workspace service, both of which may depend on
	// the database connection and the HTTP client for their operations.
	internal.RegisterHealthRoutes(router, db, redisClient)
	user.RegisterUserService(router, db, httpClient)
	workspace.RegisterWorkspaceService(router, db, httpClient)

	// Step 6: Start the HTTP server on the configured port and listen for incoming requests.
	// If the server fails to start, log the error and terminate the application.
	serverAddr := ":" + cfg.Server.Port
	log.Printf("Server is starting on port %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start the HTTP server: %v", err)
	}
}
