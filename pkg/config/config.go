package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all configuration for our application
type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	Server   ServerConfig
	Log      LogConfig
	JWT      JWTConfig
	CORS     CORSConfig
	Email    EmailConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
	URL      string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port      string
	Host      string
	ClientURL string
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret        string
	RefreshSecret string
	AccessExpiry  time.Duration // JWT_EXPIRATION
	RefreshExpiry time.Duration // SESSION_EXPIRATION
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
}

// EmailConfig holds email service configuration
type EmailConfig struct {
	Enabled      bool
	Provider     string
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	UseTLS       bool
	UseSSL       bool
	FromAddress  string
	FromName     string
	TemplateDir  string
}

// Load loads configuration from environment variables
func Load() *Config {
	config := &Config{
		Database: loadDatabaseConfig(),
		Redis:    loadRedisConfig(),
		Server:   loadServerConfig(),
		Log:      loadLogConfig(),
		JWT:      loadJWTConfig(),
		CORS:     loadCORSConfig(),
		Email:    loadEmailConfig(),
	}

	// Build DATABASE_URL from individual components
	config.Database.URL = buildDatabaseURL(config.Database)

	return config
}

func loadDatabaseConfig() DatabaseConfig {
	port, _ := strconv.Atoi(getEnv("DATABASE_PORT", "5432"))

	return DatabaseConfig{
		Host:     getEnv("DATABASE_HOST", "localhost"),
		Port:     port,
		Name:     getEnv("DATABASE_NAME", "golang_arch"),
		User:     getEnv("DATABASE_USER", "postgres"),
		Password: getEnv("DATABASE_PASSWORD", "password"),
		SSLMode:  getEnv("DATABASE_SSL_MODE", "disable"),
	}
}

func loadRedisConfig() RedisConfig {
	port, _ := strconv.Atoi(getEnv("REDIS_PORT", "6379"))
	db, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	return RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     port,
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       db,
	}
}

func loadServerConfig() ServerConfig {
	return ServerConfig{
		Port:      getEnv("SERVER_PORT", "8080"),
		Host:      getEnv("SERVER_HOST", "http://localhost:8080"),
		ClientURL: getEnv("CLIENT_URL", "http://localhost:3000"),
	}
}

func loadLogConfig() LogConfig {
	return LogConfig{
		Level: getEnv("LOG_LEVEL", "info"),
	}
}

func loadJWTConfig() JWTConfig {
	// Parse duration strings from environment (e.g., "24h", "168h")
	accessExpiry, _ := time.ParseDuration(getEnv("JWT_EXPIRATION", "24h"))
	refreshExpiry, _ := time.ParseDuration(getEnv("SESSION_EXPIRATION", "168h")) // 7 days

	return JWTConfig{
		Secret:        getEnv("JWT_SECRET", "dev-jwt-secret"),
		RefreshSecret: getEnv("JWT_REFRESH_SECRET", "dev-refresh-secret"),
		AccessExpiry:  accessExpiry,
		RefreshExpiry: refreshExpiry,
	}
}

func buildDatabaseURL(db DatabaseConfig) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		db.User, db.Password, db.Host, db.Port, db.Name, db.SSLMode)
}

func loadCORSConfig() CORSConfig {
	allowedOriginsStr := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")
	allowedOrigins := strings.Split(allowedOriginsStr, ",")

	// Trim whitespace from each origin
	for i, origin := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(origin)
	}

	return CORSConfig{
		AllowedOrigins: allowedOrigins,
	}
}

func loadEmailConfig() EmailConfig {
	port, _ := strconv.Atoi(getEnv("EMAIL_SMTP_PORT", "587"))
	enabled, _ := strconv.ParseBool(getEnv("EMAIL_SERVICE_ENABLED", "false"))
	useTLS, _ := strconv.ParseBool(getEnv("EMAIL_SMTP_USE_TLS", "true"))
	useSSL, _ := strconv.ParseBool(getEnv("EMAIL_SMTP_USE_SSL", "false"))

	return EmailConfig{
		Enabled:      enabled,
		Provider:     getEnv("EMAIL_SERVICE_PROVIDER", "smtp"),
		SMTPHost:     getEnv("EMAIL_SMTP_HOST", "localhost"),
		SMTPPort:     port,
		SMTPUsername: getEnv("EMAIL_SMTP_USERNAME", ""),
		SMTPPassword: getEnv("EMAIL_SMTP_PASSWORD", ""),
		UseTLS:       useTLS,
		UseSSL:       useSSL,
		FromAddress:  getEnv("EMAIL_FROM_ADDRESS", "noreply@localhost.com"),
		FromName:     getEnv("EMAIL_FROM_NAME", "Go Server"),
		TemplateDir:  getEnv("EMAIL_TEMPLATE_DIR", "src/email"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
