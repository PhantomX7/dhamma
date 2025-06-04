# Module Template Guide

This guide provides practical templates and step-by-step instructions for creating new modules in the Dhamma project.

## Quick Start Template

Use this template to quickly scaffold a new module:

### 1. Module Interface Definition

**File**: `modules/{module_name}/{module_name}.go`

```go
package {module_name}

import (
	"context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/{module_name}/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
	"github.com/PhantomX7/dhamma/utility/repository"
)

// Repository interface defines data access methods
type Repository interface {
	repository.BaseRepositoryInterface[entity.{EntityName}]
	// Add custom repository methods here
	// Example: FindByStatus(ctx context.Context, status string) ([]entity.{EntityName}, error)
}

// Service interface defines business logic methods
type Service interface {
	Index(ctx context.Context, pg *pagination.Pagination) ([]entity.{EntityName}, utility.PaginationMeta, error)
	Show(ctx context.Context, id uint64) (entity.{EntityName}, error)
	Create(ctx context.Context, request request.{EntityName}CreateRequest) (entity.{EntityName}, error)
	Update(ctx context.Context, id uint64, request request.{EntityName}UpdateRequest) (entity.{EntityName}, error)
	// Add custom service methods here
	// Example: ActivateItem(ctx context.Context, id uint64) (entity.{EntityName}, error)
}

// Controller interface defines HTTP handlers
type Controller interface {
	Index(ctx *gin.Context)
	Show(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	// Add custom controller methods here
	// Example: Activate(ctx *gin.Context)
}
```

### 2. Permission Definition

**File**: `modules/{module_name}/permission.go`

```go
package {module_name}

// permission defines all permissions for this module
type permission struct {
	Key string
	// Standard CRUD permissions
	Index  string // List all items
	Show   string // View single item
	Create string // Create new item
	Update string // Update existing item
	Delete string // Delete item
	// Add custom permissions here
	// Example: Activate string // Activate item
}

// Permissions contains all permission definitions for the {module_name} module
var Permissions = permission{
	Key:    "{module-name}", // kebab-case
	Index:  "index",
	Show:   "show",
	Create: "create",
	Update: "update",
	Delete: "delete",
	// Example: Activate: "activate",
}
```

### 3. Repository Implementation

**File**: `modules/{module_name}/repository/repository.go`

```go
package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/{module_name}"
	"github.com/PhantomX7/dhamma/utility/pagination"
	baseRepo "github.com/PhantomX7/dhamma/utility/repository"
)

// repository implements the {module_name}.Repository interface
type repository struct {
	base baseRepo.BaseRepositoryInterface[entity.{EntityName}]
	db   *gorm.DB
}

// New creates a new {module_name} repository instance
func New(db *gorm.DB) {module_name}.Repository {
	return &repository{
		base: baseRepo.NewBaseRepository[entity.{EntityName}](db),
		db:   db,
	}
}

// FindAll retrieves all entities with pagination
func (r *repository) FindAll(ctx context.Context, pg *pagination.Pagination) ([]entity.{EntityName}, error) {
	return r.base.FindAll(ctx, pg)
}

// FindByID retrieves an entity by its ID
func (r *repository) FindByID(ctx context.Context, id uint64, preloads ...string) (entity.{EntityName}, error) {
	return r.base.FindByID(ctx, id, preloads...)
}

// Create creates a new entity
func (r *repository) Create(ctx context.Context, entity *entity.{EntityName}, tx *gorm.DB) error {
	return r.base.Create(ctx, entity, tx)
}

// Update updates an existing entity
func (r *repository) Update(ctx context.Context, entity *entity.{EntityName}, tx *gorm.DB) error {
	return r.base.Update(ctx, entity, tx)
}

// Delete deletes an entity
func (r *repository) Delete(ctx context.Context, entity *entity.{EntityName}, tx *gorm.DB) error {
	return r.base.Delete(ctx, entity, tx)
}

// Count returns the total count of entities
func (r *repository) Count(ctx context.Context) (int64, error) {
	return r.base.Count(ctx)
}

// Example custom repository method
// func (r *repository) FindByStatus(ctx context.Context, status string) ([]entity.{EntityName}, error) {
// 	var entities []entity.{EntityName}
// 	err := r.db.WithContext(ctx).Where("status = ?", status).Find(&entities).Error
// 	return entities, err
// }
```

### 4. Service Implementation

**File**: `modules/{module_name}/service/service.go`

```go
package service

import (
	"github.com/PhantomX7/dhamma/libs/transaction_manager"
	"github.com/PhantomX7/dhamma/modules/{module_name}"
)

// service implements the {module_name}.Service interface
type service struct {
	{module_name}Repo      {module_name}.Repository
	transactionManager transaction_manager.Client
}

// New creates a new {module_name} service instance
func New(
	{module_name}Repo {module_name}.Repository,
	transactionManager transaction_manager.Client,
) {module_name}.Service {
	return &service{
		{module_name}Repo:      {module_name}Repo,
		transactionManager: transactionManager,
	}
}
```

**File**: `modules/{module_name}/service/index.go`

```go
package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
)

// Index retrieves a paginated list of entities
func (s *service) Index(ctx context.Context, pg *pagination.Pagination) ([]entity.{EntityName}, utility.PaginationMeta, error) {
	// Get entities with pagination
	entities, err := s.{module_name}Repo.FindAll(ctx, pg)
	if err != nil {
		return nil, utility.PaginationMeta{}, err
	}

	// Get total count for pagination metadata
	totalCount, err := s.{module_name}Repo.Count(ctx)
	if err != nil {
		return nil, utility.PaginationMeta{}, err
	}

	// Build pagination metadata
	paginationMeta := utility.BuildPaginationMeta(pg, totalCount)

	return entities, paginationMeta, nil
}
```

**File**: `modules/{module_name}/service/show.go`

```go
package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility/errors"
)

// Show retrieves a single entity by ID
func (s *service) Show(ctx context.Context, id uint64) (entity.{EntityName}, error) {
	entity, err := s.{module_name}Repo.FindByID(ctx, id)
	if err != nil {
		return entity, errors.NewRepositoryError("failed to find {module_name}", err)
	}

	return entity, nil
}
```

**File**: `modules/{module_name}/service/create.go`

```go
package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/{module_name}/dto/request"
	"github.com/PhantomX7/dhamma/utility/errors"
)

// Create creates a new entity
func (s *service) Create(ctx context.Context, req request.{EntityName}CreateRequest) (entity.{EntityName}, error) {
	// Create entity from request
	newEntity := entity.{EntityName}{
		// Map request fields to entity
		// Example: Name: req.Name,
	}

	// Execute creation within transaction
	result, err := s.transactionManager.WithTransaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if err := s.{module_name}Repo.Create(ctx, &newEntity, tx); err != nil {
			return nil, errors.NewRepositoryError("failed to create {module_name}", err)
		}
		return newEntity, nil
	})

	if err != nil {
		return entity.{EntityName}{}, err
	}

	return result.(entity.{EntityName}), nil
}
```

**File**: `modules/{module_name}/service/update.go`

```go
package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/{module_name}/dto/request"
	"github.com/PhantomX7/dhamma/utility/errors"
)

// Update updates an existing entity
func (s *service) Update(ctx context.Context, id uint64, req request.{EntityName}UpdateRequest) (entity.{EntityName}, error) {
	// Find existing entity
	existingEntity, err := s.{module_name}Repo.FindByID(ctx, id)
	if err != nil {
		return entity.{EntityName}{}, errors.NewRepositoryError("failed to find {module_name}", err)
	}

	// Update entity fields from request
	// Example: existingEntity.Name = req.Name

	// Execute update within transaction
	result, err := s.transactionManager.WithTransaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if err := s.{module_name}Repo.Update(ctx, &existingEntity, tx); err != nil {
			return nil, errors.NewRepositoryError("failed to update {module_name}", err)
		}
		return existingEntity, nil
	})

	if err != nil {
		return entity.{EntityName}{}, err
	}

	return result.(entity.{EntityName}), nil
}
```

### 5. Controller Implementation

**File**: `modules/{module_name}/controller/controller.go`

```go
package controller

import (
	"github.com/PhantomX7/dhamma/modules/{module_name}"
)

// controller implements the {module_name}.Controller interface
type controller struct {
	{module_name}Service {module_name}.Service
}

// New creates a new {module_name} controller instance
func New({module_name}Service {module_name}.Service) {module_name}.Controller {
	return &controller{
		{module_name}Service: {module_name}Service,
	}
}
```

**File**: `modules/{module_name}/controller/index.go`

```go
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
)

// Index handles GET /{module_name}s endpoint
func (c *controller) Index(ctx *gin.Context) {
	// Parse pagination parameters
	pg := pagination.NewFromGinContext(ctx)

	// Get entities from service
	entities, paginationMeta, err := c.{module_name}Service.Index(ctx.Request.Context(), pg)
	if err != nil {
		ctx.Error(err)
		return
	}

	// Return paginated response
	ctx.JSON(http.StatusOK, utility.BuildResponseSuccessWithMeta(
		"{module_name}s retrieved successfully",
		entities,
		paginationMeta,
	))
}
```

**File**: `modules/{module_name}/controller/show.go`

```go
package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
)

// Show handles GET /{module_name}s/:id endpoint
func (c *controller) Show(ctx *gin.Context) {
	// Parse ID parameter
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Error(&errors.AppError{
			Message: "invalid {module_name} id",
			Status:  http.StatusBadRequest,
		})
		return
	}

	// Get entity from service
	entity, err := c.{module_name}Service.Show(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess(
		"{module_name} retrieved successfully",
		entity,
	))
}
```

**File**: `modules/{module_name}/controller/create.go`

```go
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/PhantomX7/dhamma/modules/{module_name}/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
)

// Create handles POST /{module_name}s endpoint
func (c *controller) Create(ctx *gin.Context) {
	// Parse request body
	var req request.{EntityName}CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(&errors.AppError{
			Message: "invalid request body",
			Status:  http.StatusBadRequest,
		})
		return
	}

	// Create entity via service
	entity, err := c.{module_name}Service.Create(ctx.Request.Context(), req)
	if err != nil {
		ctx.Error(err)
		return
	}

	// Return created response
	ctx.JSON(http.StatusCreated, utility.BuildResponseSuccess(
		"{module_name} created successfully",
		entity,
	))
}
```

**File**: `modules/{module_name}/controller/update.go`

```go
package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/PhantomX7/dhamma/modules/{module_name}/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
)

// Update handles PUT /{module_name}s/:id endpoint
func (c *controller) Update(ctx *gin.Context) {
	// Parse ID parameter
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Error(&errors.AppError{
			Message: "invalid {module_name} id",
			Status:  http.StatusBadRequest,
		})
		return
	}

	// Parse request body
	var req request.{EntityName}UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(&errors.AppError{
			Message: "invalid request body",
			Status:  http.StatusBadRequest,
		})
		return
	}

	// Update entity via service
	entity, err := c.{module_name}Service.Update(ctx.Request.Context(), id, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	// Return updated response
	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess(
		"{module_name} updated successfully",
		entity,
	))
}
```

### 6. DTO Definitions

**File**: `modules/{module_name}/dto/request/create.go`

```go
package request

// {EntityName}CreateRequest defines the structure for creating a new {module_name}
type {EntityName}CreateRequest struct {
	// Add fields with validation tags
	// Example:
	// Name        string `json:"name" validate:"required,min=3,max=100"`
	// Description string `json:"description" validate:"max=500"`
	// Status      string `json:"status" validate:"required,oneof=active inactive"`
}
```

**File**: `modules/{module_name}/dto/request/update.go`

```go
package request

// {EntityName}UpdateRequest defines the structure for updating an existing {module_name}
type {EntityName}UpdateRequest struct {
	// Add fields with validation tags (usually same as create but with optional fields)
	// Example:
	// Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	// Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
	// Status      *string `json:"status,omitempty" validate:"omitempty,oneof=active inactive"`
}
```

### 7. Module Registration

Add your module to the dependency injection modules:

**In**: `modules/repository_module.go`
```go
// Add import
{module_name}Repository "github.com/PhantomX7/dhamma/modules/{module_name}/repository"

// Add to fx.Provide
{module_name}Repository.New,
```

**In**: `modules/service_module.go`
```go
// Add import
{module_name}Service "github.com/PhantomX7/dhamma/modules/{module_name}/service"

// Add to fx.Provide
{module_name}Service.New,
```

**In**: `modules/controller_module.go`
```go
// Add import
{module_name}Controller "github.com/PhantomX7/dhamma/modules/{module_name}/controller"

// Add to fx.Provide
{module_name}Controller.New,
```

### 8. Route Definition

Create route files in `routes/admin/` and `routes/domain/`:

**File**: `routes/admin/{module_name}_route.go`

```go
package admin

import (
	"github.com/gin-gonic/gin"

	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/modules/{module_name}"
)

// {EntityName}Routes sets up routes for {module_name} module
func {EntityName}Routes(
	router *gin.RouterGroup,
	{module_name}Controller {module_name}.Controller,
) {
	{module_name}Group := router.Group("/{module_name}s")
	{
		// GET /{module_name}s - List all {module_name}s
		{module_name}Group.GET("",
			middleware.Permission("{module_name}", "{module_name}.index"),
			{module_name}Controller.Index,
		)

		// GET /{module_name}s/:id - Get single {module_name}
		{module_name}Group.GET("/:id",
			middleware.Permission("{module_name}", "{module_name}.show"),
			{module_name}Controller.Show,
		)

		// POST /{module_name}s - Create new {module_name}
		{module_name}Group.POST("",
			middleware.Permission("{module_name}", "{module_name}.create"),
			{module_name}Controller.Create,
		)

		// PUT /{module_name}s/:id - Update {module_name}
		{module_name}Group.PUT("/:id",
			middleware.Permission("{module_name}", "{module_name}.update"),
			{module_name}Controller.Update,
		)

		// Add custom routes here
		// Example:
		// {module_name}Group.POST("/:id/activate",
		// 	middleware.Permission("{module_name}", "{module_name}.activate"),
		// 	{module_name}Controller.Activate,
		// )
	}
}
```

## Automation Script

Create a script to automate module generation:

**File**: `tools/modgen/main.go`

```go
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <module_name>")
	}

	moduleName := os.Args[1]
	entityName := strings.Title(moduleName)

	// Create module directory structure
	createDirStructure(moduleName)

	// Generate files from templates
	generateFiles(moduleName, entityName)

	fmt.Printf("Module '%s' generated successfully!\n", moduleName)
	fmt.Println("Don't forget to:")
	fmt.Println("1. Add entity definition in entity/ directory")
	fmt.Println("2. Register module in dependency injection files")
	fmt.Println("3. Add routes in routes/ directory")
	fmt.Println("4. Run database migration if needed")
}

func createDirStructure(moduleName string) {
	dirs := []string{
		fmt.Sprintf("modules/%s", moduleName),
		fmt.Sprintf("modules/%s/controller", moduleName),
		fmt.Sprintf("modules/%s/service", moduleName),
		fmt.Sprintf("modules/%s/repository", moduleName),
		fmt.Sprintf("modules/%s/dto/request", moduleName),
		fmt.Sprintf("modules/%s/dto/response", moduleName),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}
}

func generateFiles(moduleName, entityName string) {
	// Implementation would generate files from templates
	// This is a simplified version - you would expand this
	// to generate all the template files shown above
}
```

## Testing Template

Create comprehensive tests for your module:

**File**: `modules/{module_name}/service/service_test.go`

```go
package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/{module_name}/mocks"
	"github.com/PhantomX7/dhamma/modules/{module_name}/service"
	"github.com/PhantomX7/dhamma/utility/pagination"
)

func TestService_Index(t *testing.T) {
	// Setup
	mockRepo := mocks.NewRepository(t)
	mockTxManager := mocks.NewTransactionManager(t)
	svc := service.New(mockRepo, mockTxManager)

	ctx := context.Background()
	pg := &pagination.Pagination{Page: 1, Limit: 10}

	// Mock expectations
	expectedEntities := []entity.{EntityName}{
		{ID: 1, Name: "Test 1"},
		{ID: 2, Name: "Test 2"},
	}
	mockRepo.On("FindAll", ctx, pg).Return(expectedEntities, nil)
	mockRepo.On("Count", ctx).Return(int64(2), nil)

	// Execute
	entities, meta, err := svc.Index(ctx, pg)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedEntities, entities)
	assert.Equal(t, int64(2), meta.Total)
	mockRepo.AssertExpectations(t)
}
```

This template provides a solid foundation for creating consistent, well-structured modules in the Dhamma project. Remember to customize the templates based on your specific business requirements and add appropriate validation, error handling, and business logic.