package service_test

import (
	"context"
	"github.com/PhantomX7/dhamma/config"
	"github.com/PhantomX7/dhamma/constants"
	"github.com/PhantomX7/dhamma/entity"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// All methods that begin with "Test" are run as tests within a
// suite.
//
//	func (suite *ExampleTestSuite) TestExample() {
//		assert.Equal(suite.T(), 5, suite.VariableThatShouldStartAtFive)
//		suite.Equal(5, suite.VariableThatShouldStartAtFive)
//	}
func (suite *AuthServiceSuite) TestGenerateRefreshToken() {
	// Test cases
	tests := []struct {
		name         string
		contextSetup func() context.Context
		paramSetup   func() (userID uint64, tx *gorm.DB)
		mockSetup    func()
		//expectedResult string
		expectedError error
		validateToken func(t *testing.T, tokenString string, userID uint64)
	}{
		{
			name: "Success - create refresh token without transaction",
			paramSetup: func() (userID uint64, tx *gorm.DB) {
				return 1, nil
			},
			mockSetup: func() {
				suite.mockRefreshTokenRepo.On("Create",
					mock.Anything,
					mock.MatchedBy(func(rt *entity.RefreshToken) bool {
						return rt.UserID == uint64(1) &&
							rt.IsValid &&
							!rt.ExpiresAt.IsZero() &&
							rt.ID != uuid.Nil
					}),
					mock.Anything,
				).Return(nil)
			},
			validateToken: func(t *testing.T, tokenString string, userID uint64) {
				// Parse and validate token
				token, err := jwt.ParseWithClaims(tokenString, &entity.RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(config.JWT_SECRET), nil
				})

				assert.NoError(t, err)
				assert.True(t, token.Valid)

				// Check claims
				if claims, ok := token.Claims.(*entity.RefreshClaims); ok {
					// Verify refresh token UUID is valid
					_, err := uuid.Parse(claims.RefreshToken)
					assert.NoError(t, err)

					// Verify expiry
					assert.NotNil(t, claims.ExpiresAt)
					expectedExpiry := time.Now().Add(constants.RefreshTokenExpiry)
					timeDiff := claims.ExpiresAt.Sub(expectedExpiry)
					assert.Less(t, timeDiff.Abs(), time.Second) // Allow 1 second difference
				} else {
					t.Error("invalid token claims type")
				}
			},
			expectedError: nil,
		},
	}

	// Run tests
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			suite.New()
			// Setup mocks
			tt.mockSetup()

			userID, dbTransaction := tt.paramSetup()
			// Execute
			tokenString, err := suite.service.GenerateRefreshToken(userID, dbTransaction)

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
			tt.validateToken(t, tokenString, userID)
		})
	}
}
