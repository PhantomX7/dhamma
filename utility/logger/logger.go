package logger

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/PhantomX7/dhamma/config"
)

// loggerKeyType is an unexported type used for the context key.
// Using an unexported type prevents collisions with context keys defined in other packages.
type loggerKeyType struct{}

// loggerContextKey is the key used to store and retrieve the logger from the context.
var loggerContextKey = loggerKeyType{}

// Logger is a global instance of the zap logger.
var logger *zap.Logger

// NewLogger creates a new zap logger instance based on configuration.
// It supports configurable log levels, formats, file rotation, and environment-specific settings.
func NewLogger() error {
	// Create write syncers for different outputs
	writeSyncers := createWriteSyncers()

	// Create core with multiple outputs
	core := zapcore.NewTee(writeSyncers...)

	// Build logger with options
	options := []zap.Option{
		zap.AddCallerSkip(1),
	}

	// Add caller if enabled
	if config.LOG_CALLER_ENABLED {
		options = append(options, zap.AddCaller())
	}

	// Add stack trace if enabled
	if config.LOG_STACKTRACE_ENABLED {
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	// Create logger
	logger = zap.New(core, options...)

	return nil
}

// createEncoderConfig creates encoder configuration based on environment and settings
func createEncoderConfig() zapcore.EncoderConfig {
	var encoderConfig zapcore.EncoderConfig

	if config.APP_ENV == "development" {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	}

	// Common configuration
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.MessageKey = "message"
	encoderConfig.LevelKey = "level"
	encoderConfig.NameKey = "logger"

	// Configure caller and stacktrace keys
	if config.LOG_CALLER_ENABLED {
		encoderConfig.CallerKey = "caller"
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	} else {
		encoderConfig.CallerKey = ""
	}

	if config.LOG_STACKTRACE_ENABLED {
		encoderConfig.StacktraceKey = "stacktrace"
	} else {
		encoderConfig.StacktraceKey = ""
	}

	return encoderConfig
}

// createWriteSyncers creates write syncers for different output destinations
func createWriteSyncers() []zapcore.Core {
	var cores []zapcore.Core

	// Parse log level
	level, _ := zapcore.ParseLevel(config.LOG_LEVEL)

	// Create encoder
	encoderConfig := createEncoderConfig()
	var encoder zapcore.Encoder
	if strings.ToLower(config.LOG_FORMAT) == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// Add stdout/stderr cores
	for _, path := range config.LOG_OUTPUT_PATHS {
		if path == "stdout" {
			cores = append(cores, zapcore.NewCore(
				encoder,
				zapcore.AddSync(os.Stdout),
				level,
			))
		} else if path == "stderr" {
			cores = append(cores, zapcore.NewCore(
				encoder,
				zapcore.AddSync(os.Stderr),
				level,
			))
		} else {
			// File output with rotation
			fileCore := createFileCore(path, encoder, level)
			if fileCore != nil {
				cores = append(cores, fileCore)
			}
		}
	}

	return cores
}

// createFileCore creates a file-based core with log rotation
func createFileCore(filePath string, encoder zapcore.Encoder, level zapcore.Level) zapcore.Core {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		// If we can't create the directory, skip file logging
		return nil
	}

	// Create lumberjack logger for rotation
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    config.LOG_MAX_SIZE,
		MaxBackups: config.LOG_MAX_BACKUPS,
		MaxAge:     config.LOG_MAX_AGE,
		Compress:   config.LOG_COMPRESS,
	}

	return zapcore.NewCore(
		encoder,
		zapcore.AddSync(lumberjackLogger),
		level,
	)
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
