package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a global instance of the zap logger.
var Logger *zap.Logger

// InitializeLogger initializes the global logger instance.

// NewLogger creates a new zap logger instance based on the APP_ENV environment variable.
// It defaults to a production configuration and switches to development if APP_ENV is "development".
func NewLogger() error {
	// Default to production configuration
	config := zap.NewProductionConfig()

	// Check the environment variable
	env := os.Getenv("APP_ENV")
	if env == "development" {
		// Use development configuration if specified
		config = zap.NewDevelopmentConfig()
		// Optional: Adjust development logging level if needed, e.g., DebugLevel
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	// Apply common encoder settings regardless of environment
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// Keep stacktrace for development, disable for production (default in NewProductionConfig)
	if env != "development" {
		config.EncoderConfig.StacktraceKey = "" // Disable stacktrace in production JSON
	} else {
		// Development config usually includes stacktrace by default, but ensure it's set if needed
		config.EncoderConfig.StacktraceKey = "stacktrace"
	}

	// Set output paths (stdout is often default, but explicit is fine)
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	var err error
	// Build the logger from the final configuration
	Logger, err = config.Build()
	if err != nil {
		return err
	}
	return nil
}
