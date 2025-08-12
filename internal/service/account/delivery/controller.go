package delivery

import (
	"go-server/internal/service/account/domain"
	"go-server/internal/service/account/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Controller struct {
	UseCase    *domain.UseCase
	HTTPClient *http.Client
}

func NewAccountController(db *gorm.DB, httpClient *http.Client) *Controller {
	repo := repository.NewGormRepository(db)
	useCase := domain.NewUseCase(repo)
	return &Controller{
		UseCase:    useCase,
		HTTPClient: httpClient,
	}
}

func (ctl *Controller) CreateAccount(c *gin.Context) {
	var req struct {
		FullName string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate password strength
	password := domain.Password(req.Password)
	if !password.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
		return
	}

	// Validate email format
	email := domain.Email(req.Email)
	if !email.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password.String()), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	account := domain.Account{
		Email:         email,
		PasswordHash:  string(hashedPassword),
		FullName:      req.FullName,
		IsActive:      true,
		EmailVerified: false,
	}

	createdAccount, err := ctl.UseCase.CreateAccount(account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdAccount)
}

func (ctl *Controller) DeleteAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID format"})
		return
	}

	err = ctl.UseCase.DeleteAccount(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}
