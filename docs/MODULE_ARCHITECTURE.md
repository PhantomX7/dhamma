# Module Architecture Guide

This document outlines the modular architecture patterns used in the Dhamma project, providing guidelines for creating consistent, maintainable, and scalable modules.

## Table of Contents

1. [Overview](#overview)
2. [Module Structure](#module-structure)
3. [Layer Responsibilities](#layer-responsibilities)
4. [Interface Patterns](#interface-patterns)
5. [Dependency Injection](#dependency-injection)
6. [Permission System](#permission-system)
7. [Best Practices](#best-practices)
8. [Module Creation Guide](#module-creation-guide)
9. [Improvements and Recommendations](#improvements-and-recommendations)

## Overview

The Dhamma project follows a **Clean Architecture** approach with a modular design pattern. Each module represents a business domain and follows a consistent layered architecture:

- **Controller Layer**: HTTP request handling and response formatting
- **Service Layer**: Business logic and orchestration
- **Repository Layer**: Data access and persistence
- **DTO Layer**: Data transfer objects for requests and responses
- **Permission Layer**: Access control definitions

## Module Structure

Each module follows this standardized directory structure:

```
modules/
└── {module_name}/
    ├── {module_name}.go          # Interface definitions
    ├── permission.go             # Permission definitions
    ├── controller/
    │   ├── controller.go         # Controller constructor
    │   ├── index.go             # List/pagination endpoint
    │   ├── show.go              # Get single item endpoint
    │   ├── create.go            # Create endpoint
    │   ├── update.go            # Update endpoint
    │   └── {custom_action}.go   # Custom business actions
    ├── service/
    │   ├── service.go           # Service constructor
    │   ├── index.go             # List/pagination logic
    │   ├── show.go              # Get single item logic
    │   ├── create.go            # Create logic
    │   ├── update.go            # Update logic
    │   └── {custom_action}.go   # Custom business logic
    ├── repository/
    │   ├── repository.go        # Repository constructor
    │   └── {custom_query}.go    # Custom database queries
    └── dto/
        ├── request/
        │   ├── create.go        # Create request DTO
        │   ├── update.go        # Update request DTO
        │   └── {custom}.go      # Custom request DTOs
        └── response/
            ├── {module}.go      # Response DTOs
            └── {custom}.go      # Custom response DTOs
```

## Layer Responsibilities

### Controller Layer

**Responsibilities:**
- HTTP request/response handling
- Input validation and parameter extraction
- Response formatting
- Error handling delegation

**Pattern:**
```go
func (c *controller) ActionName(ctx *gin.Context) {
    // 1. Extract and validate parameters
    id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
    if err != nil {
        ctx.Error(&errors.AppError{
            Message: "invalid id",
            Status:  http.StatusBadRequest,
        })
        return
    }

    // 2. Call service layer
    result, err := c.service.ActionName(ctx.Request.Context(), id)
    if err != nil {
        ctx.Error(err)
        return
    }

    // 3. Return formatted response
    ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("success message", result))
}
```

### Service Layer

**Responsibilities:**
- Business logic implementation
- Transaction management
- Cross-module coordination
- Data validation and transformation

**Pattern:**
```go
func (s *service) ActionName(ctx context.Context, params ...) (result, error) {
    // 1. Validate business rules
    if err := s.validateBusinessRules(params); err != nil {
        return result, err
    }

    // 2. Execute business logic (with transaction if needed)
    return s.transactionManager.WithTransaction(ctx, func(tx *gorm.DB) (interface{}, error) {
        // Business logic implementation
        return s.repository.ActionName(ctx, params, tx)
    })
}
```

### Repository Layer

**Responsibilities:**
- Data access and persistence
- Database query implementation
- Entity mapping
- Base repository pattern usage

**Pattern:**
```go
type repository struct {
    base baseRepo.BaseRepositoryInterface[entity.EntityName]
    db   *gorm.DB
}

func (r *repository) CustomQuery(ctx context.Context, params ...) (result, error) {
    // Custom database operations
    // Use base repository for standard CRUD operations
    return r.base.FindByCondition(ctx, conditions)
}
```

## Interface Patterns

Each module defines three main interfaces in its root file (`{module_name}.go`):

### Repository Interface
```go
type Repository interface {
    repository.BaseRepositoryInterface[entity.EntityName]
    // Custom repository methods
    CustomMethod(ctx context.Context, params ...) (result, error)
}
```

### Service Interface
```go
type Service interface {
    Index(ctx context.Context, pg *pagination.Pagination) ([]entity.EntityName, utility.PaginationMeta, error)
    Show(ctx context.Context, id uint64) (entity.EntityName, error)
    Create(ctx context.Context, request request.CreateRequest) (entity.EntityName, error)
    Update(ctx context.Context, id uint64, request request.UpdateRequest) (entity.EntityName, error)
    // Custom service methods
    CustomAction(ctx context.Context, params ...) (result, error)
}
```

### Controller Interface
```go
type Controller interface {
    Index(ctx *gin.Context)
    Show(ctx *gin.Context)
    Create(ctx *gin.Context)
    Update(ctx *gin.Context)
    // Custom controller methods
    CustomAction(ctx *gin.Context)
}
```

## Dependency Injection

The project uses **Uber FX** for dependency injection. Each layer has its own module file:

### Controller Module (`controller_module.go`)
```go
var ControllerModule = fx.Options(
    fx.Provide(
        moduleController.New,
        // ... other controllers
    ),
)
```

### Service Module (`service_module.go`)
```go
var ServiceModule = fx.Options(
    fx.Provide(
        moduleService.New,
        // ... other services
    ),
)
```

### Repository Module (`repository_module.go`)
```go
var RepositoryModule = fx.Options(
    fx.Provide(
        moduleRepository.New,
        // ... other repositories
    ),
)
```

## Permission System

Each module defines its permissions in `permission.go`:

```go
type permission struct {
    Key string
    // Standard CRUD permissions
    Index  string // List/view all
    Show   string // View single item
    Create string // Create new item
    Update string // Update existing item
    Delete string // Delete item
    // Custom permissions
    CustomAction string // Custom business action
}

var Permissions = permission{
    Key:          "module-name",
    Index:        "index",
    Show:         "show",
    Create:       "create",
    Update:       "update",
    Delete:       "delete",
    CustomAction: "custom-action",
}
```

## Best Practices

### 1. Interface Segregation
- Keep interfaces focused and cohesive
- Separate read and write operations when appropriate
- Use composition for complex interfaces

### 2. Error Handling
- Use structured error types (`errors.AppError`, `errors.RepositoryError`)
- Provide meaningful error messages
- Handle errors at the appropriate layer

### 3. Context Usage
- Always pass `context.Context` as the first parameter
- Use context for cancellation and timeouts
- Propagate context through all layers

### 4. Transaction Management
- Use transaction manager for multi-step operations
- Keep transactions as short as possible
- Handle rollbacks properly

### 5. Validation
- Validate input at the controller layer
- Implement business rule validation in the service layer
- Use struct tags for basic validation

### 6. Testing
- Write unit tests for each layer
- Use mocks for dependencies
- Test error scenarios
- Follow TDD approach

## Module Creation Guide

### Step 1: Create Module Structure
```bash
mkdir -p modules/{module_name}/{controller,service,repository,dto/{request,response}}
```

### Step 2: Define Interfaces
Create `modules/{module_name}/{module_name}.go` with Repository, Service, and Controller interfaces.

### Step 3: Define Permissions
Create `modules/{module_name}/permission.go` with permission definitions.

### Step 4: Implement Repository
- Create repository struct with base repository composition
- Implement custom database operations
- Add constructor function

### Step 5: Implement Service
- Create service struct with repository dependency
- Implement business logic
- Add transaction management where needed
- Add constructor function

### Step 6: Implement Controller
- Create controller struct with service dependency
- Implement HTTP handlers
- Add proper error handling
- Add constructor function

### Step 7: Create DTOs
- Define request DTOs with validation tags
- Define response DTOs for API responses

### Step 8: Register Dependencies
- Add to `controller_module.go`
- Add to `service_module.go`
- Add to `repository_module.go`

### Step 9: Add Routes
- Create route definitions in `routes/` directory
- Add permission middleware where needed

### Step 10: Write Tests
- Create unit tests for each layer
- Add integration tests
- Test error scenarios

## Improvements and Recommendations

### 1. **Add Response DTOs**
**Current Issue**: Most modules lack dedicated response DTOs, returning entities directly.

**Recommendation**: Create response DTOs to:
- Control API response structure
- Hide internal entity details
- Enable API versioning
- Improve security by preventing data leaks

```go
// dto/response/module_response.go
type ModuleResponse struct {
    ID        uint64    `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
    // Only expose necessary fields
}
```

### 2. **Standardize Error Handling**
**Current Issue**: Inconsistent error handling across modules.

**Recommendation**: 
- Create module-specific error types
- Use error wrapping for better context
- Implement error translation layer

```go
// errors/module_errors.go
var (
    ErrModuleNotFound = errors.NewAppError("module not found", http.StatusNotFound)
    ErrInvalidModuleData = errors.NewAppError("invalid module data", http.StatusBadRequest)
)
```

### 3. **Add Input Validation Layer**
**Current Issue**: Limited input validation in DTOs.

**Recommendation**:
- Use comprehensive validation tags
- Create custom validators for business rules
- Add validation middleware

```go
type CreateRequest struct {
    Name        string `json:"name" validate:"required,min=3,max=100"`
    Email       string `json:"email" validate:"required,email"`
    Description string `json:"description" validate:"max=500"`
}
```

### 4. **Implement Caching Layer**
**Recommendation**: Add caching for frequently accessed data:
- Redis integration
- Cache-aside pattern
- TTL-based invalidation

### 5. **Add Audit Trail**
**Recommendation**: Implement audit logging:
- Track entity changes
- User action logging
- Compliance requirements

### 6. **Enhance Testing Structure**
**Recommendation**:
- Add test utilities package
- Create test data factories
- Implement table-driven tests
- Add benchmark tests

### 7. **API Documentation**
**Recommendation**:
- Add Swagger/OpenAPI documentation
- Generate docs from code annotations
- Include request/response examples

### 8. **Monitoring and Metrics**
**Recommendation**:
- Add Prometheus metrics
- Implement distributed tracing
- Add performance monitoring

### 9. **Configuration Management**
**Recommendation**:
- Module-specific configuration
- Environment-based config
- Configuration validation

### 10. **Database Migration Management**
**Recommendation**:
- Module-specific migrations
- Migration rollback support
- Schema versioning

This architecture provides a solid foundation for building scalable, maintainable applications while following clean architecture principles and Go best practices.