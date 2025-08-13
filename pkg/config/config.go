package config

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Environment types
const (
	EnvironmentDevelopment = "development"
	EnvironmentStaging     = "staging"
	EnvironmentProduction  = "production"
	EnvironmentTest        = "test"
)

// Config represents the application configuration
type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	Log         LogConfig
	Services    ServicesConfig
	Security    SecurityConfig
	Monitoring  MonitoringConfig
	Features    FeaturesConfig
	Email       EmailConfig
}

// ServerConfig represents server configuration
type ServerConfig struct {
	HTTPPort       string
	GRPCPort       string
	Host           string
	ClientURL      string // Client application URL for email verification links
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxHeaderBytes int
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	PostgreSQL PostgreSQLConfig
	Redis      RedisConfig
}

// PostgreSQLConfig represents PostgreSQL configuration
type PostgreSQLConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxConns        int
	MinConns        int
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
	Timeout         time.Duration
	ApplicationName string
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	Database int
	Timeout  time.Duration
	PoolSize int
}

// LogConfig represents logging configuration
type LogConfig struct {
	Level      string
	Format     string
	Output     string
	TimeFormat string
	Caller     bool
}

// ServicesConfig represents services configuration
type ServicesConfig struct {
	User UserServiceConfig
}

// UserServiceConfig represents user service configuration
type UserServiceConfig struct {
	Enabled bool
	Port    string
}

// SecurityConfig represents security configuration
type SecurityConfig struct {
	JWTSecret          string
	JWTExpiration      time.Duration
	SessionExpiration  time.Duration
	BCryptCost         int
	CORSAllowedOrigins []string
	CORSAllowedMethods []string
	CORSAllowedHeaders []string
	RateLimitRequests  int
	RateLimitWindow    time.Duration
	TrustedProxies     []string
}

// MonitoringConfig represents monitoring configuration
type MonitoringConfig struct {
	MetricsEnabled  bool
	MetricsPort     string
	HealthCheckPath string
	ReadinessPath   string
	LivenessPath    string
	PrometheusPath  string
}

// FeaturesConfig represents feature flags
type FeaturesConfig struct {
	UserRegistrationEnabled  bool
	EmailVerificationEnabled bool
	PasswordResetEnabled     bool
	MultiTenancyEnabled      bool
	AuditLoggingEnabled      bool
}

// EmailConfig represents email configuration
type EmailConfig struct {
	Enabled      bool
	Provider     string
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromAddress  string
	FromName     string
	TemplateDir  string
	UseTLS       bool
	UseSSL       bool
}

// Load loads configuration from environment variables with validation
func Load() (*Config, error) {
	// Get environment from ENVIRONMENT variable
	environment := getEnv("ENVIRONMENT", EnvironmentDevelopment)

	// Validate environment
	if !isValidEnvironment(environment) {
		return nil, fmt.Errorf("invalid environment: %s", environment)
	}

	config := &Config{
		Environment: environment,
		Server:      loadServerConfig(),
		Database:    loadDatabaseConfig(environment),
		Log:         loadLogConfig(),
		Services:    loadServicesConfig(),
		Security:    loadSecurityConfig(environment),
		Monitoring:  loadMonitoringConfig(),
		Features:    loadFeaturesConfig(),
		Email:       loadEmailConfig(environment),
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// loadServerConfig loads server configuration
func loadServerConfig() ServerConfig {
	return ServerConfig{
		HTTPPort:       normalizePort(getEnv("HTTP_PORT", ":8080")),
		GRPCPort:       normalizePort(getEnv("GRPC_PORT", ":50051")),
		Host:           getEnv("HOST", "localhost"),
		ClientURL:      getEnv("CLIENT_URL", "http://localhost:3000"), // Default to localhost for development
		ReadTimeout:    getEnvAsDuration("SERVER_READ_TIMEOUT", 30*time.Second),
		WriteTimeout:   getEnvAsDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
		IdleTimeout:    getEnvAsDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
		MaxHeaderBytes: getEnvAsInt("SERVER_MAX_HEADER_BYTES", 1<<20), // 1MB
	}
}

// loadDatabaseConfig loads database configuration based on environment
func loadDatabaseConfig(environment string) DatabaseConfig {
	return DatabaseConfig{
		PostgreSQL: loadPostgreSQLConfig(environment),
		Redis:      loadRedisConfig(environment),
	}
}

// loadPostgreSQLConfig loads PostgreSQL configuration based on environment
func loadPostgreSQLConfig(environment string) PostgreSQLConfig {
	// Check if DATABASE_URL is available (DigitalOcean)
	if databaseURL := getEnv("DATABASE_URL", ""); databaseURL != "" {
		// When DATABASE_URL is available, return a minimal config
		// The actual connection will use DATABASE_URL directly
		return PostgreSQLConfig{
			Host:            "", // Will be parsed from DATABASE_URL
			Port:            0,  // Will be parsed from DATABASE_URL
			User:            "", // Will be parsed from DATABASE_URL
			Password:        "", // Will be parsed from DATABASE_URL
			Database:        "", // Will be parsed from DATABASE_URL
			SSLMode:         getEnv("POSTGRES_SSLMODE", "require"),
			MaxConns:        getEnvAsInt("POSTGRES_MAX_CONNS", 20),
			MinConns:        getEnvAsInt("POSTGRES_MIN_CONNS", 5),
			MaxConnLifetime: getEnvAsDuration("POSTGRES_MAX_CONN_LIFETIME", 1*time.Hour),
			MaxConnIdleTime: getEnvAsDuration("POSTGRES_MAX_CONN_IDLE_TIME", 15*time.Minute),
			Timeout:         getEnvAsDuration("POSTGRES_TIMEOUT", 30*time.Second),
			ApplicationName: getEnv("POSTGRES_APP_NAME", "workluv"),
		}
	}

	switch environment {
	case EnvironmentProduction:
		return PostgreSQLConfig{
			Host:            getEnv("POSTGRES_HOST", ""),
			Port:            getEnvAsInt("POSTGRES_PORT", 5432),
			User:            getEnv("POSTGRES_USER", ""),
			Password:        getEnv("POSTGRES_PASSWORD", ""),
			Database:        getEnv("POSTGRES_DB", ""),
			SSLMode:         getEnv("POSTGRES_SSLMODE", "require"),
			MaxConns:        getEnvAsInt("POSTGRES_MAX_CONNS", 50),
			MinConns:        getEnvAsInt("POSTGRES_MIN_CONNS", 10),
			MaxConnLifetime: getEnvAsDuration("POSTGRES_MAX_CONN_LIFETIME", 1*time.Hour),
			MaxConnIdleTime: getEnvAsDuration("POSTGRES_MAX_CONN_IDLE_TIME", 15*time.Minute),
			Timeout:         getEnvAsDuration("POSTGRES_TIMEOUT", 30*time.Second),
			ApplicationName: getEnv("POSTGRES_APP_NAME", "workluv"),
		}
	case EnvironmentStaging:
		return PostgreSQLConfig{
			Host:            getEnv("POSTGRES_HOST", ""),
			Port:            getEnvAsInt("POSTGRES_PORT", 5432),
			User:            getEnv("POSTGRES_USER", ""),
			Password:        getEnv("POSTGRES_PASSWORD", ""),
			Database:        getEnv("POSTGRES_DB", ""),
			SSLMode:         getEnv("POSTGRES_SSLMODE", "require"),
			MaxConns:        getEnvAsInt("POSTGRES_MAX_CONNS", 30),
			MinConns:        getEnvAsInt("POSTGRES_MIN_CONNS", 5),
			MaxConnLifetime: getEnvAsDuration("POSTGRES_MAX_CONN_LIFETIME", 1*time.Hour),
			MaxConnIdleTime: getEnvAsDuration("POSTGRES_MAX_CONN_IDLE_TIME", 15*time.Minute),
			Timeout:         getEnvAsDuration("POSTGRES_TIMEOUT", 30*time.Second),
			ApplicationName: getEnv("POSTGRES_APP_NAME", "workluv"),
		}
	case EnvironmentDevelopment:
		// Check if we're in a cloud environment (DigitalOcean)
		if getEnv("POSTGRES_HOST", "") != "" {
			// Cloud development environment
			return PostgreSQLConfig{
				Host:            getEnv("POSTGRES_HOST", ""),
				Port:            getEnvAsInt("POSTGRES_PORT", 5432),
				User:            getEnv("POSTGRES_USER", ""),
				Password:        getEnv("POSTGRES_PASSWORD", ""),
				Database:        getEnv("POSTGRES_DB", ""),
				SSLMode:         getEnv("POSTGRES_SSLMODE", "require"),
				MaxConns:        getEnvAsInt("POSTGRES_MAX_CONNS", 20),
				MinConns:        getEnvAsInt("POSTGRES_MIN_CONNS", 5),
				MaxConnLifetime: getEnvAsDuration("POSTGRES_MAX_CONN_LIFETIME", 1*time.Hour),
				MaxConnIdleTime: getEnvAsDuration("POSTGRES_MAX_CONN_IDLE_TIME", 15*time.Minute),
				Timeout:         getEnvAsDuration("POSTGRES_TIMEOUT", 30*time.Second),
				ApplicationName: getEnv("POSTGRES_APP_NAME", "workluv"),
			}
		}
		// Local development environment
		return PostgreSQLConfig{
			Host:            getEnv("POSTGRES_HOST", "localhost"),
			Port:            getEnvAsInt("POSTGRES_PORT", 5432),
			User:            getEnv("POSTGRES_USER", "postgres"),
			Password:        getEnv("POSTGRES_PASSWORD", "password"),
			Database:        getEnv("POSTGRES_DB", "workluv"),
			SSLMode:         getEnv("POSTGRES_SSLMODE", "disable"),
			MaxConns:        getEnvAsInt("POSTGRES_MAX_CONNS", 20),
			MinConns:        getEnvAsInt("POSTGRES_MIN_CONNS", 5),
			MaxConnLifetime: getEnvAsDuration("POSTGRES_MAX_CONN_LIFETIME", 1*time.Hour),
			MaxConnIdleTime: getEnvAsDuration("POSTGRES_MAX_CONN_IDLE_TIME", 15*time.Minute),
			Timeout:         getEnvAsDuration("POSTGRES_TIMEOUT", 30*time.Second),
			ApplicationName: getEnv("POSTGRES_APP_NAME", "workluv"),
		}
	default: // fallback to local development
		return PostgreSQLConfig{
			Host:            getEnv("POSTGRES_HOST", "localhost"),
			Port:            getEnvAsInt("POSTGRES_PORT", 5432),
			User:            getEnv("POSTGRES_USER", "postgres"),
			Password:        getEnv("POSTGRES_PASSWORD", "password"),
			Database:        getEnv("POSTGRES_DB", "workluv"),
			SSLMode:         getEnv("POSTGRES_SSLMODE", "disable"),
			MaxConns:        getEnvAsInt("POSTGRES_MAX_CONNS", 20),
			MinConns:        getEnvAsInt("POSTGRES_MIN_CONNS", 5),
			MaxConnLifetime: getEnvAsDuration("POSTGRES_MAX_CONN_LIFETIME", 1*time.Hour),
			MaxConnIdleTime: getEnvAsDuration("POSTGRES_MAX_CONN_IDLE_TIME", 15*time.Minute),
			Timeout:         getEnvAsDuration("POSTGRES_TIMEOUT", 30*time.Second),
			ApplicationName: getEnv("POSTGRES_APP_NAME", "workluv"),
		}
	}
}

// loadRedisConfig loads Redis configuration based on environment
func loadRedisConfig(environment string) RedisConfig {
	switch environment {
	case EnvironmentProduction:
		return RedisConfig{
			Host:     getEnv("REDIS_HOST", ""),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			Database: getEnvAsInt("REDIS_DB", 0),
			Timeout:  getEnvAsDuration("REDIS_TIMEOUT", 5*time.Second),
			PoolSize: getEnvAsInt("REDIS_POOL_SIZE", 10),
		}
	case EnvironmentStaging:
		return RedisConfig{
			Host:     getEnv("REDIS_HOST", ""),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			Database: getEnvAsInt("REDIS_DB", 0),
			Timeout:  getEnvAsDuration("REDIS_TIMEOUT", 5*time.Second),
			PoolSize: getEnvAsInt("REDIS_POOL_SIZE", 10),
		}
	default: // development
		return RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			Database: getEnvAsInt("REDIS_DB", 0),
			Timeout:  getEnvAsDuration("REDIS_TIMEOUT", 5*time.Second),
			PoolSize: getEnvAsInt("REDIS_POOL_SIZE", 10),
		}
	}
}

// loadLogConfig loads logging configuration
func loadLogConfig() LogConfig {
	return LogConfig{
		Level:      getEnv("LOG_LEVEL", "info"),
		Format:     getEnv("LOG_FORMAT", "json"),
		Output:     getEnv("LOG_OUTPUT", "stdout"),
		TimeFormat: getEnv("LOG_TIME_FORMAT", "2006-01-02T15:04:05Z07:00"),
		Caller:     getEnvAsBool("LOG_CALLER", false),
	}
}

// loadServicesConfig loads services configuration
func loadServicesConfig() ServicesConfig {
	return ServicesConfig{
		User: UserServiceConfig{
			Enabled: getEnvAsBool("USER_SERVICE_ENABLED", true),
			Port:    normalizePort(getEnv("USER_SERVICE_PORT", ":8081")),
		},
	}
}

// loadSecurityConfig loads security configuration
func loadSecurityConfig(environment string) SecurityConfig {
	// Default CORS settings
	corsOrigins := []string{"*"}
	corsMethods := []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsHeaders := []string{"Content-Type", "Authorization", "X-Session-ID"}

	// Production CORS settings
	if environment == EnvironmentProduction {
		corsOrigins = strings.Split(getEnv("CORS_ALLOWED_ORIGINS", ""), ",")
		if len(corsOrigins) == 1 && corsOrigins[0] == "" {
			corsOrigins = []string{}
		}
	}

	return SecurityConfig{
		JWTSecret:          getEnv("JWT_SECRET", ""),
		JWTExpiration:      getEnvAsDuration("JWT_EXPIRATION", 24*time.Hour),
		SessionExpiration:  getEnvAsDuration("SESSION_EXPIRATION", 7*24*time.Hour),
		BCryptCost:         getEnvAsInt("BCRYPT_COST", 12),
		CORSAllowedOrigins: corsOrigins,
		CORSAllowedMethods: corsMethods,
		CORSAllowedHeaders: corsHeaders,
		RateLimitRequests:  getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:    getEnvAsDuration("RATE_LIMIT_WINDOW", 1*time.Minute),
		TrustedProxies:     strings.Split(getEnv("TRUSTED_PROXIES", ""), ","),
	}
}

// loadMonitoringConfig loads monitoring configuration
func loadMonitoringConfig() MonitoringConfig {
	return MonitoringConfig{
		MetricsEnabled:  getEnvAsBool("METRICS_ENABLED", true),
		MetricsPort:     normalizePort(getEnv("METRICS_PORT", ":9090")),
		HealthCheckPath: getEnv("HEALTH_CHECK_PATH", "/health"),
		ReadinessPath:   getEnv("READINESS_PATH", "/health/ready"),
		LivenessPath:    getEnv("LIVENESS_PATH", "/health/live"),
		PrometheusPath:  getEnv("PROMETHEUS_PATH", "/metrics"),
	}
}

// loadFeaturesConfig loads feature flags
func loadFeaturesConfig() FeaturesConfig {
	return FeaturesConfig{
		UserRegistrationEnabled:  getEnvAsBool("FEATURE_USER_REGISTRATION", true),
		EmailVerificationEnabled: getEnvAsBool("FEATURE_EMAIL_VERIFICATION", false),
		PasswordResetEnabled:     getEnvAsBool("FEATURE_PASSWORD_RESET", false),
		MultiTenancyEnabled:      getEnvAsBool("FEATURE_MULTI_TENANCY", false),
		AuditLoggingEnabled:      getEnvAsBool("FEATURE_AUDIT_LOGGING", true),
	}
}

// loadEmailConfig loads email configuration based on environment
func loadEmailConfig(environment string) EmailConfig {
	enabled := getEnvAsBool("EMAIL_SERVICE_ENABLED", false)
	if !enabled {
		return EmailConfig{Enabled: false}
	}

	provider := getEnv("EMAIL_SERVICE_PROVIDER", "smtp")
	smtpHost := getEnv("EMAIL_SMTP_HOST", "")
	smtpPort := getEnvAsInt("EMAIL_SMTP_PORT", 587)
	smtpUsername := getEnv("EMAIL_SMTP_USERNAME", "")
	smtpPassword := getEnv("EMAIL_SMTP_PASSWORD", "")
	fromAddress := getEnv("EMAIL_FROM_ADDRESS", "no-reply@example.com")
	fromName := getEnv("EMAIL_FROM_NAME", "Workluv")
	templateDir := getEnv("EMAIL_TEMPLATE_DIR", "./templates")
	useTLS := getEnvAsBool("EMAIL_SMTP_USE_TLS", true)
	useSSL := getEnvAsBool("EMAIL_SMTP_USE_SSL", false)

	return EmailConfig{
		Enabled:      enabled,
		Provider:     provider,
		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		SMTPUsername: smtpUsername,
		SMTPPassword: smtpPassword,
		FromAddress:  fromAddress,
		FromName:     fromName,
		TemplateDir:  templateDir,
		UseTLS:       useTLS,
		UseSSL:       useSSL,
	}
}

// normalizePort ensures port has proper format for container networking
func normalizePort(port string) string {
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	// For container environments, bind to all interfaces (0.0.0.0)
	// This ensures the application is accessible from outside the container
	if strings.HasPrefix(port, ":") {
		port = "0.0.0.0" + port
	}

	return port
}

// isValidEnvironment checks if the environment is valid
func isValidEnvironment(env string) bool {
	validEnvs := []string{
		EnvironmentDevelopment,
		EnvironmentStaging,
		EnvironmentProduction,
		EnvironmentTest,
	}

	for _, validEnv := range validEnvs {
		if env == validEnv {
			return true
		}
	}
	return false
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate server configuration
	if c.Server.HTTPPort == "" {
		return fmt.Errorf("HTTP_PORT is required")
	}

	// Validate database configuration (only if not using DATABASE_URL)
	if getEnv("DATABASE_URL", "") == "" {
		if c.Database.PostgreSQL.Host == "" {
			return fmt.Errorf("POSTGRES_HOST is required when DATABASE_URL is not set")
		}
		if c.Database.PostgreSQL.User == "" {
			return fmt.Errorf("POSTGRES_USER is required when DATABASE_URL is not set")
		}
		if c.Database.PostgreSQL.Database == "" {
			return fmt.Errorf("POSTGRES_DB is required when DATABASE_URL is not set")
		}
	}

	// Validate log configuration
	if c.Log.Level == "" {
		return fmt.Errorf("LOG_LEVEL is required")
	}

	// Validate security configuration
	if c.IsProduction() && c.Security.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required in production environment")
	}

	// Validate monitoring configuration
	if c.Monitoring.MetricsEnabled && c.Monitoring.MetricsPort == "" {
		return fmt.Errorf("METRICS_PORT is required when metrics are enabled")
	}

	// Validate email configuration
	if c.Email.Enabled && c.Email.SMTPHost == "" {
		return fmt.Errorf("EMAIL_SMTP_HOST is required when email is enabled")
	}
	if c.Email.Enabled && c.Email.SMTPPort == 0 {
		return fmt.Errorf("EMAIL_SMTP_PORT is required when email is enabled")
	}
	if c.Email.Enabled && c.Email.SMTPUsername == "" {
		return fmt.Errorf("EMAIL_SMTP_USERNAME is required when email is enabled")
	}
	if c.Email.Enabled && c.Email.SMTPPassword == "" {
		return fmt.Errorf("EMAIL_SMTP_PASSWORD is required when email is enabled")
	}
	if c.Email.Enabled && c.Email.FromAddress == "" {
		return fmt.Errorf("EMAIL_FROM_ADDRESS is required when email is enabled")
	}
	if c.Email.Enabled && c.Email.TemplateDir == "" {
		return fmt.Errorf("EMAIL_TEMPLATE_DIR is required when email is enabled")
	}

	return nil
}

// getEnv gets an environment variable with a default value
// It checks in order: system env -> .env file -> default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as an integer
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool gets an environment variable as a boolean
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvAsDuration gets an environment variable as a duration
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == EnvironmentDevelopment
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == EnvironmentProduction
}

// IsStaging returns true if the environment is staging
func (c *Config) IsStaging() bool {
	return c.Environment == EnvironmentStaging
}

// IsTest returns true if the environment is test
func (c *Config) IsTest() bool {
	return c.Environment == EnvironmentTest
}

// GetPostgreSQLDSN returns the PostgreSQL connection string
func (c *Config) GetPostgreSQLDSN() string {
	// Check if DATABASE_URL is available (DigitalOcean)
	if databaseURL := getEnv("DATABASE_URL", ""); databaseURL != "" {
		return databaseURL
	}

	// Fall back to individual variables
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s application_name=%s",
		c.Database.PostgreSQL.Host,
		c.Database.PostgreSQL.Port,
		c.Database.PostgreSQL.User,
		c.Database.PostgreSQL.Password,
		c.Database.PostgreSQL.Database,
		c.Database.PostgreSQL.SSLMode,
		c.Database.PostgreSQL.ApplicationName,
	)
}

// GetRedisAddr returns the Redis address
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Database.Redis.Host, c.Database.Redis.Port)
}

// GetTLSConfig returns TLS configuration for production
func (c *Config) GetTLSConfig() *tls.Config {
	if !c.IsProduction() {
		return nil
	}

	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
	}
}
