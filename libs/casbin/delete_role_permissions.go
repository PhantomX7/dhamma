package casbin

import (
	"errors" // Import errors package
	"fmt"

	// "github.com/PhantomX7/dhamma/constants/permissions" // No longer needed directly here if parsePermissionCode handles it
	"github.com/PhantomX7/dhamma/utility/logger"
	"go.uber.org/zap"
)

// DeleteRolePermissions removes specified permissions for a role within a domain.
// It collects parsing errors and returns them if any occur.
// It attempts to remove all validly parsed permissions in a single batch operation.
func (c *client) DeleteRolePermissions(roleID uint64, domainID uint64, permissionCodes []string) error {
	roleStr := fmt.Sprintf("%d", roleID)
	domainStr := fmt.Sprintf("%d", domainID)
	var encounteredErrors []error // Slice to collect errors
	rules := make([][]string, 0, len(permissionCodes))

	// Parse all permission codes first and collect valid rules and parsing errors
	for _, code := range permissionCodes {
		object, action, permissionType, parseErr := parsePermissionCode(code)
		if parseErr != nil {
			// Log the parsing error and collect it
			logger.Get().Warn("Skipping invalid permission code format during deletion", zap.String("code", code), zap.Error(parseErr))
			encounteredErrors = append(encounteredErrors, parseErr) // Collect the error
			continue                                                // Skip to the next permission code
		}
		// Rule format: [subject, domain, object, action, type]
		rules = append(rules, []string{roleStr, domainStr, object, action, permissionType})
	}

	// If any parsing errors occurred, return them immediately before attempting removal
	if len(encounteredErrors) > 0 {
		// Use errors.Join for Go 1.20+
		return errors.Join(encounteredErrors...)
		// For older Go:
		// return fmt.Errorf("encountered %d error(s) parsing permission codes, first error: %w", len(encounteredErrors), encounteredErrors[0])
	}

	// If no valid rules were generated (e.g., all input codes were invalid), return an error.
	if len(rules) == 0 {
		// Check if the original input was also empty
		if len(permissionCodes) == 0 {
			logger.Get().Info("No permissions provided to delete.", zap.Uint64("roleID", roleID), zap.Uint64("domainID", domainID))
			return nil // Nothing to do
		}
		// Input was provided, but all were invalid
		return fmt.Errorf("no valid permissions provided to delete after parsing")
	}

	// Attempt to remove the valid policies in a batch
	_, err := c.enforcer.RemovePolicies(rules)
	if err != nil {
		logger.Get().Error("Failed to remove policies from Casbin",
			zap.Uint64("roleID", roleID),
			zap.Uint64("domainID", domainID),
			zap.Error(err),
		)
		// Wrap the Casbin error for context
		return fmt.Errorf("casbin policy removal failed: %w", err)
	}

	// Persist changes if auto-save is disabled or not guaranteed by the adapter
	err = c.enforcer.SavePolicy()
	if err != nil {
		logger.Get().Error("Failed to save policy after deletion",
			zap.Uint64("roleID", roleID),
			zap.Uint64("domainID", domainID),
			zap.Error(err),
		)
		// Wrap the save error
		return fmt.Errorf("casbin policy save failed: %w", err)
	}

	// Log successful processing
	logger.Get().Info("Successfully processed policy removal request",
		zap.Uint64("roleID", roleID),
		zap.Uint64("domainID", domainID),
		zap.Int("valid_rules_count", len(rules)), // Log how many were attempted
		zap.Strings("original_permissions_requested", permissionCodes),
	)

	return nil // Success
}
