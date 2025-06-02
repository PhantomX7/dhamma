package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var (
	APP_ENV string
	PORT    string

	DATABASE_HOST     string
	DATABASE_PORT     string
	DATABASE_NAME     string
	DATABASE_USERNAME string
	DATABASE_PASSWORD string

	JWT_SECRET string

	// Logging Configuration
	LOG_LEVEL              string
	LOG_FORMAT             string
	LOG_OUTPUT_PATHS       []string
	LOG_ERROR_PATHS        []string
	LOG_FILE_ENABLED       bool
	LOG_FILE_PATH          string
	LOG_MAX_SIZE           int // MB
	LOG_MAX_BACKUPS        int
	LOG_MAX_AGE            int // days
	LOG_COMPRESS           bool
	LOG_CALLER_ENABLED     bool
	LOG_STACKTRACE_ENABLED bool
)

func LoadEnv() {
	// load environment
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	APP_ENV = os.Getenv("APP_ENV")
	PORT = os.Getenv("PORT")

	DATABASE_HOST = os.Getenv("DATABASE_HOST")
	DATABASE_PORT = os.Getenv("DATABASE_PORT")
	DATABASE_NAME = os.Getenv("DATABASE_NAME")
	DATABASE_USERNAME = os.Getenv("DATABASE_USERNAME")
	DATABASE_PASSWORD = os.Getenv("DATABASE_PASSWORD")

	JWT_SECRET = os.Getenv("JWT_SECRET")

	// Load logging configuration with defaults
	loadLoggingConfig()
}

// loadLoggingConfig loads logging-specific configuration with sensible defaults
func loadLoggingConfig() {
	// Log Level (debug, info, warn, error, dpanic, panic, fatal)
	LOG_LEVEL = getEnvWithDefault("LOG_LEVEL", "info")
	if APP_ENV == "development" {
		LOG_LEVEL = getEnvWithDefault("LOG_LEVEL", "debug")
	}

	// Log Format (console, json)
	LOG_FORMAT = getEnvWithDefault("LOG_FORMAT", "json")
	if APP_ENV == "development" {
		LOG_FORMAT = getEnvWithDefault("LOG_FORMAT", "console")
	}

	// Output paths
	outputPaths := getEnvWithDefault("LOG_OUTPUT_PATHS", "stdout")
	LOG_OUTPUT_PATHS = strings.Split(outputPaths, ",")

	// Error output paths
	errorPaths := getEnvWithDefault("LOG_ERROR_PATHS", "stderr")
	LOG_ERROR_PATHS = strings.Split(errorPaths, ",")

	// File logging configuration
	LOG_FILE_ENABLED = getEnvBoolWithDefault("LOG_FILE_ENABLED", false)
	LOG_FILE_PATH = getEnvWithDefault("LOG_FILE_PATH", "logs/app.log")
	LOG_MAX_SIZE = getEnvIntWithDefault("LOG_MAX_SIZE", 100) // 100MB
	LOG_MAX_BACKUPS = getEnvIntWithDefault("LOG_MAX_BACKUPS", 3)
	LOG_MAX_AGE = getEnvIntWithDefault("LOG_MAX_AGE", 28) // 28 days
	LOG_COMPRESS = getEnvBoolWithDefault("LOG_COMPRESS", true)

	// Additional logging features
	LOG_CALLER_ENABLED = getEnvBoolWithDefault("LOG_CALLER_ENABLED", true)
	LOG_STACKTRACE_ENABLED = getEnvBoolWithDefault("LOG_STACKTRACE_ENABLED", false)
	if APP_ENV == "development" {
		LOG_STACKTRACE_ENABLED = getEnvBoolWithDefault("LOG_STACKTRACE_ENABLED", true)
	}

	// Add file path to output paths if file logging is enabled
	if LOG_FILE_ENABLED {
		LOG_OUTPUT_PATHS = append(LOG_OUTPUT_PATHS, LOG_FILE_PATH)
	}
}

// Helper functions for environment variable parsing with defaults
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBoolWithDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvIntWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
