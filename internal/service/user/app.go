package user

import (
	delivery "go-server/internal/service/user/delivery"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserService(r *gin.Engine, db *gorm.DB, httpClient *http.Client) {
	controller := delivery.NewUserController(db, httpClient)
	router := delivery.NewUserRouter(r)

	router.CreateUser(controller.CreateUser)
	router.DeleteUser(controller.DeleteUser)
}
