package account

import (
	delivery "go-server/internal/service/account/delivery"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAccountService(r *gin.Engine, db *gorm.DB, httpClient *http.Client) {
	controller := delivery.NewAccountController(db, httpClient)
	router := delivery.NewAccountRouter(r)

	router.CreateAccount(controller.CreateAccount)
	router.DeleteAccount(controller.DeleteAccount)
}
