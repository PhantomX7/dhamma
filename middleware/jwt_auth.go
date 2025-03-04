package middleware

import (
	"errors"
	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/utility"
	"net/http"
	"strings"

	"github.com/PhantomX7/dhamma/config"
	"github.com/PhantomX7/dhamma/constants"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func (m *Middleware) AuthHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate the token
		claims := &entity.AccessClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(config.JWT_SECRET), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if refresh token is valid
		var count int64
		count, _ = m.refreshTokenRepo.GetValidCountByUserID(c.Request.Context(), claims.UserID)
		if count == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			c.Abort()
			return
		}

		// Set user details in context

		c.Request = c.Request.WithContext(utility.NewContextWithValues(
			c.Request.Context(),
			utility.ContextValues{
				UserID: claims.UserID,
				IsRoot: claims.Role == constants.EnumRoleRoot,
			},
		))

		c.Next()
	}
}
