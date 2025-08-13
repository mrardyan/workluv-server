package delivery

import (
	"database/sql"
	"net/http"
	"workluv/internal/infrastructure"
	"workluv/internal/service/account/domain"
	"workluv/internal/service/account/dto"
	"workluv/internal/service/account/repository"
	"workluv/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Controller struct {
	UseCase    *domain.UseCase
	HTTPClient *http.Client
}

func NewAccountController(db *sql.DB, httpClient *http.Client, emailService infrastructure.EmailService, cfg *config.Config) *Controller {
	repo := repository.NewRepository(db)
	useCase := domain.NewUseCase(repo, emailService, cfg)
	return &Controller{
		UseCase:    useCase,
		HTTPClient: httpClient,
	}
}

func (ctl *Controller) CreateAccount(c *gin.Context) {
	var req dto.CreateAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
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

	createdAccount, err := ctl.UseCase.CreateAccount(c.Request.Context(), account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.CreateAccountResponse{
		ID:            createdAccount.ID.String(),
		Email:         string(createdAccount.Email),
		FullName:      createdAccount.FullName,
		IsActive:      createdAccount.IsActive,
		EmailVerified: createdAccount.EmailVerified,
		CreatedAt:     createdAccount.CreatedAt,
		Message:       "Account created successfully. Please check your email to verify your account.",
	}

	c.JSON(http.StatusCreated, response)
}

func (ctl *Controller) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		// Try to get token from request body
		var req dto.VerifyEmailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token is required"})
			return
		}
		token = req.Token
	}

	err := ctl.UseCase.VerifyEmail(c.Request.Context(), token)
	if err != nil {
		response := dto.EmailVerificationResponse{
			Success: false,
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := dto.EmailVerificationResponse{
		Success: true,
		Message: "Email verified successfully! You can now log in.",
	}
	c.JSON(http.StatusOK, response)
}

func (ctl *Controller) ResendVerificationEmail(c *gin.Context) {
	var req dto.ResendVerificationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	email := domain.Email(req.Email)
	if !email.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	err := ctl.UseCase.ResendVerificationEmail(c.Request.Context(), email)
	if err != nil {
		response := dto.EmailVerificationResponse{
			Success: false,
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := dto.EmailVerificationResponse{
		Success: true,
		Message: "Verification email sent successfully. Please check your email.",
	}
	c.JSON(http.StatusOK, response)
}

func (ctl *Controller) DeleteAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID format"})
		return
	}

	err = ctl.UseCase.DeleteAccount(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}
