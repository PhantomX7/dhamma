package constants

import (
	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/auth"
	"github.com/PhantomX7/dhamma/modules/role"
	"github.com/PhantomX7/dhamma/modules/user"
)

const EnumPermissionTypeApi = "API"
const EnumPermissionTypeWeb = "WEB"

var ApiPermissions = []entity.Permission{
	{
		Name:        "get me",
		Object:      auth.Permissions.Key,
		Action:      auth.Permissions.GetMe,
		Description: "get me route",
		Type:        EnumPermissionTypeApi,
	},
	{
		Name:        "update password",
		Object:      auth.Permissions.Key,
		Action:      auth.Permissions.UpdatePassword,
		Description: "update password route",
		Type:        EnumPermissionTypeApi,
	},
	{
		Name:        "role - index",
		Object:      role.Permissions.Key,
		Action:      role.Permissions.Index,
		Description: "role index route",
		Type:        EnumPermissionTypeApi,
	},
	{
		Name:        "role - show",
		Object:      role.Permissions.Key,
		Action:      role.Permissions.Show,
		Description: "role show route",
		Type:        EnumPermissionTypeApi,
	},
	{
		Name:        "role - create",
		Object:      role.Permissions.Key,
		Action:      role.Permissions.Create,
		Description: "role create route",
		Type:        EnumPermissionTypeApi,
	},
	{
		Name:        "role - update",
		Object:      role.Permissions.Key,
		Action:      role.Permissions.Update,
		Description: "role update route",
		Type:        EnumPermissionTypeApi,
	},
	{
		Name:        "user - index",
		Object:      user.Permissions.Key,
		Action:      user.Permissions.Index,
		Description: "user index route",
		Type:        EnumPermissionTypeApi,
	},
	{
		Name:        "user - show",
		Object:      user.Permissions.Key,
		Action:      user.Permissions.Show,
		Description: "user show route",
		Type:        EnumPermissionTypeApi,
	},
	{
		Name:        "user - create",
		Object:      user.Permissions.Key,
		Action:      user.Permissions.Create,
		Description: "user create route",
		Type:        EnumPermissionTypeApi,
	},
}
