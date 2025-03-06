package constants

import "time"

const (
	EnumRoleRoot  = "root"
	EnumRoleAdmin = "admin"

	EnumRunProduction = "production"
	ENUM_RUN_TESTING  = "testing"

	EnumRequestIDKey = "request_id"
)

const AccessTokenExpiry = 30 * time.Minute
const RefreshTokenExpiry = 24 * time.Hour
