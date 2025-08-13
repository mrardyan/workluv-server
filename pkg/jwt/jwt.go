package jwt

import (
	"errors"
	"time"
	"workluv/internal/shared"
	"workluv/pkg/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrTokenExpired   = errors.New("token has expired")
	ErrTokenInvalid   = errors.New("token is invalid")
	ErrTokenMalformed = errors.New("token is malformed")
	ErrTokenSignature = errors.New("invalid token signature")
)

type Claims struct {
	AccountID uuid.UUID `json:"account_id"`
	Email     string    `json:"email"`
	TokenType string    `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type Service struct {
	config *config.Config
}

func NewService(cfg *config.Config) *Service {
	return &Service{config: cfg}
}

// GenerateAccessToken creates a new JWT access token
func (s *Service) GenerateAccessToken(accountID uuid.UUID, email string) (string, shared.Time, error) {
	expiryTime := time.Now().Add(s.config.JWT.AccessExpiry)
	claims := Claims{
		AccountID: accountID,
		Email:     email,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   accountID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", shared.Time{}, err
	}

	return tokenString, shared.NewTime(expiryTime), nil
}

// GenerateRefreshToken creates a new JWT refresh token
func (s *Service) GenerateRefreshToken(accountID uuid.UUID, email string) (string, shared.Time, error) {
	expiryTime := time.Now().Add(s.config.JWT.RefreshExpiry)
	claims := Claims{
		AccountID: accountID,
		Email:     email,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   accountID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWT.RefreshSecret))
	if err != nil {
		return "", shared.Time{}, err
	}

	return tokenString, shared.NewTime(expiryTime), nil
}

// GenerateTokenPair creates both access and refresh tokens
func (s *Service) GenerateTokenPair(accountID uuid.UUID, email string) (*TokenPair, shared.Time, error) {
	accessToken, accessExpiry, err := s.GenerateAccessToken(accountID, email)
	if err != nil {
		return nil, shared.Time{}, err
	}

	refreshToken, refreshExpiry, err := s.GenerateRefreshToken(accountID, email)
	if err != nil {
		return nil, shared.Time{}, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiry.ToEpoch(),
	}, refreshExpiry, nil
}

// ValidateAccessToken validates and parses an access token
func (s *Service) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenSignature
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		}
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid || claims.TokenType != "access" {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

// ValidateRefreshToken validates and parses a refresh token
func (s *Service) ValidateRefreshToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenSignature
		}
		return []byte(s.config.JWT.RefreshSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		}
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid || claims.TokenType != "refresh" {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}
