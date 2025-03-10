package casbin

import (
	"fmt"

	"github.com/PhantomX7/dhamma/constants"
)

func (c *client) GetRolePermissions(roleID uint64, domainID uint64) []string {
	policies, err := c.enforcer.GetFilteredPolicy(0, fmt.Sprintf("%d", roleID), fmt.Sprintf("%d", domainID))
	if err != nil {
		return []string{}
	}
	permissions := make([]string, 0)

	for _, policy := range policies {
		// policy format: [roleID domainID object action type]
		if len(policy) >= 5 {
			permissionType := policy[4]
			if permissionType == constants.EnumPermissionTypeWeb {
				permissions = append(permissions, "web:"+policy[2]+"/"+policy[3])
			} else {
				permissions = append(permissions, policy[2]+"/"+policy[3])
			}
		}
	}

	return permissions
}
