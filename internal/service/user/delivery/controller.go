package delivery

import (
	"go-server/internal/service/user/domain"
	"go-server/internal/service/user/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	UseCase    *domain.UseCase
	HTTPClient *http.Client
}

func NewUserController(db *gorm.DB, httpClient *http.Client) *Controller {
	repo := repository.NewGormRepository(db)
	useCase := domain.NewUseCase(repo)
	return &Controller{
		UseCase:    useCase,
		HTTPClient: httpClient,
	}
}

func (ctl *Controller) CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := ctl.UseCase.CreateUser(domain.User{
		Username: req.Username,
		Email:    domain.Email(req.Email),
		Password: domain.Password(req.Password),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (ctl *Controller) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = ctl.UseCase.DeleteUser(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
