package service_test

import (
	transactionManagerMocks "github.com/PhantomX7/dhamma/libs/transaction_manager/mocks"
	"github.com/PhantomX7/dhamma/modules/auth"
	authService "github.com/PhantomX7/dhamma/modules/auth/service"
	refreshTokenMocks "github.com/PhantomX7/dhamma/modules/refresh_token/mocks"
	userMocks "github.com/PhantomX7/dhamma/modules/user/mocks"
	userRoleMocks "github.com/PhantomX7/dhamma/modules/user_role/mocks"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AuthServiceSuite struct {
	suite.Suite
	service                auth.Service
	mockUserRepo           *userMocks.Repository
	mockUserRoleRepo       *userRoleMocks.Repository
	mockRefreshTokenRepo   *refreshTokenMocks.Repository
	mockTransactionManager *transactionManagerMocks.Client
}

// before each test
func (suite *AuthServiceSuite) New() {
	// Initialize mocks
	mockUserRepo := new(userMocks.Repository)
	mockUserRoleRepo := new(userRoleMocks.Repository)
	mockRefreshTokenRepo := new(refreshTokenMocks.Repository)
	mockTransactionManager := new(transactionManagerMocks.Client)

	suite.service = authService.New(
		mockUserRepo,
		mockUserRoleRepo,
		mockRefreshTokenRepo,
		mockTransactionManager,
	)

	suite.mockUserRepo = mockUserRepo
	suite.mockUserRoleRepo = mockUserRoleRepo
	suite.mockRefreshTokenRepo = mockRefreshTokenRepo
	suite.mockTransactionManager = mockTransactionManager
}

func TestAuthServiceSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceSuite))
}
