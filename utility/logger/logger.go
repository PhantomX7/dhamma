package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// loggerKeyType is an unexported type used for the context key.
// Using an unexported type prevents collisions with context keys defined in other packages.
type loggerKeyType struct{}

// loggerContextKey is the key used to store and retrieve the logger from the context.
var loggerContextKey = loggerKeyType{}

// Logger is a global instance of the zap logger.
var logger *zap.Logger

// NewLogger creates a new zap logger instance based on the APP_ENV environment variable.
// It uses Console encoding for "development" and JSON encoding otherwise (production).
func NewLogger() error {
	var config zap.Config
	var encoderConfig zapcore.EncoderConfig

	// Check the environment variable to decide on base config and encoding
	env := os.Getenv("APP_ENV")
	if env == "development" {
		// Use development config (defaults to console encoding)
		config = zap.NewDevelopmentConfig()
		encoderConfig = zap.NewDevelopmentEncoderConfig()

		// Customize development encoder
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		// encoderConfig.StacktraceKey = "stacktrace"
		encoderConfig.StacktraceKey = ""
		encoderConfig.CallerKey = "caller"
		encoderConfig.MessageKey = "msg"
		encoderConfig.NameKey = "name"

	} else {
		// Use production config (defaults to JSON encoding)
		config = zap.NewProductionConfig()
		encoderConfig = zap.NewProductionEncoderConfig()

		// Customize production encoder (JSON)
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.StacktraceKey = ""
		// LevelKey, CallerKey, MessageKey, NameKey are standard in production JSON encoder
	}

	// Apply the configured encoder settings
	config.EncoderConfig = encoderConfig

	// Set output paths
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	var err error
	// Build the logger from the final configuration
	// AddCallerSkip(1) helps ensure the caller field shows the correct location in your code
	logger, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}
	return nil
}

// Get returns the global logger instance.
// Note: Prefer using FromCtx when context is available.
func Get() *zap.Logger {
	// Initialize logger if it hasn't been already (e.g., if Get is called before NewLogger)
	// This is a defensive measure; ideally NewLogger is called first during app setup.
	if logger == nil {
		err := NewLogger()
		if err != nil {
			// Fallback to Nop logger if initialization fails here
			return zap.NewNop()
		}
	}
	return logger
}

// FromCtx retrieves the logger instance from the context.
// It falls back to the global logger if not found in context, and finally to a Nop logger.
func FromCtx(ctx context.Context) *zap.Logger {
	// Attempt to retrieve logger from context using the unexported key type
	if l, ok := ctx.Value(loggerContextKey).(*zap.Logger); ok {
		return l
	}
	// Fallback to global logger if not in context
	if l := Get(); l != nil { // Use Get() to ensure logger is initialized
		return l
	}
	// Final fallback to Nop logger
	return zap.NewNop()
}

// WithCtx embeds the provided logger instance into the context using the unexported key type.
// It returns the new context with the logger embedded.
func WithCtx(ctx context.Context, l *zap.Logger) context.Context {
	// Avoid creating a new context if the logger is already the one present
	if lp, ok := ctx.Value(loggerContextKey).(*zap.Logger); ok {
		if lp == l {
			return ctx
		}
	}
	// Embed the logger into the context
	return context.WithValue(ctx, loggerContextKey, l)
}
