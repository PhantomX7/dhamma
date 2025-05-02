package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// ModuleConfig holds the configuration for generating a module
type ModuleConfig struct {
	// Module name in snake_case (e.g., "user_profile")
	ModuleName string
	// Module name in camelCase (e.g., "userProfile")
	ModuleNameCamel string
	// Module name in PascalCase (e.g., "UserProfile")
	ModuleNamePascal string
	// Table name for the module (e.g., "user_profiles")
	TableName string
	// Entity name in PascalCase (e.g., "UserProfile")
	EntityName string
	// Fields for the entity
	Fields []FieldConfig
}

// FieldConfig holds the configuration for a field in the entity
type FieldConfig struct {
	// Field name in camelCase (e.g., "firstName")
	Name string
	// Field name in snake_case (e.g., "first_name")
	NameSnake string
	// Field type (e.g., "string", "uint64", "bool")
	Type string
	// Whether the field is required in create request
	Required bool
	// Whether the field should be unique
	Unique bool
	// Description of the field
	Description string
}

// main is the entry point for the module generator
func main() {
	// Parse command line arguments
	moduleName := flag.String("module", "", "Module name in snake_case (required)")
	tableName := flag.String("table", "", "Table name (defaults to plural of module name)")
	entityName := flag.String("entity", "", "Entity name in PascalCase (defaults to PascalCase of module name)")
	outputDir := flag.String("output", "../../modules", "Output directory for the generated module")
	flag.Parse()

	// Validate required arguments
	if *moduleName == "" {
		fmt.Println("Error: Module name is required")
		flag.Usage()
		os.Exit(1)
	}

	// Set default values if not provided
	if *tableName == "" {
		*tableName = *moduleName + "s"
	}
	if *entityName == "" {
		*entityName = snakeToPascal(*moduleName)
	}

	// Create module config
	config := ModuleConfig{
		ModuleName:       *moduleName,
		ModuleNameCamel:  snakeToCamel(*moduleName),
		ModuleNamePascal: snakeToPascal(*moduleName),
		TableName:        *tableName,
		EntityName:       *entityName,
		Fields: []FieldConfig{
			{
				Name:        "name",
				NameSnake:   "name",
				Type:        "string",
				Required:    true,
				Unique:      true,
				Description: "Name of the " + *moduleName,
			},
			{
				Name:        "code",
				NameSnake:   "code",
				Type:        "string",
				Required:    true,
				Unique:      true,
				Description: "Code of the " + *moduleName,
			},
			{
				Name:        "description",
				NameSnake:   "description",
				Type:        "string",
				Required:    false,
				Unique:      false,
				Description: "Description of the " + *moduleName,
			},
			{
				Name:        "isActive",
				NameSnake:   "is_active",
				Type:        "bool",
				Required:    true,
				Unique:      false,
				Description: "Whether the " + *moduleName + " is active",
			},
		},
	}

	// Generate module
	err := generateModule(config, *outputDir)
	if err != nil {
		fmt.Printf("Error generating module: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Module '%s' generated successfully in %s/%s\n", config.ModuleName, *outputDir, config.ModuleName)
}

// generateModule generates all files for a module
func generateModule(config ModuleConfig, outputDir string) error {
	// Create module directory
	moduleDir := filepath.Join(outputDir, config.ModuleName)
	if err := os.MkdirAll(moduleDir, 0755); err != nil {
		return fmt.Errorf("failed to create module directory: %w", err)
	}

	// Generate module.go
	if err := generateModuleFile(config, moduleDir); err != nil {
		return fmt.Errorf("failed to generate module file: %w", err)
	}

	// Generate controller files
	if err := generateControllerFiles(config, moduleDir); err != nil {
		return fmt.Errorf("failed to generate controller files: %w", err)
	}

	// Generate service files
	if err := generateServiceFiles(config, moduleDir); err != nil {
		return fmt.Errorf("failed to generate service files: %w", err)
	}

	// Generate repository files
	if err := generateRepositoryFiles(config, moduleDir); err != nil {
		return fmt.Errorf("failed to generate repository files: %w", err)
	}

	// Generate DTO files
	if err := generateDTOFiles(config, moduleDir); err != nil {
		return fmt.Errorf("failed to generate DTO files: %w", err)
	}

	return nil
}

// generateModuleFile generates the main module file
func generateModuleFile(config ModuleConfig, moduleDir string) error {
	tmpl := template.Must(template.New("module").Parse(`package {{.ModuleName}}

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/{{.ModuleName}}/dto/request"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
	"github.com/PhantomX7/dhamma/utility/repository"

	"github.com/gin-gonic/gin"
)

type Repository interface {
	repository.BaseRepositoryInterface[entity.{{.EntityName}}]
}

type Service interface {
	Index(ctx context.Context, pg *pagination.Pagination) ([]entity.{{.EntityName}}, utility.PaginationMeta, error)
	Show(ctx context.Context, {{.ModuleNameCamel}}ID uint64) (entity.{{.EntityName}}, error)
	Update(ctx context.Context, {{.ModuleNameCamel}}ID uint64, request request.{{.EntityName}}UpdateRequest) (entity.{{.EntityName}}, error)
	Create(ctx context.Context, request request.{{.EntityName}}CreateRequest) (entity.{{.EntityName}}, error)
}

type Controller interface {
	Index(ctx *gin.Context)
	Show(ctx *gin.Context)
	Update(ctx *gin.Context)
	Create(ctx *gin.Context)
}
`))

	return renderTemplate(tmpl, config, filepath.Join(moduleDir, fmt.Sprintf("%s.go", config.ModuleName)))
}

// generateControllerFiles generates all controller files
func generateControllerFiles(config ModuleConfig, moduleDir string) error {
	// Create controller directory
	controllerDir := filepath.Join(moduleDir, "controller")
	if err := os.MkdirAll(controllerDir, 0755); err != nil {
		return fmt.Errorf("failed to create controller directory: %w", err)
	}

	// Generate controller.go
	controllerTmpl := template.Must(template.New("controller").Parse(`package controller

import (
	"github.com/PhantomX7/dhamma/modules/{{.ModuleName}}"
)

type controller struct {
	{{.ModuleNameCamel}}Service {{.ModuleName}}.Service
}

func New({{.ModuleNameCamel}}Service {{.ModuleName}}.Service) {{.ModuleName}}.Controller {
	return &controller{
		{{.ModuleNameCamel}}Service: {{.ModuleNameCamel}}Service,
	}
}
`))

	if err := renderTemplate(controllerTmpl, config, filepath.Join(controllerDir, "controller.go")); err != nil {
		return fmt.Errorf("failed to generate controller.go: %w", err)
	}

	// Generate index.go
	indexTmpl := template.Must(template.New("index").Parse(`package controller

import (
	"net/http"

	"github.com/PhantomX7/dhamma/modules/{{.ModuleName}}/dto/request"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/gin-gonic/gin"
)

func (c *controller) Index(ctx *gin.Context) {
	res, meta, err := c.{{.ModuleNameCamel}}Service.Index(ctx.Request.Context(), request.New{{.EntityName}}Pagination(ctx.Request.URL.Query()))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildPaginationResponseSuccess("ok", res, meta))
}
`))

	if err := renderTemplate(indexTmpl, config, filepath.Join(controllerDir, "index.go")); err != nil {
		return fmt.Errorf("failed to generate index.go: %w", err)
	}

	// Generate show.go
	showTmpl := template.Must(template.New("show").Parse(`package controller

import (
	"net/http"
	"strconv"

	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/errors"
	"github.com/gin-gonic/gin"
)

func (c *controller) Show(ctx *gin.Context) {
	{{.ModuleNameCamel}}ID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Error(&errors.AppError{
			Message: "invalid {{.ModuleName}} id",
			Status:  http.StatusBadRequest,
		})
		return
	}

	res, err := c.{{.ModuleNameCamel}}Service.Show(ctx.Request.Context(), {{.ModuleNameCamel}}ID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", res))
}
`))

	if err := renderTemplate(showTmpl, config, filepath.Join(controllerDir, "show.go")); err != nil {
		return fmt.Errorf("failed to generate show.go: %w", err)
	}

	// Generate create.go
	createTmpl := template.Must(template.New("create").Parse(`package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/PhantomX7/dhamma/modules/{{.ModuleName}}/dto/request"
	"github.com/PhantomX7/dhamma/utility"
)

func (c *controller) Create(ctx *gin.Context) {
	var req request.{{.EntityName}}CreateRequest

	// validate request
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}

	res, err := c.{{.ModuleNameCamel}}Service.Create(ctx.Request.Context(), req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", res))
}
`))

	if err := renderTemplate(createTmpl, config, filepath.Join(controllerDir, "create.go")); err != nil {
		return fmt.Errorf("failed to generate create.go: %w", err)
	}

	// Generate update.go
	updateTmpl := template.Must(template.New("update").Parse(`package controller

import (
	"net/http"
	"strconv"

	"github.com/PhantomX7/dhamma/modules/{{.ModuleName}}/dto/request"
	"github.com/PhantomX7/dhamma/utility/errors"

	"github.com/gin-gonic/gin"

	"github.com/PhantomX7/dhamma/utility"
)

func (c *controller) Update(ctx *gin.Context) {
	var req request.{{.EntityName}}UpdateRequest

	// validate request
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}

	{{.ModuleNameCamel}}ID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Error(&errors.AppError{
			Message: "invalid {{.ModuleName}} id",
			Status:  http.StatusBadRequest,
		})
		return
	}

	res, err := c.{{.ModuleNameCamel}}Service.Update(ctx.Request.Context(), {{.ModuleNameCamel}}ID, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utility.BuildResponseSuccess("ok", res))
}
`))

	if err := renderTemplate(updateTmpl, config, filepath.Join(controllerDir, "update.go")); err != nil {
		return fmt.Errorf("failed to generate update.go: %w", err)
	}

	return nil
}

// generateServiceFiles generates all service files
func generateServiceFiles(config ModuleConfig, moduleDir string) error {
	// Create service directory
	serviceDir := filepath.Join(moduleDir, "service")
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		return fmt.Errorf("failed to create service directory: %w", err)
	}

	// Generate service.go
	serviceTmpl := template.Must(template.New("service").Parse(`package service

import (
	"github.com/PhantomX7/dhamma/modules/{{.ModuleName}}"
)

type service struct {
	{{.ModuleNameCamel}}Repo {{.ModuleName}}.Repository
}

func New(
	{{.ModuleNameCamel}}Repo {{.ModuleName}}.Repository,
) {{.ModuleName}}.Service {
	return &service{
		{{.ModuleNameCamel}}Repo: {{.ModuleNameCamel}}Repo,
	}
}
`))

	if err := renderTemplate(serviceTmpl, config, filepath.Join(serviceDir, "service.go")); err != nil {
		return fmt.Errorf("failed to generate service.go: %w", err)
	}

	// Generate index.go
	indexTmpl := template.Must(template.New("index").Parse(`package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/dhamma/utility/pagination"
)

// Index implements {{.ModuleName}}.Service.
func (s *service) Index(ctx context.Context, pg *pagination.Pagination) (
	{{.ModuleNameCamel}}s []entity.{{.EntityName}}, meta utility.PaginationMeta, err error,
) {
	{{.ModuleNameCamel}}s, err = s.{{.ModuleNameCamel}}Repo.FindAll(ctx, pg)
	if err != nil {
		return
	}

	count, err := s.{{.ModuleNameCamel}}Repo.Count(ctx, pg)
	if err != nil {
		return
	}

	meta.Limit = pg.Limit
	meta.Offset = pg.Offset
	meta.Total = count

	return
}
`))

	if err := renderTemplate(indexTmpl, config, filepath.Join(serviceDir, "index.go")); err != nil {
		return fmt.Errorf("failed to generate index.go: %w", err)
	}

	// Generate show.go
	showTmpl := template.Must(template.New("show").Parse(`package service

import (
	"context"

	"github.com/PhantomX7/dhamma/entity"
)

// Show implements {{.ModuleName}}.Service
func (s *service) Show(ctx context.Context, {{.ModuleNameCamel}}ID uint64) ({{.ModuleNameCamel}} entity.{{.EntityName}}, err error) {
	{{.ModuleNameCamel}}, err = s.{{.ModuleNameCamel}}Repo.FindByID(ctx, {{.ModuleNameCamel}}ID)
	if err != nil {
		return
	}

	return
}
`))

	if err := renderTemplate(showTmpl, config, filepath.Join(serviceDir, "show.go")); err != nil {
		return fmt.Errorf("failed to generate show.go: %w", err)
	}

	// Generate create.go
	createTmpl := template.Must(template.New("create").Parse(`package service

import (
	"context"
	"github.com/jinzhu/copier"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/{{.ModuleName}}/dto/request"
)

func (s *service) Create(ctx context.Context, request request.{{.EntityName}}CreateRequest) ({{.ModuleNameCamel}} entity.{{.EntityName}}, err error) {
	{{.ModuleNameCamel}} = entity.{{.EntityName}}{
		IsActive: true,
	}

	err = copier.Copy(&{{.ModuleNameCamel}}, &request)
	if err != nil {
		return
	}

	err = s.{{.ModuleNameCamel}}Repo.Create(ctx, &{{.ModuleNameCamel}}, nil)
	if err != nil {
		return
	}

	return
}
`))

	if err := renderTemplate(createTmpl, config, filepath.Join(serviceDir, "create.go")); err != nil {
		return fmt.Errorf("failed to generate create.go: %w", err)
	}

	// Generate update.go
	updateTmpl := template.Must(template.New("update").Parse(`package service

import (
	"context"
	"github.com/jinzhu/copier"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/{{.ModuleName}}/dto/request"
)

func (s *service) Update(ctx context.Context, {{.ModuleNameCamel}}ID uint64, request request.{{.EntityName}}UpdateRequest) ({{.ModuleNameCamel}} entity.{{.EntityName}}, err error) {
	{{.ModuleNameCamel}}, err = s.{{.ModuleNameCamel}}Repo.FindByID(ctx, {{.ModuleNameCamel}}ID)
	if err != nil {
		return
	}

	err = copier.Copy(&{{.ModuleNameCamel}}, &request)
	if err != nil {
		return
	}

	err = s.{{.ModuleNameCamel}}Repo.Update(ctx, &{{.ModuleNameCamel}}, nil)
	if err != nil {
		return
	}

	return
}
`))

	if err := renderTemplate(updateTmpl, config, filepath.Join(serviceDir, "update.go")); err != nil {
		return fmt.Errorf("failed to generate update.go: %w", err)
	}

	return nil
}

// generateRepositoryFiles generates all repository files
func generateRepositoryFiles(config ModuleConfig, moduleDir string) error {
	// Create repository directory
	repoDir := filepath.Join(moduleDir, "repository")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		return fmt.Errorf("failed to create repository directory: %w", err)
	}

	// Generate repository.go
	repoTmpl := template.Must(template.New("repository").Parse(`package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/{{.ModuleName}}"
	"github.com/PhantomX7/dhamma/utility/pagination"
	baseRepo "github.com/PhantomX7/dhamma/utility/repository"
)

type repository struct {
	base baseRepo.BaseRepositoryInterface[entity.{{.EntityName}}] // Use the interface type
	db   *gorm.DB
}

// New creates a new {{.ModuleName}} repository instance.
func New(db *gorm.DB) {{.ModuleName}}.Repository {
	return &repository{
		base: baseRepo.NewBaseRepository[entity.{{.EntityName}}](db), // Instantiate the concrete base repository
		db:   db,
	}
}

// FindAll retrieves all {{.ModuleName}} entities with pagination.
func (r *repository) FindAll(ctx context.Context, pg *pagination.Pagination) ([]entity.{{.EntityName}}, error) {
	return r.base.FindAll(ctx, pg)
}

// FindByID retrieves a {{.ModuleName}} entity by its ID.
func (r *repository) FindByID(ctx context.Context, {{.ModuleNameCamel}}ID uint64, preloads ...string) (entity.{{.EntityName}}, error) {
	return r.base.FindByID(ctx, {{.ModuleNameCamel}}ID, preloads...)
}

// Create creates a new {{.ModuleName}} entity.
func (r *repository) Create(ctx context.Context, {{.ModuleNameCamel}} *entity.{{.EntityName}}, tx *gorm.DB) error {
	return r.base.Create(ctx, {{.ModuleNameCamel}}, tx)
}

// Update updates an existing {{.ModuleName}} entity.
func (r *repository) Update(ctx context.Context, {{.ModuleNameCamel}} *entity.{{.EntityName}}, tx *gorm.DB) error {
	return r.base.Update(ctx, {{.ModuleNameCamel}}, tx)
}

// Delete deletes a {{.ModuleName}} entity.
func (r *repository) Delete(ctx context.Context, {{.ModuleNameCamel}} *entity.{{.EntityName}}, tx *gorm.DB) error {
	return r.base.Delete(ctx, {{.ModuleNameCamel}}, tx)
}

// Count counts {{.ModuleName}} entities matching pagination filters.
func (r *repository) Count(ctx context.Context, pg *pagination.Pagination) (int64, error) {
	return r.base.Count(ctx, pg)
}

// FindByField retrieves {{.ModuleName}} entities where a specific field matches the given value.
func (r *repository) FindByField(ctx context.Context, fieldName string, value any, preloads ...string) ([]entity.{{.EntityName}}, error) {
	return r.base.FindByField(ctx, fieldName, value, preloads...)
}

// FindOneByField retrieves a single {{.ModuleName}} entity where a specific field matches the given value.
func (r *repository) FindOneByField(ctx context.Context, fieldName string, value any, preloads ...string) (entity.{{.EntityName}}, error) {
	return r.base.FindOneByField(ctx, fieldName, value, preloads...)
}

// FindByFields retrieves {{.ModuleName}} entities matching multiple field conditions.
func (r *repository) FindByFields(ctx context.Context, conditions map[string]any, preloads ...string) ([]entity.{{.EntityName}}, error) {
	return r.base.FindByFields(ctx, conditions, preloads...)
}

// FindOneByFields retrieves a single {{.ModuleName}} entity matching multiple field conditions.
func (r *repository) FindOneByFields(ctx context.Context, conditions map[string]any, preloads ...string) (entity.{{.EntityName}}, error) {
	return r.base.FindOneByFields(ctx, conditions, preloads...)
}

// Exists checks if any {{.ModuleName}} records match the given conditions.
func (r *repository) Exists(ctx context.Context, conditions map[string]any) (bool, error) {
	return r.base.Exists(ctx, conditions)
}
`))

	if err := renderTemplate(repoTmpl, config, filepath.Join(repoDir, "repository.go")); err != nil {
		return fmt.Errorf("failed to generate repository.go: %w", err)
	}

	return nil
}

// generateDTOFiles generates all DTO files
func generateDTOFiles(config ModuleConfig, moduleDir string) error {
	// Create DTO directories
	dtoDir := filepath.Join(moduleDir, "dto")
	if err := os.MkdirAll(dtoDir, 0755); err != nil {
		return fmt.Errorf("failed to create dto directory: %w", err)
	}

	requestDir := filepath.Join(dtoDir, "request")
	if err := os.MkdirAll(requestDir, 0755); err != nil {
		return fmt.Errorf("failed to create request directory: %w", err)
	}

	// Generate request.go
	var createFields, updateFields, filterFields, sortFields strings.Builder

	// Build create fields
	for _, field := range config.Fields {
		fieldType := field.Type
		if field.Required {
			createFields.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\" form:\"%s\" binding:\"required",
				snakeToPascal(field.Name), fieldType, field.NameSnake, field.NameSnake))
			if field.Unique {
				createFields.WriteString(fmt.Sprintf(",unique=%s.%s", config.TableName, field.NameSnake))
			}
			createFields.WriteString("\"`\n")
		} else {
			createFields.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\" form:\"%s\"`\n",
				snakeToPascal(field.Name), fieldType, field.NameSnake, field.NameSnake))
		}
	}

	// Build update fields
	for _, field := range config.Fields {
		fieldType := "*" + field.Type
		updateFields.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\" form:\"%s\" binding:\"omitempty",
			snakeToPascal(field.Name), fieldType, field.NameSnake, field.NameSnake))
		if field.Unique {
			updateFields.WriteString(fmt.Sprintf(",unique=%s.%s", config.TableName, field.NameSnake))
		}
		updateFields.WriteString("\"`\n")
	}

	// Build filter fields
	for _, field := range config.Fields {
		filterType := "pagination.FilterTypeString"
		operators := "pagination.OperatorIn, pagination.OperatorEquals, pagination.OperatorLike"

		if field.Type == "uint64" {
			filterType = "pagination.FilterTypeID"
			operators = "pagination.OperatorIn, pagination.OperatorEquals"
		} else if field.Type == "bool" {
			filterType = "pagination.FilterTypeBool"
			operators = "pagination.OperatorEquals"
		}

		filterFields.WriteString(fmt.Sprintf(`		AddFilter("%s", pagination.FilterConfig{
			Field: "%s",
			Type:  %s,
			Operators: []pagination.FilterOperator{
				%s,
			},
		}).
`, field.NameSnake, field.NameSnake, filterType, operators))
	}

	// Add created_at filter
	filterFields.WriteString(`		AddFilter("created_at", pagination.FilterConfig{
			Field:     "created_at",
			Type:      pagination.FilterTypeDateTime,
			Operators: []pagination.FilterOperator{pagination.OperatorBetween},
		}).
`)

	// Build sort fields
	for _, field := range config.Fields {
		sortFields.WriteString(fmt.Sprintf(`		AddSort("%s", pagination.SortConfig{
			Field:   "%s",
			Allowed: true,
		}).
`, field.NameSnake, field.NameSnake))
	}

	requestTmpl := template.Must(template.New("request").Parse(`package request

import "github.com/PhantomX7/dhamma/utility/pagination"

type {{.EntityName}}CreateRequest struct {
{{.CreateFields}}
}

type {{.EntityName}}UpdateRequest struct {
{{.UpdateFields}}
}

func New{{.EntityName}}Pagination(conditions map[string][]string) *pagination.Pagination {
	filterDef := pagination.NewFilterDefinition().
{{.FilterFields}}

{{.SortFields}}

	return pagination.NewPagination(
		conditions,
		filterDef,
		pagination.PaginationOptions{
			DefaultLimit: 20,
			MaxLimit:     100,
			DefaultOrder: "id desc",
		},
	)
}
`))

	data := struct {
		ModuleConfig
		CreateFields string
		UpdateFields string
		FilterFields string
		SortFields   string
	}{
		ModuleConfig: config,
		CreateFields: createFields.String(),
		UpdateFields: updateFields.String(),
		FilterFields: filterFields.String(),
		SortFields:   sortFields.String(),
	}

	if err := renderTemplateWithData(requestTmpl, data, filepath.Join(requestDir, "request.go")); err != nil {
		return fmt.Errorf("failed to generate request.go: %w", err)
	}

	return nil
}

// renderTemplate renders a template with the given config and writes it to the specified file
func renderTemplate(tmpl *template.Template, config ModuleConfig, filePath string) error {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

// renderTemplateWithData renders a template with the given data and writes it to the specified file
func renderTemplateWithData(tmpl *template.Template, data interface{}, filePath string) error {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

// snakeToCamel converts a snake_case string to camelCase
func snakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}

// snakeToPascal converts a snake_case string to PascalCase
func snakeToPascal(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}
