package dto

import "workluv/internal/shared"

// CreateAccountRequest represents the request to create a new account
type CreateAccountRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// CreateAccountResponse represents the response after creating an account
type CreateAccountResponse struct {
	ID            string      `json:"id"`
	Email         string      `json:"email"`
	FullName      string      `json:"full_name"`
	IsActive      bool        `json:"is_active"`
	EmailVerified bool        `json:"email_verified"`
	CreatedAt     shared.Time `json:"created_at"`
	Message       string      `json:"message"`
}

// VerifyEmailRequest represents the request to verify an email
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

// ResendVerificationRequest represents the request to resend verification email
type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// EmailVerificationResponse represents the response for email verification operations
type EmailVerificationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
