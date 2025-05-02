package user

// Define permissions for this module
type permission struct {
	Key string
	// List all users
	Index string
	// View user details
	Show string
	// Create a new user
	Create string
	// Update user information
	Update string
	// Assign a role to a user
	AssignRole string
	// Remove a role from a user
	RemoveRole string
	// Force logout a user
	ForceLogout string
}

// Permissions defines all permissions for the user module
var Permissions = permission{
	Key:         "user",
	Index:       "index",
	Show:        "show",
	Create:      "create",
	Update:      "update",
	AssignRole:  "assign-role",
	RemoveRole:  "remove-role",
	ForceLogout: "force-logout",
}
