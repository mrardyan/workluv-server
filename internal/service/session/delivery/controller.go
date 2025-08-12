package delivery

import (
	accountDomain "go-server/internal/service/account/domain"
	accountRepo "go-server/internal/service/account/repository"
	"go-server/internal/service/session/domain"
	"go-server/internal/service/session/dto"
	"go-server/internal/service/session/repository"
	"go-server/pkg/config"
	"go-server/pkg/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Controller struct {
	UseCase *domain.UseCase
}

// NewSessionController creates a new session controller
func NewSessionController(db *gorm.DB, cfg *config.Config) *Controller {
	sessionRepo := repository.NewGormRepository(db)
	accountRepoImpl := accountRepo.NewGormRepository(db)
	jwtService := jwt.NewService(cfg)
	useCase := domain.NewUseCase(sessionRepo, accountRepoImpl, jwtService)

	return &Controller{
		UseCase: useCase,
	}
}

// Login handles user login and creates a session
func (ctl *Controller) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	loginReq := domain.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	response, err := ctl.UseCase.Login(loginReq)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
		}
		return
	}

	// Convert account to DTO
	account, ok := response.Account.(accountDomain.Account)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	loginResponse := dto.LoginResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		ExpiresAt:    response.ExpiresAt,
		Account: dto.AccountInfo{
			ID:            account.ID.String(),
			Email:         string(account.Email),
			FullName:      account.FullName,
			IsActive:      account.IsActive,
			EmailVerified: account.EmailVerified,
		},
	}

	c.JSON(http.StatusOK, loginResponse)
}

// Logout handles user logout and revokes the session
func (ctl *Controller) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := ctl.UseCase.Logout(req.RefreshToken)
	if err != nil {
		switch err {
		case domain.ErrSessionNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		case domain.ErrSessionRevoked:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Session already revoked"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed"})
		}
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Successfully logged out"})
}

// LogoutAll handles logout from all devices and revokes all sessions
func (ctl *Controller) LogoutAll(c *gin.Context) {
	var req dto.LogoutAllRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	accountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID format"})
		return
	}

	err = ctl.UseCase.LogoutAll(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logout all failed"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Successfully logged out from all devices"})
}

// RefreshToken handles token refresh using a valid refresh token
func (ctl *Controller) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	tokenPair, err := ctl.UseCase.RefreshToken(req.RefreshToken)
	if err != nil {
		switch err {
		case jwt.ErrTokenExpired:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token has expired"})
		case jwt.ErrTokenInvalid:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		case domain.ErrSessionNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		case domain.ErrSessionRevoked:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session has been revoked"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token refresh failed"})
		}
		return
	}

	response := dto.RefreshTokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
	}

	c.JSON(http.StatusOK, response)
}
