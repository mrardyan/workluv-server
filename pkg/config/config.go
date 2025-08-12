package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for our application
type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	Server   ServerConfig
	Log      LogConfig
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
	Port string
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level string
}

// Load loads configuration from environment variables
func Load() *Config {
	config := &Config{
		Database: loadDatabaseConfig(),
		Redis:    loadRedisConfig(),
		Server:   loadServerConfig(),
		Log:      loadLogConfig(),
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
		Port: getEnv("SERVER_PORT", "8080"),
	}
}

func loadLogConfig() LogConfig {
	return LogConfig{
		Level: getEnv("LOG_LEVEL", "info"),
	}
}

func buildDatabaseURL(db DatabaseConfig) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		db.User, db.Password, db.Host, db.Port, db.Name, db.SSLMode)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
