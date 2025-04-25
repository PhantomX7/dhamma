package middleware

import (
	"context"

	"go.uber.org/zap"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/libs/casbin"
	"github.com/PhantomX7/dhamma/modules/domain"
	"github.com/PhantomX7/dhamma/modules/permission"
	"github.com/PhantomX7/dhamma/modules/refresh_token"
	"github.com/PhantomX7/dhamma/modules/user"
	"github.com/PhantomX7/dhamma/modules/user_domain"
	"github.com/PhantomX7/dhamma/utility/logger"
	"github.com/PhantomX7/dhamma/utility/pagination"
)

type Middleware struct {
	userRepo              user.Repository
	refreshTokenRepo      refresh_token.Repository
	userDomainRepo        user_domain.Repository
	domainRepo            domain.Repository
	permissionRepo        permission.Repository
	casbin                casbin.Client
	permissionDefinitions map[string]entity.Permission // Add map to store definitions
}

// New creates a new Middleware instance.
// It loads permission definitions from the repository at startup.
func New(
	userRepo user.Repository,
	refreshTokenRepo refresh_token.Repository,
	userDomainRepo user_domain.Repository,
	domainRepo domain.Repository,
	permissionRepo permission.Repository,
	casbin casbin.Client,
) *Middleware {

	// Load permission definitions at startup
	permissionDefsMap := make(map[string]entity.Permission)

	// Using nil pagination to fetch all permissions. Adjust if you have many permissions.
	allPermissions, err := permissionRepo.FindAll(
		context.Background(), pagination.NewPagination(nil, nil, pagination.PaginationOptions{
			DefaultLimit: 1000,
		}))
	if err != nil {
		logger.Get().Fatal("Failed to load permission definitions", zap.Error(err))
	}
	for _, p := range allPermissions {
		permissionDefsMap[p.Code] = p
	}
	logger.Get().Info("Loaded permission definitions", zap.Int("count", len(permissionDefsMap)))

	return &Middleware{
		userRepo:              userRepo,
		refreshTokenRepo:      refreshTokenRepo,
		userDomainRepo:        userDomainRepo,
		domainRepo:            domainRepo,
		permissionRepo:        permissionRepo,
		casbin:                casbin,
		permissionDefinitions: permissionDefsMap,
	}
}
