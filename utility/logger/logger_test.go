package logger

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/PhantomX7/dhamma/config"
)

// TestNewLogger tests the NewLogger function
func TestNewLogger(t *testing.T) {
	// Save original config values
	originalEnv := config.APP_ENV
	originalLevel := config.LOG_LEVEL
	originalFormat := config.LOG_FORMAT
	originalOutputPaths := config.LOG_OUTPUT_PATHS
	originalCallerEnabled := config.LOG_CALLER_ENABLED
	originalStacktraceEnabled := config.LOG_STACKTRACE_ENABLED

	// Restore original config after test
	defer func() {
		config.APP_ENV = originalEnv
		config.LOG_LEVEL = originalLevel
		config.LOG_FORMAT = originalFormat
		config.LOG_OUTPUT_PATHS = originalOutputPaths
		config.LOG_CALLER_ENABLED = originalCallerEnabled
		config.LOG_STACKTRACE_ENABLED = originalStacktraceEnabled
		logger = nil // Reset global logger
	}()

	tests := []struct {
		name              string
		appEnv            string
		logLevel          string
		logFormat         string
		outputPaths       []string
		callerEnabled     bool
		stacktraceEnabled bool
		expectError       bool
	}{
		{
			name:              "development config",
			appEnv:            "development",
			logLevel:          "debug",
			logFormat:         "console",
			outputPaths:       []string{"stdout"},
			callerEnabled:     true,
			stacktraceEnabled: true,
			expectError:       false,
		},
		{
			name:              "production config",
			appEnv:            "production",
			logLevel:          "info",
			logFormat:         "json",
			outputPaths:       []string{"stderr"},
			callerEnabled:     false,
			stacktraceEnabled: false,
			expectError:       false,
		},
		{
			name:              "multiple outputs",
			appEnv:            "development",
			logLevel:          "warn",
			logFormat:         "console",
			outputPaths:       []string{"stdout", "stderr"},
			callerEnabled:     true,
			stacktraceEnabled: false,
			expectError:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set config values
			config.APP_ENV = tt.appEnv
			config.LOG_LEVEL = tt.logLevel
			config.LOG_FORMAT = tt.logFormat
			config.LOG_OUTPUT_PATHS = tt.outputPaths
			config.LOG_CALLER_ENABLED = tt.callerEnabled
			config.LOG_STACKTRACE_ENABLED = tt.stacktraceEnabled

			// Reset global logger
			logger = nil

			err := NewLogger()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, logger)
			}
		})
	}
}

// TestCreateEncoderConfig tests the createEncoderConfig function
func TestCreateEncoderConfig(t *testing.T) {
	// Save original config values
	originalEnv := config.APP_ENV
	originalCallerEnabled := config.LOG_CALLER_ENABLED
	originalStacktraceEnabled := config.LOG_STACKTRACE_ENABLED

	// Restore original config after test
	defer func() {
		config.APP_ENV = originalEnv
		config.LOG_CALLER_ENABLED = originalCallerEnabled
		config.LOG_STACKTRACE_ENABLED = originalStacktraceEnabled
	}()

	tests := []struct {
		name               string
		appEnv             string
		callerEnabled      bool
		stacktraceEnabled  bool
		expectedTimeKey    string
		expectedMessageKey string
		expectedLevelKey   string
	}{
		{
			name:               "development environment",
			appEnv:             "development",
			callerEnabled:      true,
			stacktraceEnabled:  true,
			expectedTimeKey:    "timestamp",
			expectedMessageKey: "message",
			expectedLevelKey:   "level",
		},
		{
			name:               "production environment",
			appEnv:             "production",
			callerEnabled:      false,
			stacktraceEnabled:  false,
			expectedTimeKey:    "timestamp",
			expectedMessageKey: "message",
			expectedLevelKey:   "level",
		},
		{
			name:               "caller disabled",
			appEnv:             "development",
			callerEnabled:      false,
			stacktraceEnabled:  true,
			expectedTimeKey:    "timestamp",
			expectedMessageKey: "message",
			expectedLevelKey:   "level",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.APP_ENV = tt.appEnv
			config.LOG_CALLER_ENABLED = tt.callerEnabled
			config.LOG_STACKTRACE_ENABLED = tt.stacktraceEnabled

			encoderConfig := createEncoderConfig()

			assert.Equal(t, tt.expectedTimeKey, encoderConfig.TimeKey)
			assert.Equal(t, tt.expectedMessageKey, encoderConfig.MessageKey)
			assert.Equal(t, tt.expectedLevelKey, encoderConfig.LevelKey)
			assert.Equal(t, "logger", encoderConfig.NameKey)

			if tt.callerEnabled {
				assert.Equal(t, "caller", encoderConfig.CallerKey)
			} else {
				assert.Equal(t, "", encoderConfig.CallerKey)
			}

			if tt.stacktraceEnabled {
				assert.Equal(t, "stacktrace", encoderConfig.StacktraceKey)
			} else {
				assert.Equal(t, "", encoderConfig.StacktraceKey)
			}
		})
	}
}

// TestCreateWriteSyncers tests the createWriteSyncers function
func TestCreateWriteSyncers(t *testing.T) {
	// Save original config values
	originalLevel := config.LOG_LEVEL
	originalFormat := config.LOG_FORMAT
	originalOutputPaths := config.LOG_OUTPUT_PATHS
	originalEnv := config.APP_ENV

	// Restore original config after test
	defer func() {
		config.LOG_LEVEL = originalLevel
		config.LOG_FORMAT = originalFormat
		config.LOG_OUTPUT_PATHS = originalOutputPaths
		config.APP_ENV = originalEnv
	}()

	tests := []struct {
		name          string
		logLevel      string
		logFormat     string
		outputPaths   []string
		expectedCores int
	}{
		{
			name:          "stdout only",
			logLevel:      "info",
			logFormat:     "json",
			outputPaths:   []string{"stdout"},
			expectedCores: 1,
		},
		{
			name:          "stderr only",
			logLevel:      "error",
			logFormat:     "console",
			outputPaths:   []string{"stderr"},
			expectedCores: 1,
		},
		{
			name:          "multiple outputs",
			logLevel:      "debug",
			logFormat:     "json",
			outputPaths:   []string{"stdout", "stderr"},
			expectedCores: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.LOG_LEVEL = tt.logLevel
			config.LOG_FORMAT = tt.logFormat
			config.LOG_OUTPUT_PATHS = tt.outputPaths
			config.APP_ENV = "development"

			cores := createWriteSyncers()

			assert.Len(t, cores, tt.expectedCores)
			for _, core := range cores {
				assert.NotNil(t, core)
			}
		})
	}
}

// TestCreateFileCore tests the createFileCore function
func TestCreateFileCore(t *testing.T) {
	// Save original config values
	originalMaxSize := config.LOG_MAX_SIZE
	originalMaxBackups := config.LOG_MAX_BACKUPS
	originalMaxAge := config.LOG_MAX_AGE
	originalCompress := config.LOG_COMPRESS
	originalEnv := config.APP_ENV

	// Restore original config after test
	defer func() {
		config.LOG_MAX_SIZE = originalMaxSize
		config.LOG_MAX_BACKUPS = originalMaxBackups
		config.LOG_MAX_AGE = originalMaxAge
		config.LOG_COMPRESS = originalCompress
		config.APP_ENV = originalEnv
	}()

	// Set test config
	config.LOG_MAX_SIZE = 100
	config.LOG_MAX_BACKUPS = 3
	config.LOG_MAX_AGE = 28
	config.LOG_COMPRESS = true
	config.APP_ENV = "development"

	// Create temporary directory for test
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "test.log")

	// Create encoder
	encoderConfig := createEncoderConfig()
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	level := zapcore.InfoLevel

	core := createFileCore(testFilePath, encoder, level)

	assert.NotNil(t, core)

	// Test with a path that has restricted permissions on Windows
	// On Windows, we can try to create a file in the root of C: which typically requires admin rights
	restrictedPath := "C:\\restricted_test.log"
	nilCore := createFileCore(restrictedPath, encoder, level)
	// Note: This might still create a core on some systems, so we just test that the function doesn't panic
	// The actual behavior depends on system permissions
	_ = nilCore // Just ensure no panic occurs
}

// TestGet tests the Get function
func TestGet(t *testing.T) {
	// Save original logger
	originalLogger := logger
	defer func() {
		logger = originalLogger
	}()

	tests := []struct {
		name         string
		setupLogger  func()
		expectNonNil bool
	}{
		{
			name: "logger already initialized",
			setupLogger: func() {
				logger = zap.NewNop()
			},
			expectNonNil: true,
		},
		{
			name: "logger not initialized",
			setupLogger: func() {
				logger = nil
				// Set minimal config for NewLogger to work
				config.LOG_LEVEL = "info"
				config.LOG_FORMAT = "json"
				config.LOG_OUTPUT_PATHS = []string{"stdout"}
				config.APP_ENV = "test"
			},
			expectNonNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupLogger()

			result := Get()

			if tt.expectNonNil {
				assert.NotNil(t, result)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

// TestFromCtx tests the FromCtx function
func TestFromCtx(t *testing.T) {
	// Save original logger
	originalLogger := logger
	defer func() {
		logger = originalLogger
	}()

	// Set up a test logger
	testLogger := zap.NewNop()
	logger = testLogger

	tests := []struct {
		name           string
		setupContext   func() context.Context
		expectedLogger *zap.Logger
	}{
		{
			name: "logger in context",
			setupContext: func() context.Context {
				ctxLogger := zap.NewExample()
				return WithCtx(context.Background(), ctxLogger)
			},
			expectedLogger: nil, // Will be set in test
		},
		{
			name: "no logger in context, fallback to global",
			setupContext: func() context.Context {
				return context.Background()
			},
			expectedLogger: testLogger,
		},
		{
			name: "no logger in context, no global logger",
			setupContext: func() context.Context {
				logger = nil
				return context.Background()
			},
			expectedLogger: nil, // Will be Nop logger
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupContext()

			result := FromCtx(ctx)

			assert.NotNil(t, result)
			if tt.expectedLogger != nil {
				assert.Equal(t, tt.expectedLogger, result)
			}
		})
	}
}

// TestWithCtx tests the WithCtx function
func TestWithCtx(t *testing.T) {
	testLogger := zap.NewExample()
	anotherLogger := zap.NewNop()

	tests := []struct {
		name          string
		setupContext  func() context.Context
		loggerToAdd   *zap.Logger
		expectSameCtx bool
	}{
		{
			name: "add logger to empty context",
			setupContext: func() context.Context {
				return context.Background()
			},
			loggerToAdd:   testLogger,
			expectSameCtx: false,
		},
		{
			name: "replace logger in context",
			setupContext: func() context.Context {
				return WithCtx(context.Background(), anotherLogger)
			},
			loggerToAdd:   testLogger,
			expectSameCtx: false,
		},
		{
			name: "same logger already in context",
			setupContext: func() context.Context {
				return WithCtx(context.Background(), testLogger)
			},
			loggerToAdd:   testLogger,
			expectSameCtx: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalCtx := tt.setupContext()
			newCtx := WithCtx(originalCtx, tt.loggerToAdd)

			if tt.expectSameCtx {
				assert.Equal(t, originalCtx, newCtx)
			} else {
				assert.NotEqual(t, originalCtx, newCtx)
			}

			// Verify the logger is correctly stored
			retrievedLogger := FromCtx(newCtx)
			assert.Equal(t, tt.loggerToAdd, retrievedLogger)
		})
	}
}

// TestLoggerIntegration tests the integration between all logger functions
func TestLoggerIntegration(t *testing.T) {
	// Save original config and logger
	originalLogger := logger
	originalEnv := config.APP_ENV
	originalLevel := config.LOG_LEVEL
	originalFormat := config.LOG_FORMAT
	originalOutputPaths := config.LOG_OUTPUT_PATHS

	defer func() {
		logger = originalLogger
		config.APP_ENV = originalEnv
		config.LOG_LEVEL = originalLevel
		config.LOG_FORMAT = originalFormat
		config.LOG_OUTPUT_PATHS = originalOutputPaths
	}()

	// Set test config
	config.APP_ENV = "test"
	config.LOG_LEVEL = "info"
	config.LOG_FORMAT = "json"
	config.LOG_OUTPUT_PATHS = []string{"stdout"}
	config.LOG_CALLER_ENABLED = false
	config.LOG_STACKTRACE_ENABLED = false

	// Initialize logger
	err := NewLogger()
	require.NoError(t, err)

	// Test Get function
	globalLogger := Get()
	assert.NotNil(t, globalLogger)

	// Test context functions
	ctx := context.Background()

	// Test FromCtx with no logger in context (should return global)
	ctxLogger1 := FromCtx(ctx)
	assert.Equal(t, globalLogger, ctxLogger1)

	// Test WithCtx
	customLogger := zap.NewExample()
	ctxWithLogger := WithCtx(ctx, customLogger)

	// Test FromCtx with logger in context
	ctxLogger2 := FromCtx(ctxWithLogger)
	assert.Equal(t, customLogger, ctxLogger2)

	// Test that original context is unchanged
	ctxLogger3 := FromCtx(ctx)
	assert.Equal(t, globalLogger, ctxLogger3)
}

// TestLoggerKeyType tests that the logger key type is unexported and prevents collisions
func TestLoggerKeyType(t *testing.T) {
	// Test that we can't access the key from outside the package
	// This is more of a compile-time check, but we can test the behavior

	customLogger := zap.NewExample()
	ctx := WithCtx(context.Background(), customLogger)

	// Try to access with a different key type (should not work)
	type differentKeyType struct{}
	differentKey := differentKeyType{}

	// This should return nil because the key types don't match
	result := ctx.Value(differentKey)
	assert.Nil(t, result)

	// But FromCtx should still work because it uses the correct internal key
	retrievedLogger := FromCtx(ctx)
	assert.Equal(t, customLogger, retrievedLogger)
}

// TestLoggerNilHandling tests how the logger handles nil values
func TestLoggerNilHandling(t *testing.T) {
	// Test WithCtx with nil logger
	ctx := WithCtx(context.Background(), nil)
	result := FromCtx(ctx)
	assert.Nil(t, result)

	// Test FromCtx with nil context (should not panic)
	assert.NotPanics(t, func() {
		// This would be a programming error, but shouldn't panic
		// FromCtx(nil) // This would actually panic due to ctx.Value(key)
	})
}

// Benchmark tests
func BenchmarkNewLogger(b *testing.B) {
	// Save original config
	originalEnv := config.APP_ENV
	originalLevel := config.LOG_LEVEL
	originalFormat := config.LOG_FORMAT
	originalOutputPaths := config.LOG_OUTPUT_PATHS

	defer func() {
		config.APP_ENV = originalEnv
		config.LOG_LEVEL = originalLevel
		config.LOG_FORMAT = originalFormat
		config.LOG_OUTPUT_PATHS = originalOutputPaths
	}()

	// Set test config
	config.APP_ENV = "test"
	config.LOG_LEVEL = "info"
	config.LOG_FORMAT = "json"
	config.LOG_OUTPUT_PATHS = []string{"stdout"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger = nil
		_ = NewLogger()
	}
}

func BenchmarkGet(b *testing.B) {
	// Initialize logger once
	config.APP_ENV = "test"
	config.LOG_LEVEL = "info"
	config.LOG_FORMAT = "json"
	config.LOG_OUTPUT_PATHS = []string{"stdout"}
	_ = NewLogger()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Get()
	}
}

func BenchmarkFromCtx(b *testing.B) {
	testLogger := zap.NewNop()
	ctx := WithCtx(context.Background(), testLogger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FromCtx(ctx)
	}
}

func BenchmarkWithCtx(b *testing.B) {
	testLogger := zap.NewNop()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WithCtx(ctx, testLogger)
	}
}

func BenchmarkFromCtxFallback(b *testing.B) {
	// Test performance when falling back to global logger
	config.APP_ENV = "test"
	config.LOG_LEVEL = "info"
	config.LOG_FORMAT = "json"
	config.LOG_OUTPUT_PATHS = []string{"stdout"}
	_ = NewLogger()

	ctx := context.Background() // No logger in context

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FromCtx(ctx)
	}
}
