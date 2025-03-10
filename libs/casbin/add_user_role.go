package casbin

import "fmt"

// AddUserRole assigns a role to a user within a specific domain
func (c *client) AddUserRole(userID uint64, roleID uint64, domainID uint64) error {
	_, err := c.enforcer.AddGroupingPolicy(
		fmt.Sprintf("%d", userID),
		fmt.Sprintf("%d", roleID),
		fmt.Sprintf("%d", domainID),
	)
	return err
}
