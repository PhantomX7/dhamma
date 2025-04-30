package role

type permission struct {
	Key string
	// Index all roles
	Index string
	// View role details
	Show string
	// Create a new role
	Create string
	// Update role information
	Update string
	// Add permissions to a role
	AddPermissions string
	// Delete permissions from a role
	DeletePermissions string
}

var Permissions = permission{
	Key:               "role",
	Index:             "index",
	Show:              "show",
	Create:            "create",
	Update:            "update",
	AddPermissions:    "add-permissions",
	DeletePermissions: "delete-permissions",
}
