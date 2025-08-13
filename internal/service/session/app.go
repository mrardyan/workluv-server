package session

import (
	"database/sql"
	delivery "workluv/internal/service/session/delivery"
	"workluv/pkg/config"

	"github.com/gin-gonic/gin"
)

// RegisterSessionService registers all session service routes
func RegisterSessionService(r *gin.Engine, db *sql.DB, cfg *config.Config) {
	controller := delivery.NewSessionController(db, cfg)
	router := delivery.NewSessionRouter(r)

	// Register all auth endpoints
	router.Login(controller.Login)
	router.Logout(controller.Logout)
	router.LogoutAll(controller.LogoutAll)
	router.RefreshToken(controller.RefreshToken)
}
