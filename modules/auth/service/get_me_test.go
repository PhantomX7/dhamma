package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/auth/dto/response"
	"github.com/PhantomX7/dhamma/utility"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// All methods that begin with "Test" are run as tests within a
// suite.
//
//	func (suite *ExampleTestSuite) TestExample() {
//		assert.Equal(suite.T(), 5, suite.VariableThatShouldStartAtFive)
//		suite.Equal(5, suite.VariableThatShouldStartAtFive)
//	}
func (suite *AuthServiceSuite) TestGetMe() {
	// Test cases
	tests := []struct {
		name           string
		contextSetup   func() context.Context
		mockSetup      func()
		expectedResult response.MeResponse
		expectedError  error
	}{
		{
			name: "Success without domain ID from context",
			contextSetup: func() context.Context {
				return utility.NewContextWithValues(context.Background(), utility.ContextValues{
					UserID: 1,
				})
			},
			mockSetup: func() {
				suite.mockUserRepo.On("FindByID", mock.Anything, uint64(1), true).Return(
					entity.User{
						ID:       1,
						Username: "test",
					}, nil)
			},
			expectedResult: response.MeResponse{
				User: entity.User{
					ID:       1,
					Username: "test",
				},
			},
			expectedError: nil,
		},
		{
			name: "Success with domain ID from context",
			contextSetup: func() context.Context {
				return utility.NewContextWithValues(context.Background(), utility.ContextValues{
					UserID:   1,
					DomainID: utility.PointOf(uint64(1)),
				})
			},
			mockSetup: func() {
				suite.mockUserRepo.On("FindByID", mock.Anything, uint64(1), true).Return(
					entity.User{
						ID:       1,
						Username: "test",
					}, nil)
				suite.mockUserRoleRepo.On("FindByUserIDAndDomainID", mock.Anything, uint64(1), uint64(1), true).Return(
					[]entity.UserRole{
						{
							DomainID: 1,
							UserID:   1,
						},
					}, nil)
			},
			expectedResult: response.MeResponse{
				User: entity.User{
					ID:       1,
					Username: "test",
					UserRoles: []entity.UserRole{
						{
							DomainID: 1,
							UserID:   1,
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Error - Invalid context",
			contextSetup: func() context.Context {
				return context.Background()
			},
			mockSetup:      func() {},
			expectedResult: response.MeResponse{},
			expectedError:  errors.New("context values not found"),
		},
		{
			name: "Error - User not found",
			contextSetup: func() context.Context {
				return utility.NewContextWithValues(context.Background(), utility.ContextValues{
					UserID: 1,
				})
			},
			mockSetup: func() {
				suite.mockUserRepo.On("FindByID", mock.Anything, uint64(1), true).Return(
					entity.User{}, errors.New("user not found"))
			},
			expectedResult: response.MeResponse{},
			expectedError:  errors.New("user not found"),
		},
		{
			name: "Error - Failed to get user roles",
			contextSetup: func() context.Context {
				return utility.NewContextWithValues(context.Background(), utility.ContextValues{
					UserID:   1,
					DomainID: utility.PointOf(uint64(1)),
				})
			},
			mockSetup: func() {
				suite.mockUserRepo.On("FindByID", mock.Anything, uint64(1), true).Return(
					entity.User{
						ID: 1,
					}, nil)

				suite.mockUserRoleRepo.On("FindByUserIDAndDomainID", mock.Anything, uint64(1), uint64(1), true).Return(
					[]entity.UserRole{}, errors.New("failed to get user roles"))
			},
			expectedResult: response.MeResponse{},
			expectedError:  errors.New("failed to get user roles"),
		},
	}

	// Run tests
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			suite.New()
			// Setup mocks
			tt.mockSetup()

			// Execute
			result, err := suite.service.GetMe(tt.contextSetup())

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				if err != nil {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			// Verify all expectations were met
			suite.mockUserRepo.AssertExpectations(t)
			suite.mockUserRoleRepo.AssertExpectations(t)
		})
	}
}
