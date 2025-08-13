package infrastructure

import (
	"time"
	"workluv/pkg/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// Configure CORS middleware
	corsConfig := cors.Config{
		AllowOrigins:     cfg.Security.CORSAllowedOrigins,
		AllowMethods:     cfg.Security.CORSAllowedMethods,
		AllowHeaders:     cfg.Security.CORSAllowedHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(corsConfig))

	return router
}
