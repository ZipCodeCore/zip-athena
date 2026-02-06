package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all application configuration
type Config struct {
	// Server configuration
	Server ServerConfig

	// Database configuration
	Database DatabaseConfig

	// Security configuration
	Security SecurityConfig

	// CORS configuration
	CORS CORSConfig

	// Logging configuration
	Logging LoggingConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Type string // "sqlite" or "postgres"
	Path string // For SQLite
	URL  string // For PostgreSQL

	// Connection pool settings
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	JWTSecret       string
	JWTExpiration   time.Duration
	BCryptCost      int
	SessionDuration time.Duration
}

// CORSConfig holds CORS-related configuration
type CORSConfig struct {
	Enabled            bool
	AllowedOrigins     []string
	AllowedMethods     []string
	AllowedHeaders     []string
	AllowedCredentials bool
}

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level  string // "debug", "info", "warn", "error"
	Format string // "json" or "text"
	Output string // "stdout", "stderr", or file path
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		Server:   loadServerConfig(),
		Database: loadDatabaseConfig(),
		Security: loadSecurityConfig(),
		CORS:     loadCORSConfig(),
		Logging:  loadLoggingConfig(),
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// loadServerConfig loads server configuration from environment
func loadServerConfig() ServerConfig {
	port := getEnv("PORT", "9206")
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	return ServerConfig{
		Port:         port,
		Host:         getEnv("HOST", "0.0.0.0"),
		ReadTimeout:  getDurationEnv("READ_TIMEOUT", 15*time.Second),
		WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:  getDurationEnv("IDLE_TIMEOUT", 60*time.Second),
	}
}

// loadDatabaseConfig loads database configuration from environment
func loadDatabaseConfig() DatabaseConfig {
	dbType := getEnv("DB_TYPE", "sqlite")
	dbPath := getEnv("DB_PATH", "./athena.db")
	dbURL := getEnv("DATABASE_URL", "")

	return DatabaseConfig{
		Type:            dbType,
		Path:            dbPath,
		URL:             dbURL,
		MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
	}
}

// loadSecurityConfig loads security configuration from environment
func loadSecurityConfig() SecurityConfig {
	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		// Generate a warning-worthy default for development
		jwtSecret = "INSECURE_DEFAULT_SECRET_CHANGE_IN_PRODUCTION"
	}

	return SecurityConfig{
		JWTSecret:       jwtSecret,
		JWTExpiration:   getDurationEnv("JWT_EXPIRATION", 24*time.Hour),
		BCryptCost:      getIntEnv("BCRYPT_COST", 10),
		SessionDuration: getDurationEnv("SESSION_DURATION", 24*time.Hour),
	}
}

// loadCORSConfig loads CORS configuration from environment
func loadCORSConfig() CORSConfig {
	enabled := getBoolEnv("CORS_ENABLED", true)

	// Parse allowed origins
	originsStr := getEnv("CORS_ALLOWED_ORIGINS", "*")
	origins := strings.Split(originsStr, ",")
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}

	// Parse allowed methods
	methodsStr := getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS")
	methods := strings.Split(methodsStr, ",")
	for i, method := range methods {
		methods[i] = strings.TrimSpace(method)
	}

	// Parse allowed headers
	headersStr := getEnv("CORS_ALLOWED_HEADERS", "Content-Type,Authorization,X-User-ID")
	headers := strings.Split(headersStr, ",")
	for i, header := range headers {
		headers[i] = strings.TrimSpace(header)
	}

	return CORSConfig{
		Enabled:            enabled,
		AllowedOrigins:     origins,
		AllowedMethods:     methods,
		AllowedHeaders:     headers,
		AllowedCredentials: getBoolEnv("CORS_ALLOW_CREDENTIALS", true),
	}
}

// loadLoggingConfig loads logging configuration from environment
func loadLoggingConfig() LoggingConfig {
	return LoggingConfig{
		Level:  getEnv("LOG_LEVEL", "info"),
		Format: getEnv("LOG_FORMAT", "text"),
		Output: getEnv("LOG_OUTPUT", "stdout"),
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate server config
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	// Validate database config
	if c.Database.Type != "sqlite" && c.Database.Type != "postgres" {
		return fmt.Errorf("invalid database type: %s (must be 'sqlite' or 'postgres')", c.Database.Type)
	}

	if c.Database.Type == "sqlite" && c.Database.Path == "" {
		return fmt.Errorf("database path is required for SQLite")
	}

	if c.Database.Type == "postgres" && c.Database.URL == "" {
		return fmt.Errorf("database URL is required for PostgreSQL")
	}

	// Validate security config
	if c.Security.JWTSecret == "INSECURE_DEFAULT_SECRET_CHANGE_IN_PRODUCTION" {
		fmt.Fprintf(os.Stderr, "WARNING: Using insecure default JWT secret. Set JWT_SECRET environment variable.\n")
	}

	if c.Security.BCryptCost < 4 || c.Security.BCryptCost > 31 {
		return fmt.Errorf("bcrypt cost must be between 4 and 31")
	}

	// Validate logging config
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLogLevels[c.Logging.Level] {
		return fmt.Errorf("invalid log level: %s (must be debug, info, warn, or error)", c.Logging.Level)
	}

	if c.Logging.Format != "json" && c.Logging.Format != "text" {
		return fmt.Errorf("invalid log format: %s (must be json or text)", c.Logging.Format)
	}

	return nil
}

// Helper functions for reading environment variables

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// GetAddress returns the full server address (host:port)
func (s *ServerConfig) GetAddress() string {
	if s.Host == "0.0.0.0" || s.Host == "" {
		return s.Port
	}
	return s.Host + s.Port
}
