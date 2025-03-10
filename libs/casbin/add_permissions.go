package casbin

import (
	"fmt"
	"strings"

	"github.com/PhantomX7/dhamma/constants"
)

func (c *client) AddPermissions(roleID uint64, domainID uint64, permissionsCodes []string) {
	for _, permissionCode := range permissionsCodes {
		codes := strings.Split(permissionCode, ":")
		permissionType := constants.EnumPermissionTypeApi
		var permissionObject, permissionAction string
		if len(codes) > 1 {
			permissionType = constants.EnumPermissionTypeWeb
			permission := strings.Split(codes[1], "/")
			permissionObject = permission[0]
			permissionAction = permission[1]
		} else {
			permission := strings.Split(codes[0], "/")
			permissionObject = permission[0]
			permissionAction = permission[1]
		}

		_, err := c.enforcer.AddPolicy(fmt.Sprintf("%d", roleID), fmt.Sprintf("%d", domainID), permissionObject, permissionAction, permissionType)
		if err != nil {
			return
		}

	}
}
