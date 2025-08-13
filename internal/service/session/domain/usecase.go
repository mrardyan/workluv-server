package domain

import (
	"context"
	"errors"
	accountDomain "workluv/internal/service/account/domain"
	"workluv/internal/shared"
	"workluv/pkg/jwt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrSessionNotFound    = errors.New("session not found")
	ErrSessionExpired     = errors.New("session has expired")
	ErrSessionRevoked     = errors.New("session has been revoked")
)

// UseCase handles session business logic with minimal methods
type UseCase struct {
	sessionRepo Repository
	accountRepo accountDomain.Repository
	jwtService  *jwt.Service
}

// NewUseCase creates a new session use case
func NewUseCase(sessionRepo Repository, accountRepo accountDomain.Repository, jwtService *jwt.Service) *UseCase {
	return &UseCase{
		sessionRepo: sessionRepo,
		accountRepo: accountRepo,
		jwtService:  jwtService,
	}
}

// LoginRequest represents the login request data
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response data
type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresAt    int64       `json:"expires_at"`
	Account      interface{} `json:"account"`
}

// Login authenticates user and creates a new session
func (uc *UseCase) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Find account by email (assuming FindByEmail exists in account repository)
	account, err := uc.findAccountByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate token pair
	tokenPair, refreshExpiry, err := uc.jwtService.GenerateTokenPair(account.ID, string(account.Email))
	if err != nil {
		return nil, err
	}

	// Create session in database
	session := Session{
		AccountID:    account.ID,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    refreshExpiry,
		RevokedAt:    nil,
		CreatedAt:    shared.Now(),
		UpdatedAt:    shared.Now(),
	}

	_, err = uc.sessionRepo.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		Account:      account,
	}, nil
}

// Logout revokes a specific session
func (uc *UseCase) Logout(ctx context.Context, refreshToken string) error {
	session, err := uc.sessionRepo.FindByRefreshToken(ctx, refreshToken)
	if err != nil {
		return ErrSessionNotFound
	}

	if !session.IsActive() {
		return ErrSessionRevoked
	}

	return uc.sessionRepo.RevokeByID(ctx, session.ID)
}

// LogoutAll revokes all sessions for a user
func (uc *UseCase) LogoutAll(ctx context.Context, accountID uuid.UUID) error {
	return uc.sessionRepo.RevokeAllByAccountID(ctx, accountID)
}

// RefreshToken generates a new access token using a valid refresh token
func (uc *UseCase) RefreshToken(ctx context.Context, refreshToken string) (*jwt.TokenPair, error) {
	// Validate refresh token
	claims, err := uc.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Find session in database
	session, err := uc.sessionRepo.FindByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	// Check if session is active
	if !session.IsActive() {
		return nil, ErrSessionRevoked
	}

	// Generate new token pair
	tokenPair, _, err := uc.jwtService.GenerateTokenPair(claims.AccountID, claims.Email)
	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

// findAccountByEmail finds an account by email address
func (uc *UseCase) findAccountByEmail(ctx context.Context, email string) (accountDomain.Account, error) {
	return uc.accountRepo.FindByEmail(ctx, accountDomain.Email(email))
}
