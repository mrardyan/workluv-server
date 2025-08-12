package session

import (
	delivery "go-server/internal/service/session/delivery"
	"go-server/pkg/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterSessionService registers all session service routes
func RegisterSessionService(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	controller := delivery.NewSessionController(db, cfg)
	router := delivery.NewSessionRouter(r)

	// Register all auth endpoints
	router.Login(controller.Login)
	router.Logout(controller.Logout)
	router.LogoutAll(controller.LogoutAll)
	router.RefreshToken(controller.RefreshToken)
}
