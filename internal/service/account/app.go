package account

import (
	"go-server/internal/infrastructure"
	delivery "go-server/internal/service/account/delivery"
	"go-server/pkg/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAccountService(r *gin.Engine, db *gorm.DB, httpClient *http.Client, emailService infrastructure.EmailService, cfg *config.Config) {
	controller := delivery.NewAccountController(db, httpClient, emailService, cfg)
	router := delivery.NewAccountRouter(r)

	router.CreateAccount(controller.CreateAccount)
	router.VerifyEmail(controller.VerifyEmail)
	router.ResendVerificationEmail(controller.ResendVerificationEmail)
	router.DeleteAccount(controller.DeleteAccount)
}
