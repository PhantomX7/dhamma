package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/modules/user"
	"github.com/PhantomX7/dhamma/utility/pagination"
	baseRepo "github.com/PhantomX7/dhamma/utility/repository"
)

type repository struct {
	base baseRepo.BaseRepository[entity.User]
	db   *gorm.DB
}

func New(db *gorm.DB) user.Repository {
	return &repository{
		base: baseRepo.NewBaseRepository[entity.User](db),
		db:   db,
	}
}

// FindAll finds all users with pagination
func (r *repository) FindAll(ctx context.Context, pg *pagination.Pagination) ([]entity.User, error) {
	return r.base.FindAll(ctx, pg)
}

// FindByID finds a user by ID
func (r *repository) FindByID(ctx context.Context, id uint64, withPreload bool) (entity.User, error) {
	var preloads []string
	if withPreload {
		preloads = []string{"Domains", "UserRoles.Role"}
	}
	return r.base.FindByID(ctx, id, preloads...)
}

// Create creates a new user
func (r *repository) Create(ctx context.Context, user *entity.User, tx *gorm.DB) error {
	return r.base.Create(ctx, user, tx)
}

// Update updates a user
func (r *repository) Update(ctx context.Context, user *entity.User, tx *gorm.DB) error {
	return r.base.Update(ctx, user, tx)
}

// Count counts users with pagination
func (r *repository) Count(ctx context.Context, pg *pagination.Pagination) (int64, error) {
	return r.base.Count(ctx, pg)
}
