package utility

import (
	"fmt"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// GetRoleFromContext this is not a heavyweight operation
// it accesses payload from map in gin context
func GetRoleFromContext(c *gin.Context) (role string) {
	claims := jwt.ExtractClaims(c)
	role = fmt.Sprint(claims["role"])
	return
}

// GetIDFromContext this is not a heavyweight operation
// it accesses payload from map in gin context
func GetIDFromContext(c *gin.Context) (userID uint64) {
	claims := jwt.ExtractClaims(c)
	userID, _ = strconv.ParseUint(fmt.Sprint(claims["id"]), 10, 64)
	return
}

// GetUsernameFromContext this is not a heavyweight operation
// it accesses payload from map in gin context
func GetUsernameFromContext(c *gin.Context) (username string) {
	claims := jwt.ExtractClaims(c)
	username = fmt.Sprint(claims["username"])
	return
}

func GetKeyFromContext(c *gin.Context) (key string) {
	claims := jwt.ExtractClaims(c)
	key = fmt.Sprint(claims["key"])
	return
}

// GetDeviceTokenFromContext get device token from header
func GetDeviceTokenFromContext(c *gin.Context) string {
	return c.GetHeader("device_token")
}
