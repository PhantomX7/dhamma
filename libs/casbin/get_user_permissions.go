package casbin

import (
	"fmt"

	"github.com/PhantomX7/dhamma/constants/permissions"
)

// GetUserPermissions returns all permissions a user has in a specific domain through their roles
func (c *client) GetUserPermissions(userID uint64, domainID uint64) []string {
	// Get all policies that apply to this user in this domain
	policies, err := c.enforcer.GetImplicitPermissionsForUser(fmt.Sprintf("%d", userID), fmt.Sprintf("%d", domainID))
	if err != nil {
		return []string{}
	}

	permissionList := make([]string, 0)
	for _, policy := range policies {
		// policy format: [userID domainID object action type]
		if len(policy) >= 5 {
			permissionType := policy[4]
			if permissionType == permissions.PermissionTypeWeb {
				permissionList = append(permissionList, "web:"+policy[2]+"/"+policy[3])
			} else {
				permissionList = append(permissionList, policy[2]+"/"+policy[3])
			}
		}
	}

	return permissionList
}
