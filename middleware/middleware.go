package middleware

import (
	"context"
	"log"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/libs/casbin"
	"github.com/PhantomX7/dhamma/modules/domain"
	"github.com/PhantomX7/dhamma/modules/permission" // Import permission module
	"github.com/PhantomX7/dhamma/modules/refresh_token"
	"github.com/PhantomX7/dhamma/modules/user"
	"github.com/PhantomX7/dhamma/modules/user_domain"
	"github.com/PhantomX7/dhamma/utility/pagination"

	"go.uber.org/zap"
)

type Middleware struct {
	userRepo              user.Repository
	refreshTokenRepo      refresh_token.Repository
	userDomainRepo        user_domain.Repository
	domainRepo            domain.Repository
	permissionRepo        permission.Repository
	casbin                casbin.Client
	logger                *zap.Logger
	permissionDefinitions map[string]entity.Permission // Add map to store definitions
}

func New(
	userRepo user.Repository,
	refreshTokenRepo refresh_token.Repository,
	userDomainRepo user_domain.Repository,
	domainRepo domain.Repository,
	permissionRepo permission.Repository,
	casbin casbin.Client,
	logger *zap.Logger,
) *Middleware {

	// Load permission definitions at startup
	permissionDefsMap := make(map[string]entity.Permission)

	// Using nil pagination to fetch all permissions. Adjust if you have many permissions.
	allPermissions, err := permissionRepo.FindAll(
		context.Background(), pagination.NewPagination(nil, nil, pagination.PaginationOptions{
			DefaultLimit: 1000,
		}))
	if err != nil {
		// Log the error and potentially panic or handle gracefully
		log.Fatalf("Failed to load permission definitions: %v", err)
	}
	for _, p := range allPermissions {
		permissionDefsMap[p.Code] = p
	}
	log.Printf("Loaded %d permission definitions", len(permissionDefsMap))

	return &Middleware{
		userRepo:              userRepo,
		refreshTokenRepo:      refreshTokenRepo,
		userDomainRepo:        userDomainRepo,
		domainRepo:            domainRepo,
		permissionRepo:        permissionRepo, // Store injected repo
		casbin:                casbin,
		logger:                logger,
		permissionDefinitions: permissionDefsMap, // Store loaded definitions
	}
}
