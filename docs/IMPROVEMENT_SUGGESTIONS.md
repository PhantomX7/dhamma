# Module Architecture Improvement Suggestions

Based on the analysis of the current module patterns in the Dhamma project, here are comprehensive suggestions for improvements across different aspects of the architecture.

## 1. Code Organization & Structure

### Current Issues
- Inconsistent file organization across modules
- Mixed responsibilities in some files
- Lack of standardized naming conventions

### Improvements

#### A. Standardize Directory Structure
```
modules/{module_name}/
├── {module_name}.go          # Interface definitions
├── permission.go             # Permission constants
├── controller/
│   ├── controller.go         # Controller struct and constructor
│   ├── index.go             # GET /items
│   ├── show.go              # GET /items/:id
│   ├── create.go            # POST /items
│   ├── update.go            # PUT /items/:id
│   ├── delete.go            # DELETE /items/:id
│   └── {custom_action}.go   # Custom endpoints
├── service/
│   ├── service.go           # Service struct and constructor
│   ├── index.go             # List logic
│   ├── show.go              # Show logic
│   ├── create.go            # Create logic
│   ├── update.go            # Update logic
│   ├── delete.go            # Delete logic
│   └── {custom_logic}.go    # Custom business logic
├── repository/
│   ├── repository.go        # Repository struct and constructor
│   └── {custom_queries}.go  # Custom database queries
├── dto/
│   ├── request/
│   │   ├── create.go        # Create request DTOs
│   │   ├── update.go        # Update request DTOs
│   │   └── {custom}.go      # Custom request DTOs
│   └── response/
│       ├── detail.go        # Detailed response DTOs
│       ├── list.go          # List response DTOs
│       └── {custom}.go      # Custom response DTOs
├── validator/
│   └── {module_name}.go     # Custom validation rules
└── tests/
    ├── controller_test.go   # Controller tests
    ├── service_test.go      # Service tests
    └── repository_test.go   # Repository tests
```

#### B. Implement Response DTOs
Currently, entities are returned directly. Implement response DTOs for better API design:

```go
// modules/{module_name}/dto/response/detail.go
package response

import "time"

type {EntityName}DetailResponse struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	// Exclude sensitive fields like internal IDs, etc.
}

func From{EntityName}(entity entity.{EntityName}) {EntityName}DetailResponse {
	return {EntityName}DetailResponse{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		Status:      entity.Status,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
```

## 2. Error Handling & Logging

### Current Issues
- Inconsistent error handling patterns
- Limited error context
- No structured logging in some areas

### Improvements

#### A. Standardized Error Types
```go
// utility/errors/module_errors.go
package errors

type ModuleError struct {
	Module    string
	Operation string
	Message   string
	Cause     error
	Context   map[string]interface{}
}

func NewModuleError(module, operation, message string, cause error) *ModuleError {
	return &ModuleError{
		Module:    module,
		Operation: operation,
		Message:   message,
		Cause:     cause,
		Context:   make(map[string]interface{}),
	}
}

func (e *ModuleError) WithContext(key string, value interface{}) *ModuleError {
	e.Context[key] = value
	return e
}
```

#### B. Enhanced Service Error Handling
```go
// Example in service layer
func (s *service) Show(ctx context.Context, id uint64) (entity.{EntityName}, error) {
	logger := s.logger.WithFields(logrus.Fields{
		"module":    "{module_name}",
		"operation": "show",
		"entity_id": id,
	})

	logger.Info("retrieving entity")

	entity, err := s.{module_name}Repo.FindByID(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to retrieve entity")
		return entity, errors.NewModuleError(
			"{module_name}", "show", "failed to retrieve entity", err,
		).WithContext("entity_id", id)
	}

	logger.Info("entity retrieved successfully")
	return entity, nil
}
```

## 3. Validation & Security

### Current Issues
- Basic validation only
- No input sanitization
- Limited security checks

### Improvements

#### A. Enhanced Validation
```go
// modules/{module_name}/validator/{module_name}.go
package validator

import (
	"context"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/PhantomX7/dhamma/modules/{module_name}"
)

type {EntityName}Validator struct {
	repo {module_name}.Repository
}

func New{EntityName}Validator(repo {module_name}.Repository) *{EntityName}Validator {
	return &{EntityName}Validator{repo: repo}
}

// Custom validation for unique name
func (v *{EntityName}Validator) ValidateUniqueName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	ctx := context.Background()
	
	// Check if name already exists
	exists, err := v.repo.ExistsByName(ctx, name)
	if err != nil {
		return false // Fail safe
	}
	return !exists
}

// Sanitize input
func SanitizeString(input string) string {
	// Remove potentially harmful characters
	re := regexp.MustCompile(`[<>"'&]`)
	return re.ReplaceAllString(input, "")
}
```

#### B. Request Validation Middleware
```go
// middleware/validation.go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
)

func ValidateRequest[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(&errors.AppError{
				Message: "Invalid request format",
				Status:  http.StatusBadRequest,
				Details: err.Error(),
			})
			return
		}

		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			c.Error(&errors.AppError{
				Message: "Validation failed",
				Status:  http.StatusBadRequest,
				Details: utility.FormatValidationErrors(err),
			})
			return
		}

		c.Set("validated_request", req)
		c.Next()
	}
}
```

## 4. Performance Optimization

### Current Issues
- No caching strategy
- N+1 query problems in some areas
- No database query optimization

### Improvements

#### A. Implement Caching Layer
```go
// utility/cache/cache.go
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	InvalidatePattern(ctx context.Context, pattern string) error
}

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) Cache {
	return &redisCache{client: client}
}

func (r *redisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (r *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}
```

#### B. Cached Service Layer
```go
// modules/{module_name}/service/cached_service.go
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/{module_name}"
	"github.com/PhantomX7/dhamma/utility/cache"
)

type cachedService struct {
	baseService {module_name}.Service
	cache       cache.Cache
	cacheTTL    time.Duration
}

func NewCachedService(baseService {module_name}.Service, cache cache.Cache) {module_name}.Service {
	return &cachedService{
		baseService: baseService,
		cache:       cache,
		cacheTTL:    15 * time.Minute,
	}
}

func (s *cachedService) Show(ctx context.Context, id uint64) (entity.{EntityName}, error) {
	cacheKey := fmt.Sprintf("{module_name}:show:%d", id)
	
	// Try to get from cache first
	var cached entity.{EntityName}
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}

	// Get from database
	entity, err := s.baseService.Show(ctx, id)
	if err != nil {
		return entity, err
	}

	// Cache the result
	_ = s.cache.Set(ctx, cacheKey, entity, s.cacheTTL)

	return entity, nil
}
```

#### C. Database Query Optimization
```go
// modules/{module_name}/repository/optimized_queries.go
package repository

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility/pagination"
)

// FindAllWithRelations optimizes queries with proper preloading
func (r *repository) FindAllWithRelations(ctx context.Context, pg *pagination.Pagination, relations []string) ([]entity.{EntityName}, error) {
	var entities []entity.{EntityName}
	
	query := r.db.WithContext(ctx)
	
	// Preload specified relations to avoid N+1 queries
	for _, relation := range relations {
		query = query.Preload(relation)
	}
	
	// Apply pagination
	offset := (pg.Page - 1) * pg.Limit
	err := query.Offset(offset).Limit(pg.Limit).Find(&entities).Error
	
	return entities, err
}

// FindByIDsOptimized batch fetch multiple entities
func (r *repository) FindByIDsOptimized(ctx context.Context, ids []uint64) ([]entity.{EntityName}, error) {
	var entities []entity.{EntityName}
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&entities).Error
	return entities, err
}
```

## 5. Testing Strategy

### Current Issues
- Limited test coverage
- No integration tests
- Missing mock implementations

### Improvements

#### A. Comprehensive Test Structure
```go
// modules/{module_name}/tests/service_test.go
package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/PhantomX7/dhamma/modules/{module_name}/mocks"
	"github.com/PhantomX7/dhamma/modules/{module_name}/service"
)

type ServiceTestSuite struct {
	suite.Suite
	mockRepo       *mocks.Repository
	mockTxManager  *mocks.TransactionManager
	mockCache      *mocks.Cache
	service        {module_name}.Service
}

func (suite *ServiceTestSuite) SetupTest() {
	suite.mockRepo = mocks.NewRepository(suite.T())
	suite.mockTxManager = mocks.NewTransactionManager(suite.T())
	suite.mockCache = mocks.NewCache(suite.T())
	suite.service = service.New(
		suite.mockRepo,
		suite.mockTxManager,
	)
}

func (suite *ServiceTestSuite) TestShow_Success() {
	// Test implementation
}

func (suite *ServiceTestSuite) TestShow_NotFound() {
	// Test implementation
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
```

#### B. Integration Tests
```go
// tests/integration/{module_name}_test.go
package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/PhantomX7/dhamma/config"
	"github.com/PhantomX7/dhamma/database"
)

func Test{EntityName}CRUD_Integration(t *testing.T) {
	// Setup test database
	db := database.SetupTestDB()
	defer database.CleanupTestDB(db)

	// Setup test server
	router := setupTestRouter(db)

	// Test Create
	createReq := map[string]interface{}{
		"name": "Test Item",
	}
	createBody, _ := json.Marshal(createReq)
	createResp := httptest.NewRecorder()
	createHttpReq, _ := http.NewRequest("POST", "/{module_name}s", bytes.NewBuffer(createBody))
	router.ServeHTTP(createResp, createHttpReq)

	assert.Equal(t, http.StatusCreated, createResp.Code)

	// Test Read, Update, Delete...
}
```

## 6. Monitoring & Observability

### Improvements

#### A. Metrics Collection
```go
// utility/metrics/metrics.go
package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests.",
		},
		[]string{"module", "method", "endpoint", "status"},
	)

	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "database_query_duration_seconds",
			Help: "Duration of database queries.",
		},
		[]string{"module", "operation"},
	)
)

func RecordRequestDuration(module, method, endpoint, status string, duration time.Duration) {
	RequestDuration.WithLabelValues(module, method, endpoint, status).Observe(duration.Seconds())
}
```

#### B. Health Checks
```go
// utility/health/health.go
package health

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type HealthChecker struct {
	db *gorm.DB
}

func NewHealthChecker(db *gorm.DB) *HealthChecker {
	return &HealthChecker{db: db}
}

func (h *HealthChecker) CheckDatabase(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	sqlDB, err := h.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.PingContext(ctx)
}
```

## 7. Configuration Management

### Improvements

#### A. Module-Specific Configuration
```go
// config/modules.go
package config

type ModuleConfig struct {
	CacheEnabled bool          `mapstructure:"cache_enabled"`
	CacheTTL     time.Duration `mapstructure:"cache_ttl"`
	RateLimit    int           `mapstructure:"rate_limit"`
	MaxPageSize  int           `mapstructure:"max_page_size"`
}

type ModulesConfig struct {
	Auth         ModuleConfig `mapstructure:"auth"`
	Domain       ModuleConfig `mapstructure:"domain"`
	ChatTemplate ModuleConfig `mapstructure:"chat_template"`
	// Add other modules
}
```

## 8. API Documentation

### Improvements

#### A. OpenAPI/Swagger Integration
```go
// Add swagger annotations to controllers
// @Summary List {module_name}s
// @Description Get a paginated list of {module_name}s
// @Tags {module_name}
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utility.Response{data=[]entity.{EntityName}}
// @Failure 400 {object} utility.ErrorResponse
// @Router /{module_name}s [get]
func (c *controller) Index(ctx *gin.Context) {
	// Implementation
}
```

## Implementation Priority

1. **High Priority**
   - Standardize directory structure
   - Implement response DTOs
   - Enhanced error handling
   - Basic caching layer

2. **Medium Priority**
   - Comprehensive testing
   - Performance optimizations
   - Validation improvements
   - Monitoring setup

3. **Low Priority**
   - Advanced caching strategies
   - Complex health checks
   - Advanced metrics
   - API documentation automation

## Migration Strategy

1. **Phase 1**: Create templates and standards
2. **Phase 2**: Migrate one module as a reference
3. **Phase 3**: Gradually migrate other modules
4. **Phase 4**: Implement advanced features

This comprehensive improvement plan will enhance the maintainability, performance, and reliability of the Dhamma project's module architecture.