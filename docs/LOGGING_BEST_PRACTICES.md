# Logging and Error Handling Best Practices

This document outlines the logging and error handling patterns implemented in the Dhamma project, providing guidelines and examples for consistent usage across the application.

## Table of Contents

1. [Overview](#overview)
2. [Configuration](#configuration)
3. [Error Types and Handling](#error-types-and-handling)
4. [Logging Patterns](#logging-patterns)
5. [Middleware Usage](#middleware-usage)
6. [Circuit Breaker and Recovery](#circuit-breaker-and-recovery)
7. [Health Monitoring](#health-monitoring)
8. [Best Practices](#best-practices)
9. [Examples](#examples)

## Overview

The Dhamma project implements a comprehensive logging and error handling system with the following key components:

- **Structured Logging**: Using Zap logger with configurable levels and formats
- **Error Types**: Standardized error types with proper HTTP status codes
- **Request Tracing**: Request ID and trace ID for distributed tracing
- **Circuit Breaker**: Resilience patterns for handling failures
- **Recovery Middleware**: Panic recovery with detailed logging
- **Health Monitoring**: System health and circuit breaker statistics

## Configuration

### Environment Variables

Configure logging behavior through environment variables:

```bash
# Log Level (debug, info, warn, error, fatal, panic)
LOG_LEVEL=info

# Log Format (json, console)
LOG_FORMAT=json

# Output destinations (comma-separated: stdout, stderr, file)
LOG_OUTPUT=stdout,file

# File logging settings
LOG_FILE_PATH=logs/app.log
LOG_FILE_MAX_SIZE=100
LOG_FILE_MAX_BACKUPS=5
LOG_FILE_MAX_AGE=30
LOG_FILE_COMPRESS=true

# Additional logging features
LOG_ENABLE_CALLER=true
LOG_ENABLE_STACKTRACE=true
LOG_ENABLE_SAMPLING=false
```

### Code Configuration

```go
// Initialize logger with configuration
logger.Init()

// Get logger instance
log := logger.Get()

// Get context-aware logger
ctxLogger := logger.FromCtx(ctx)
```

## Error Types and Handling

### Available Error Types

| Error Type | HTTP Status | Use Case |
|------------|-------------|----------|
| `ValidationError` | 400 | Input validation failures |
| `AuthenticationError` | 401 | Authentication failures |
| `AuthorizationError` | 403 | Permission denied |
| `NotFoundError` | 404 | Resource not found |
| `ConflictError` | 409 | Resource conflicts |
| `TooManyRequestsError` | 429 | Rate limiting |
| `InternalError` | 500 | Internal server errors |
| `ServiceUnavailableError` | 503 | Service unavailable |
| `DatabaseError` | 500 | Database operation failures |
| `ExternalServiceError` | 502 | External service failures |
| `BusinessError` | 422 | Business logic violations |
| `RateLimitError` | 429 | Rate limit exceeded |

### Creating Errors

```go
import "github.com/PhantomX7/dhamma/utility/errors"

// Validation error
err := errors.NewValidationError("Invalid email format", nil)

// Authentication error
err := errors.NewAuthenticationError("Invalid credentials", nil)

// With additional details
err := errors.NewValidationError("Validation failed", nil).
    WithDetails(map[string]interface{}{
        "field": "email",
        "value": "invalid-email",
    }).
    WithRequestID("req-123")

// Database error with underlying error
err := errors.NewDatabaseError("Failed to create user", dbErr)
```

### Error Response Format

```json
{
  "success": false,
  "error": {
    "code": "E1001",
    "message": "Validation failed",
    "type": "VALIDATION_ERROR",
    "details": {
      "field": "email",
      "value": "invalid-email"
    }
  },
  "request_id": "req-123",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Logging Patterns

### Structured Logging

```go
import (
    "go.uber.org/zap"
    "github.com/PhantomX7/dhamma/utility/logger"
)

// Basic logging
log := logger.Get()
log.Info("User created successfully",
    zap.String("user_id", userID),
    zap.String("email", email),
    zap.Duration("duration", time.Since(start)),
)

// Error logging
log.Error("Failed to create user",
    zap.Error(err),
    zap.String("user_id", userID),
    zap.String("operation", "create_user"),
)

// Context-aware logging
ctxLogger := logger.FromCtx(c.Request.Context())
ctxLogger.Info("Processing request",
    zap.String("endpoint", c.Request.URL.Path),
    zap.String("method", c.Request.Method),
)
```

### Log Levels

- **Debug**: Detailed information for debugging
- **Info**: General information about application flow
- **Warn**: Warning conditions that should be noted
- **Error**: Error conditions that need attention
- **Fatal**: Critical errors that cause application termination
- **Panic**: Critical errors that cause panic

### Request Context Logging

```go
// In handlers, use context logger for request-scoped logging
func (h *UserHandler) CreateUser(c *gin.Context) {
    // Get context logger with request information
    log := logger.FromCtx(c.Request.Context())
    
    log.Info("Creating new user",
        zap.String("email", request.Email),
        zap.String("role", request.Role),
    )
    
    // ... business logic ...
    
    if err != nil {
        log.Error("Failed to create user",
            zap.Error(err),
            zap.String("email", request.Email),
        )
        c.Error(errors.NewDatabaseError("Failed to create user", err))
        return
    }
    
    log.Info("User created successfully",
        zap.String("user_id", user.ID),
        zap.String("email", user.Email),
    )
}
```

## Middleware Usage

### Logger Middleware

```go
// Basic usage with default configuration
router.Use(middleware.Logger())

// Custom configuration
config := middleware.LoggerConfig{
    SkipPaths: []string{"/health", "/metrics"},
    EnableRequestBody: true,
    EnableResponseBody: true,
    MaxBodySize: 1024,
    SensitiveHeaders: []string{"authorization", "x-api-key"},
    SlowRequestThreshold: 5 * time.Second,
}
router.Use(middleware.LoggerWithConfig(config))
```

### Error Handler Middleware

```go
// Add error handler middleware (should be last)
router.Use(middleware.ErrorHandler())

// In handlers, use c.Error() to set errors
func (h *Handler) SomeHandler(c *gin.Context) {
    if err := someOperation(); err != nil {
        c.Error(errors.NewDatabaseError("Operation failed", err))
        return
    }
    
    c.JSON(http.StatusOK, response)
}
```

### Recovery Middleware

```go
// Basic usage with default configuration
router.Use(middleware.Recovery())

// Custom configuration
config := middleware.RecoveryConfig{
    EnableStackTrace: true,
    EnableCircuitBreaker: true,
    MaxRequests: 10,
    Interval: 60 * time.Second,
    Timeout: 30 * time.Second,
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
        return counts.Requests >= 5 && failureRatio >= 0.6
    },
}
router.Use(middleware.RecoveryWithConfig(config))
```

## Circuit Breaker and Recovery

### Circuit Breaker States

- **Closed**: Normal operation, requests pass through
- **Open**: Circuit is open, requests are rejected immediately
- **Half-Open**: Limited requests allowed to test if service recovered

### Configuration

```go
config := middleware.RecoveryConfig{
    EnableCircuitBreaker: true,
    MaxRequests: 10,                    // Max requests in half-open state
    Interval: 60 * time.Second,         // Interval to clear counts
    Timeout: 30 * time.Second,          // Timeout before half-open
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        // Custom logic to determine when to trip
        failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
        return counts.Requests >= 5 && failureRatio >= 0.6
    },
    OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
        logger.Get().Warn("Circuit breaker state changed",
            zap.String("name", name),
            zap.String("from", from.String()),
            zap.String("to", to.String()),
        )
    },
}
```

### Monitoring Circuit Breakers

```go
// Get circuit breaker statistics
stats := middleware.GetCircuitBreakerStats()
for endpoint, counts := range stats {
    log.Info("Circuit breaker stats",
        zap.String("endpoint", endpoint),
        zap.Uint32("requests", counts.Requests),
        zap.Uint32("failures", counts.TotalFailures),
        zap.Uint32("successes", counts.TotalSuccesses),
    )
}

// Reset a circuit breaker
success := middleware.ResetCircuitBreaker("GET:/api/users")
if success {
    log.Info("Circuit breaker reset successfully")
}
```

## Health Monitoring

### Health Check Endpoints

```go
// Basic health check
GET /health

// Detailed health with circuit breaker stats
GET /health/detailed

// Circuit breaker statistics
GET /health/circuit-breakers

// Reset specific circuit breaker
POST /health/circuit-breaker/{endpoint}/reset
```

### Health Response Example

```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "uptime": "2h30m15s",
  "system": {
    "go_version": "go1.23.0",
    "goroutines": 25,
    "memory_alloc_mb": 45,
    "memory_total_mb": 120,
    "memory_sys_mb": 75,
    "gc_count": 12,
    "last_gc_time": "2024-01-15T10:29:45Z"
  },
  "circuits": {
    "GET:/api/users": {
      "state": "closed",
      "requests": 150,
      "total_successes": 148,
      "total_failures": 2,
      "consecutive_successes": 25,
      "consecutive_failures": 0,
      "failure_rate": 1.33
    }
  }
}
```

## Best Practices

### 1. Use Structured Logging

```go
// ✅ Good: Structured logging
log.Info("User login successful",
    zap.String("user_id", userID),
    zap.String("ip", clientIP),
    zap.Duration("duration", loginDuration),
)

// ❌ Bad: String formatting
log.Info(fmt.Sprintf("User %s logged in from %s in %v", userID, clientIP, loginDuration))
```

### 2. Include Context Information

```go
// ✅ Good: Include relevant context
log.Error("Database query failed",
    zap.Error(err),
    zap.String("query", "SELECT * FROM users"),
    zap.String("user_id", userID),
    zap.String("operation", "get_user_profile"),
)

// ❌ Bad: Minimal context
log.Error("Query failed", zap.Error(err))
```

### 3. Use Appropriate Log Levels

```go
// ✅ Good: Appropriate levels
log.Debug("Cache miss for key", zap.String("key", cacheKey))
log.Info("User created", zap.String("user_id", userID))
log.Warn("Rate limit approaching", zap.Int("current", current), zap.Int("limit", limit))
log.Error("Failed to send email", zap.Error(err))

// ❌ Bad: Wrong levels
log.Error("User created", zap.String("user_id", userID)) // Should be Info
log.Info("Database connection failed", zap.Error(err))   // Should be Error
```

### 4. Handle Errors Consistently

```go
// ✅ Good: Consistent error handling
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    log := logger.FromCtx(ctx)
    
    if err := s.validateUser(req); err != nil {
        log.Warn("User validation failed", zap.Error(err))
        return nil, errors.NewValidationError("Invalid user data", err)
    }
    
    user, err := s.repo.Create(ctx, req)
    if err != nil {
        log.Error("Failed to create user in database", zap.Error(err))
        return nil, errors.NewDatabaseError("Failed to create user", err)
    }
    
    log.Info("User created successfully", zap.String("user_id", user.ID))
    return user, nil
}
```

### 5. Use Request Context

```go
// ✅ Good: Use context logger in handlers
func (h *UserHandler) GetUser(c *gin.Context) {
    log := logger.FromCtx(c.Request.Context())
    userID := c.Param("id")
    
    log.Info("Fetching user", zap.String("user_id", userID))
    
    user, err := h.service.GetUser(c.Request.Context(), userID)
    if err != nil {
        log.Error("Failed to fetch user", zap.Error(err))
        c.Error(err)
        return
    }
    
    c.JSON(http.StatusOK, user)
}
```

### 6. Monitor Circuit Breakers

```go
// ✅ Good: Monitor circuit breaker health
func (h *HealthHandler) checkCircuitBreakers() {
    stats := h.middleware.GetCircuitBreakerStats()
    
    for endpoint, counts := range stats {
        failureRate := float64(counts.TotalFailures) / float64(counts.Requests) * 100
        
        if failureRate > 50 {
            logger.Get().Warn("High failure rate detected",
                zap.String("endpoint", endpoint),
                zap.Float64("failure_rate", failureRate),
                zap.Uint32("total_requests", counts.Requests),
            )
        }
    }
}
```

### 7. Sensitive Data Handling

```go
// ✅ Good: Avoid logging sensitive data
log.Info("User authentication attempt",
    zap.String("user_id", userID),
    zap.String("ip", clientIP),
    // Don't log password, tokens, etc.
)

// ❌ Bad: Logging sensitive data
log.Info("User login",
    zap.String("password", password), // Never log passwords
    zap.String("token", authToken),   // Never log tokens
)
```

## Examples

### Complete Handler Example

```go
package handler

import (
    "net/http"
    
    "github.com/PhantomX7/dhamma/utility/errors"
    "github.com/PhantomX7/dhamma/utility/logger"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type UserHandler struct {
    service UserService
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    // Get context logger with request information
    log := logger.FromCtx(c.Request.Context())
    
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        log.Warn("Invalid request body", zap.Error(err))
        c.Error(errors.NewValidationError("Invalid request body", err))
        return
    }
    
    log.Info("Creating new user",
        zap.String("email", req.Email),
        zap.String("role", req.Role),
    )
    
    user, err := h.service.CreateUser(c.Request.Context(), req)
    if err != nil {
        log.Error("Failed to create user", zap.Error(err))
        c.Error(err) // Error middleware will handle the response
        return
    }
    
    log.Info("User created successfully",
        zap.String("user_id", user.ID),
        zap.String("email", user.Email),
    )
    
    c.JSON(http.StatusCreated, gin.H{
        "success": true,
        "data": user,
    })
}
```

### Service Layer Example

```go
package service

import (
    "context"
    
    "github.com/PhantomX7/dhamma/utility/errors"
    "github.com/PhantomX7/dhamma/utility/logger"
    "go.uber.org/zap"
)

type UserService struct {
    repo UserRepository
}

func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    log := logger.FromCtx(ctx)
    
    // Validate request
    if err := s.validateCreateUserRequest(req); err != nil {
        log.Warn("User validation failed",
            zap.Error(err),
            zap.String("email", req.Email),
        )
        return nil, errors.NewValidationError("Invalid user data", err)
    }
    
    // Check if user already exists
    existing, err := s.repo.GetByEmail(ctx, req.Email)
    if err != nil && !errors.IsNotFoundError(err) {
        log.Error("Failed to check existing user",
            zap.Error(err),
            zap.String("email", req.Email),
        )
        return nil, errors.NewDatabaseError("Failed to check existing user", err)
    }
    
    if existing != nil {
        log.Warn("User already exists",
            zap.String("email", req.Email),
            zap.String("existing_user_id", existing.ID),
        )
        return nil, errors.NewConflictError("User already exists", nil)
    }
    
    // Create user
    user, err := s.repo.Create(ctx, req)
    if err != nil {
        log.Error("Failed to create user in database",
            zap.Error(err),
            zap.String("email", req.Email),
        )
        return nil, errors.NewDatabaseError("Failed to create user", err)
    }
    
    log.Info("User created successfully",
        zap.String("user_id", user.ID),
        zap.String("email", user.Email),
        zap.String("role", user.Role),
    )
    
    return user, nil
}
```

### Middleware Setup Example

```go
package main

import (
    "github.com/PhantomX7/dhamma/middleware"
    "github.com/gin-gonic/gin"
)

func setupMiddleware(router *gin.Engine, mw *middleware.Middleware) {
    // Recovery middleware (should be first)
    router.Use(mw.Recovery())
    
    // Logger middleware
    router.Use(mw.Logger())
    
    // CORS middleware
    router.Use(mw.CORS())
    
    // Rate limiting
    router.Use(mw.RateLimit())
    
    // Error handler (should be last)
    router.Use(mw.ErrorHandler())
}
```

This documentation provides a comprehensive guide to using the logging and error handling system in the Dhamma project. Follow these patterns for consistent, maintainable, and observable code.