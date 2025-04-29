package casbin

import "fmt"

// RemoveUserRole removes a role from a user within a specific domain
func (c *client) RemoveUserRole(userID uint64, roleID uint64, domainID uint64) error {
	_, err := c.enforcer.RemoveGroupingPolicy(
		fmt.Sprintf("%d", userID),
		fmt.Sprintf("%d", roleID),
		fmt.Sprintf("%d", domainID),
	)
	return err
}
