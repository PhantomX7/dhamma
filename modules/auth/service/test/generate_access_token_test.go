package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/PhantomX7/dhamma/config"
	"github.com/PhantomX7/dhamma/constants"
	"github.com/PhantomX7/dhamma/entity"
	"github.com/golang-jwt/jwt/v4"

	"github.com/stretchr/testify/assert"
)

func (suite *AuthServiceSuite) TestGenerateAccessToken() {
	// Test cases
	tests := []struct {
		name         string
		contextSetup func() context.Context
		paramSetup   func() (userID uint64, role string)
		mockSetup    func()
		//expectedResult string
		expectedError error
		validateToken func(t *testing.T, tokenString string, userID uint64, role string)
	}{
		{
			name: "Success - create token with user id and admin role",
			paramSetup: func() (userID uint64, role string) {
				return 1, constants.EnumRoleAdmin
			},
			validateToken: func(t *testing.T, tokenString string, userID uint64, role string) {
				// Parse and validate token
				token, err := jwt.ParseWithClaims(tokenString, &entity.AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(config.JWT_SECRET), nil
				})

				assert.NoError(t, err)
				assert.True(t, token.Valid)

				// Check claims
				if claims, ok := token.Claims.(*entity.AccessClaims); ok {
					assert.Equal(t, userID, claims.UserID)
					assert.Equal(t, role, claims.Role)

					// Validate time claims
					now := time.Now()
					assert.True(t, claims.ExpiresAt.After(now))
					assert.True(t, claims.IssuedAt.Before(now) || claims.IssuedAt.Equal(now))
					assert.True(t, claims.NotBefore.Before(now) || claims.NotBefore.Equal(now))

					// Validate expiry duration
					expectedExpiry := now.Add(constants.AccessTokenExpiry)
					timeDiff := claims.ExpiresAt.Sub(expectedExpiry)
					assert.Less(t, timeDiff.Abs(), time.Second) // Allow 1 second difference due to execution time
				} else {
					t.Error("invalid token claims type")
				}
			},
			expectedError: nil,
		},
		{
			name: "Error - no role provided",
			paramSetup: func() (userID uint64, role string) {
				return 1, ""
			},
			expectedError: errors.New("empty role"),
		},
	}

	// Run tests
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			suite.New()
			// Setup mocks
			//tt.mockSetup()

			userID, role := tt.paramSetup()
			// Execute
			tokenString, err := suite.service.GenerateAccessToken(userID, role)

			// Check error
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Empty(t, tokenString)
				return
			}

			// Verify success case
			assert.NoError(t, err)
			assert.NotEmpty(t, tokenString)

			// Validate token
			tt.validateToken(t, tokenString, userID, role)
		})
	}
}
