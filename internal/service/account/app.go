package account

import (
	"database/sql"
	"net/http"
	"workluv/internal/infrastructure"
	delivery "workluv/internal/service/account/delivery"
	"workluv/pkg/config"

	"github.com/gin-gonic/gin"
)

func RegisterAccountService(r *gin.Engine, db *sql.DB, httpClient *http.Client, emailService infrastructure.EmailService, cfg *config.Config) {
	controller := delivery.NewAccountController(db, httpClient, emailService, cfg)
	router := delivery.NewAccountRouter(r)

	router.CreateAccount(controller.CreateAccount)
	router.VerifyEmail(controller.VerifyEmail)
	router.ResendVerificationEmail(controller.ResendVerificationEmail)
	router.DeleteAccount(controller.DeleteAccount)
}
