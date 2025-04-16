package permissions

import (
	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/auth"
	"github.com/PhantomX7/dhamma/modules/permission"
	"github.com/PhantomX7/dhamma/modules/role"
	"github.com/PhantomX7/dhamma/modules/user"
)

const PermissionTypeApi = "API"
const PermissionTypeWeb = "WEB"

// Define permissions using module keys and action constants
var ApiPermissions = []entity.Permission{
	// Auth Module
	{
		Name:             "auth - update password",
		Object:           auth.Permissions.Key, // Use module key
		Action:           ActionUpdatePassword, // Use action constant
		Description:      "update password route",
		Type:             PermissionTypeApi,
		IsDomainSpecific: false, // Auth actions are typically not domain-specific
	},

	// Role Module
	{
		Name:             "role - list",
		Object:           role.Permissions.Key,
		Action:           ActionList, // Use consistent action name (e.g., list or index)
		Description:      "role list route",
		Type:             PermissionTypeApi,
		IsDomainSpecific: false,
	},
	{
		Name:             "role - show",
		Object:           role.Permissions.Key,
		Action:           ActionShow,
		Description:      "role show route",
		Type:             PermissionTypeApi,
		IsDomainSpecific: false,
	},
	{
		Name:             "role - create",
		Object:           role.Permissions.Key,
		Action:           ActionCreate,
		Description:      "role create route",
		Type:             PermissionTypeApi,
		IsDomainSpecific: false,
	},
	{
		Name:             "role - update",
		Object:           role.Permissions.Key,
		Action:           ActionUpdate,
		Description:      "role update route",
		Type:             PermissionTypeApi,
		IsDomainSpecific: false,
	},
	{
		Name:             "role - add permission",
		Object:           role.Permissions.Key,
		Action:           ActionAddPermissions,
		Description:      "role add permission route",
		Type:             PermissionTypeApi,
		IsDomainSpecific: false,
	},
	// Add role - delete if needed
	// {
	// 	Name:             "role - delete",
	// 	Object:           role.Permissions.Key,
	// 	Action:           ActionDelete,
	// 	Description:      "role delete route",
	// 	Type:             PermissionTypeApi,
	// 	IsDomainSpecific: true,
	// },

	// User Module
	{
		Name:             "user - list",
		Object:           user.Permissions.Key,
		Action:           ActionList, // Use consistent action name
		Description:      "user list route",
		Type:             PermissionTypeApi,
		IsDomainSpecific: false, // Users might be listed per domain or globally, adjust as needed
	},
	{
		Name:             "user - show",
		Object:           user.Permissions.Key,
		Action:           ActionShow,
		Description:      "user show route",
		Type:             PermissionTypeApi,
		IsDomainSpecific: false, // Showing a user might depend on domain context
	},
	{
		Name:             "user - create",
		Object:           user.Permissions.Key,
		Action:           ActionCreate,
		Description:      "user create route",
		Type:             PermissionTypeApi,
		IsDomainSpecific: false, // Creating a user might be a global action
	},
	// Add user - update/delete/assign-role/assign-domain if needed
	// {
	// 	Name:             "user - update",
	// 	Object:           user.Permissions.Key,
	// 	Action:           ActionUpdate,
	// 	Description:      "user update route",
	// 	Type:             PermissionTypeApi,
	// 	IsDomainSpecific: true,
	// },

	// Permission Module
	{
		Name:             "permission - list",
		Object:           permission.Permissions.Key,
		Action:           ActionList, // Use consistent action name
		Description:      "permission list route",
		Type:             PermissionTypeApi,
		IsDomainSpecific: false, // Permissions are likely global
	},
	// Add permission - show if needed
	// {
	// 	Name:             "permission - show",
	// 	Object:           permission.Permissions.Key,
	// 	Action:           ActionShow,
	// 	Description:      "permission show route",
	// 	Type:             PermissionTypeApi,
	// 	IsDomainSpecific: false,
	// },
}

// Helper function to get all permission codes defined above
func GetAllPermissionCodes() []string {
	codes := make([]string, len(ApiPermissions))
	for i, p := range ApiPermissions {
		codes[i] = p.Object + "/" + p.Action
	}
	return codes
}
