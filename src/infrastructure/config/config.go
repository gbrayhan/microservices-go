package config

import (
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTime    int64
	RefreshTime   int64
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnvOrDefault("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			User:     getEnvOrDefault("DB_USER", "postgres"),
			Password: getEnvOrDefault("DB_PASSWORD", "password"),
			DBName:   getEnvOrDefault("DB_NAME", "microservices_go"),
			SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			AccessSecret:  getEnvOrDefault("JWT_ACCESS_SECRET", "default_access_secret"),
			RefreshSecret: getEnvOrDefault("JWT_REFRESH_SECRET", "default_refresh_secret"),
			AccessTime:    getEnvAsInt64OrDefault("JWT_ACCESS_TIME_MINUTE", 60),
			RefreshTime:   getEnvAsInt64OrDefault("JWT_REFRESH_TIME_HOUR", 24),
		},
	}
}

// LoadTestConfig loads configuration for testing
func LoadTestConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8081",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "test_user",
			Password: "test_password",
			DBName:   "test_db",
			SSLMode:  "disable",
		},
		JWT: JWTConfig{
			AccessSecret:  "test_access_secret",
			RefreshSecret: "test_refresh_secret",
			AccessTime:    60,
			RefreshTime:   24,
		},
	}
}

// Helper functions
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt64OrDefault(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}
