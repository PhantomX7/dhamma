package casbin

import (
	// Import errors package
	"errors"
	"fmt"

	"github.com/PhantomX7/dhamma/utility/logger"
	"go.uber.org/zap" // Assuming you use zap logger as elsewhere
)

// AddRolePermissions adds multiple permissions for a given role and domain to Casbin.
// It returns an error if any permission fails to be added.
func (c *client) AddRolePermissions(roleID uint64, domainID uint64, permissionsCodes []string) error {
	var encounteredErrors []error // Slice to collect errors

	roleStr := fmt.Sprintf("%d", roleID)
	domainStr := fmt.Sprintf("%d", domainID)

	for _, permissionCode := range permissionsCodes {
		// Use the helper function to parse the code
		permissionObject, permissionAction, permissionType, parseErr := parsePermissionCode(permissionCode)

		if parseErr != nil {
			// Log and collect parsing error, then continue to next code
			// Use c.logger instead of logger.Logger if 'c' has the logger instance
			logger.Logger.Warn("Skipping invalid permission code", zap.String("code", permissionCode), zap.Error(parseErr))
			encounteredErrors = append(encounteredErrors, parseErr)
			continue // Skip to the next permission code
		}

		// Add the policy
		added, err := c.enforcer.AddPolicy(roleStr, domainStr, permissionObject, permissionAction, permissionType)
		if err != nil {
			// Log and collect error, then continue
			logger.Logger.Error(
				"Failed to add policy to Casbin",
				zap.Uint64("roleID", roleID),
				zap.Uint64("domainID", domainID),
				zap.String("object", permissionObject),
				zap.String("action", permissionAction),
				zap.String("type", permissionType),
				zap.Error(err),
			)
			encounteredErrors = append(encounteredErrors, fmt.Errorf("failed to add permission '%s': %w", permissionCode, err))
			continue // Continue processing other permissions
		}
		if !added {
			// Policy already existed, which might be okay or might indicate an issue depending on requirements.
			// Log it for information.
			logger.Logger.Info(
				"Policy already exists in Casbin",
				zap.Uint64("roleID", roleID),
				zap.Uint64("domainID", domainID),
				zap.String("object", permissionObject),
				zap.String("action", permissionAction),
				zap.String("type", permissionType),
			)
			// Optionally, add a specific error if duplicates are not allowed:
			// encounteredErrors = append(encounteredErrors, fmt.Errorf("permission '%s' already exists", permissionCode))
		}
	}

	// Return a combined error if any occurred (requires Go 1.20+)
	// If using older Go, you might return the first error or a custom error type.
	if len(encounteredErrors) > 0 {
		// Use errors.Join for Go 1.20+
		return errors.Join(encounteredErrors...)

		// For older Go versions, return a generic error or the first one:
		// return fmt.Errorf("encountered %d error(s) while adding permissions, first error: %w", len(encounteredErrors), encounteredErrors[0])
	}

	// Persist changes if auto-save is disabled
	err := c.enforcer.SavePolicy()
	if err != nil {
		logger.Logger.Error("Failed to save policy after adding permissions", zap.Error(err))
		return fmt.Errorf("casbin policy save failed: %w", err)
	}

	return nil // Success
}

// Helper function (can be shared or kept private)
