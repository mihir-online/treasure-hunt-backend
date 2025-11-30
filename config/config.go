package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	App      AppConfig
	QRCode   QRCodeConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host string
	Port string
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// AppConfig holds general application configuration
type AppConfig struct {
	Environment    string
	LogLevel       string
	AllowedOrigins string
}

// QRCodeConfig holds QR code generation configuration
type QRCodeConfig struct {
	Size       int
	BaseURL    string
	QRPath     string
	StorageDir string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host: getEnv("HOST"),
			Port: getEnv("PORT"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST"),
			Port:     getEnv("DB_PORT"),
			User:     getEnv("DB_USER"),
			Password: getEnv("DB_PASSWORD"),
			DBName:   getEnv("DB_NAME"),
			SSLMode:  getEnv("DB_SSLMODE"),
		},
		App: AppConfig{
			Environment:    getEnv("APP_ENV"),
			LogLevel:       getEnv("LOG_LEVEL"),
			AllowedOrigins: getEnv("ALLOWED_ORIGINS"),
		},
		QRCode: QRCodeConfig{
			Size:       getEnvAsInt("QR_CODE_SIZE"),
			BaseURL:    getEnv("QR_CODE_BASE_URL"),
			QRPath:     getEnv("QR_PATH"),
			StorageDir: getEnv("QR_STORAGE_DIR"),
		},
	}
}

// getEnv gets an environment variable (required)
func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Environment variable %s is required but not set", key))
	}
	return value
}

// getEnvAsInt gets an environment variable as an integer (required)
func getEnvAsInt(key string) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		panic(fmt.Sprintf("Environment variable %s is required but not set", key))
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		panic(fmt.Sprintf("Environment variable %s must be an integer, got: %s", key, valueStr))
	}
	return value
}
