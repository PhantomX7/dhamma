package example

// permission defines all permissions for the example module
type permission struct {
	// Key is the module identifier used as prefix for all permissions
	Key string

	// explanatory (dont include in code) // User management permissions

	// List all examples
	// @group:example-management
	Index string

	// View example details
	// @group:example-management
	Show string

	// Create a new example
	// @group:example-management
	Create string

	// Update example information
	// @group:example-management
	Update string

	// Delete an example
	// @group:example-management
	Delete string

	// explanatory (dont include in code) // Example-specific operations

	// Process an example
	// @group:example-operations
	Process string

	// Export example data
	// @group:example-operations
	Export string

	// Import example data
	// @group:example-operations
	Import string
}

// Permissions defines all permissions for the example module
var Permissions = permission{
	Key:     "example",
	Index:   "list",
	Show:    "show",
	Create:  "create",
	Update:  "update",
	Delete:  "delete",
	Process: "process",
	Export:  "export",
	Import:  "import",
}

// Service defines the example service interface
type Service interface {
	// Add your service methods here
}

// Repository defines the example repository interface
type Repository interface {
	// Add your repository methods here
}

// Controller defines the example controller interface
type Controller interface {
	// Add your controller methods here
}
