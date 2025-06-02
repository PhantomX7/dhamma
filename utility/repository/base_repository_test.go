package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/PhantomX7/dhamma/utility/pagination"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestUser represents a test entity for repository testing
type TestUser struct {
	ID        uint64       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string       `gorm:"size:100;not null" json:"name"`
	Email     string       `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Age       int          `gorm:"not null" json:"age"`
	Active    bool         `json:"active"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Profile   *TestProfile `gorm:"foreignKey:UserID" json:"profile,omitempty"`
}

// TestProfile represents a related entity for testing preloads
type TestProfile struct {
	ID     uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID uint64 `gorm:"not null" json:"user_id"`
	Bio    string `gorm:"size:500" json:"bio"`
}

// setupTestDB creates an in-memory SQLite database and returns a GORM database connection
func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	// Connect to in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Enable logging for debugging
	})
	require.NoError(t, err)

	// Auto migrate test tables
	err = db.AutoMigrate(&TestUser{}, &TestProfile{})
	require.NoError(t, err)

	// Verify tables were created
	var count int64
	err = db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name IN ('test_users', 'test_profiles')").Scan(&count).Error
	require.NoError(t, err)
	t.Logf("Created %d tables", count)

	// Return cleanup function (no-op for in-memory database)
	cleanup := func() {
		// No cleanup needed for in-memory database
	}

	return db, cleanup
}

// seedTestData creates sample data for testing
func seedTestData(t *testing.T, db *gorm.DB) {
	users := []TestUser{
		{Name: "John Doe", Email: "john@example.com", Age: 30, Active: true},
		{Name: "Jane Smith", Email: "jane@example.com", Age: 25, Active: true},
		{Name: "Bob Johnson", Email: "bob@example.com", Age: 35, Active: false},
	}

	for i, user := range users {
		err := db.Create(&user).Error
		require.NoError(t, err)
		t.Logf("Created user %d: %+v", i+1, user)
	}

	// Verify users were inserted
	var userCount int64
	err := db.Model(&TestUser{}).Count(&userCount).Error
	require.NoError(t, err)
	t.Logf("Total users in database: %d", userCount)

	// Create profiles for first two users
	profiles := []TestProfile{
		{UserID: 1, Bio: "Software Engineer"},
		{UserID: 2, Bio: "Product Manager"},
	}

	for i, profile := range profiles {
		err := db.Create(&profile).Error
		require.NoError(t, err)
		t.Logf("Created profile %d: %+v", i+1, profile)
	}
}

func TestNewBaseRepository(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBaseRepository[TestUser](db)

	assert.NotNil(t, repo.DB)
	assert.Equal(t, db, repo.DB)
}

func TestBaseRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBaseRepository[TestUser](db)
	ctx := context.Background()

	t.Run("successful create", func(t *testing.T) {
		user := &TestUser{
			Name:   "Test User",
			Email:  "test@example.com",
			Age:    28,
			Active: true,
		}

		err := repo.Create(ctx, user, nil)
		assert.NoError(t, err)
		assert.NotZero(t, user.ID)
		assert.NotZero(t, user.CreatedAt)
		assert.NotZero(t, user.UpdatedAt)
	})

	t.Run("create with transaction", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		user := &TestUser{
			Name:   "TX User",
			Email:  "tx@example.com",
			Age:    30,
			Active: true,
		}

		err := repo.Create(ctx, user, tx)
		assert.NoError(t, err)
		assert.NotZero(t, user.ID)
	})

	t.Run("create with duplicate email should fail", func(t *testing.T) {
		// First user
		user1 := &TestUser{
			Name:   "User 1",
			Email:  "duplicate@example.com",
			Age:    25,
			Active: true,
		}
		err := repo.Create(ctx, user1, nil)
		assert.NoError(t, err)

		// Second user with same email
		user2 := &TestUser{
			Name:   "User 2",
			Email:  "duplicate@example.com",
			Age:    30,
			Active: true,
		}
		err = repo.Create(ctx, user2, nil)
		assert.Error(t, err)
	})
}

func TestBaseRepository_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBaseRepository[TestUser](db)
	ctx := context.Background()

	// Create initial user
	user := &TestUser{
		Name:   "Original Name",
		Email:  "original@example.com",
		Age:    25,
		Active: true,
	}
	err := repo.Create(ctx, user, nil)
	require.NoError(t, err)

	t.Run("successful update", func(t *testing.T) {
		user.Name = "Updated Name"
		user.Age = 30

		err := repo.Update(ctx, user, nil)
		assert.NoError(t, err)

		// Verify update
		updatedUser, err := repo.FindByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Name", updatedUser.Name)
		assert.Equal(t, 30, updatedUser.Age)
	})

	t.Run("update with transaction", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		user.Name = "TX Updated Name"
		err := repo.Update(ctx, user, tx)
		assert.NoError(t, err)
	})
}

func TestBaseRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBaseRepository[TestUser](db)
	ctx := context.Background()

	// Create user to delete
	user := &TestUser{
		Name:   "To Delete",
		Email:  "delete@example.com",
		Age:    25,
		Active: true,
	}
	err := repo.Create(ctx, user, nil)
	require.NoError(t, err)

	t.Run("successful delete", func(t *testing.T) {
		err := repo.Delete(ctx, user, nil)
		assert.NoError(t, err)

		// Verify deletion
		_, err = repo.FindByID(ctx, user.ID)
		assert.Error(t, err)
	})

	t.Run("delete with transaction", func(t *testing.T) {
		// Create another user
		txUser := &TestUser{
			Name:   "TX Delete",
			Email:  "txdelete@example.com",
			Age:    25,
			Active: true,
		}
		err := repo.Create(ctx, txUser, nil)
		require.NoError(t, err)

		tx := db.Begin()
		defer tx.Rollback()

		err = repo.Delete(ctx, txUser, tx)
		assert.NoError(t, err)
	})
}

func TestBaseRepository_FindByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBaseRepository[TestUser](db)
	ctx := context.Background()

	seedTestData(t, db)

	t.Run("find existing user", func(t *testing.T) {
		user, err := repo.FindByID(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, uint64(1), user.ID)
		assert.Equal(t, "John Doe", user.Name)
		assert.Equal(t, "john@example.com", user.Email)
	})

	t.Run("find with preload", func(t *testing.T) {
		user, err := repo.FindByID(ctx, 1, "Profile")
		assert.NoError(t, err)
		assert.Equal(t, uint64(1), user.ID)
		assert.NotNil(t, user.Profile)
		assert.Equal(t, "Software Engineer", user.Profile.Bio)
	})

	t.Run("find non-existing user", func(t *testing.T) {
		_, err := repo.FindByID(ctx, 999)
		assert.Error(t, err)
	})
}

func TestBaseRepository_FindAll(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBaseRepository[TestUser](db)
	ctx := context.Background()

	seedTestData(t, db)

	t.Run("find all without pagination", func(t *testing.T) {
		pg := pagination.NewPagination(map[string][]string{}, &pagination.FilterDefinition{}, pagination.PaginationOptions{
			DefaultLimit: 10,
			MaxLimit:     100,
			DefaultOrder: "id asc",
		})

		users, err := repo.FindAll(ctx, pg)
		assert.NoError(t, err)
		assert.Len(t, users, 3)
	})

	t.Run("find all with limit", func(t *testing.T) {
		pg := pagination.NewPagination(map[string][]string{
			"limit": {"2"},
		}, &pagination.FilterDefinition{}, pagination.PaginationOptions{
			DefaultLimit: 10,
			MaxLimit:     100,
			DefaultOrder: "id asc",
		})

		users, err := repo.FindAll(ctx, pg)
		assert.NoError(t, err)
		assert.Len(t, users, 2)
	})

	t.Run("find all with offset", func(t *testing.T) {
		pg := pagination.NewPagination(map[string][]string{
			"limit":  {"2"},
			"offset": {"1"},
		}, &pagination.FilterDefinition{}, pagination.PaginationOptions{
			DefaultLimit: 10,
			MaxLimit:     100,
			DefaultOrder: "id asc",
		})

		users, err := repo.FindAll(ctx, pg)
		assert.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, "Jane Smith", users[0].Name)
	})

	t.Run("find all with filters", func(t *testing.T) {
		// Create filter definition for active field
		filterDef := pagination.NewFilterDefinition()
		filterDef.AddFilter("active", pagination.FilterConfig{
			Field:     "active",
			Type:      pagination.FilterTypeBool,
			Operators: []pagination.FilterOperator{pagination.OperatorEquals},
		})

		pg := pagination.NewPagination(map[string][]string{
			"active": {"true"},
		}, filterDef, pagination.PaginationOptions{
			DefaultLimit: 10,
			MaxLimit:     100,
			DefaultOrder: "id asc",
		})
		users, err := repo.FindAll(ctx, pg)
		assert.NoError(t, err)
		assert.Len(t, users, 2) // Only active users
	})

	t.Run("find all with preloads", func(t *testing.T) {
		pg := pagination.NewPagination(map[string][]string{}, &pagination.FilterDefinition{}, pagination.PaginationOptions{
			DefaultLimit: 10,
			MaxLimit:     100,
			DefaultOrder: "id asc",
		})

		// Add preload scope
		pg.AddCustomScope(func(db *gorm.DB) *gorm.DB {
			return db.Preload("Profile")
		})
		users, err := repo.FindAll(ctx, pg)
		assert.NoError(t, err)
		assert.Len(t, users, 3)
		// Check that profiles are loaded for users that have them
		for _, user := range users {
			if user.ID <= 2 {
				assert.NotNil(t, user.Profile)
			}
		}
	})
}

func TestBaseRepository_Count(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBaseRepository[TestUser](db)
	ctx := context.Background()

	seedTestData(t, db)

	t.Run("count all users", func(t *testing.T) {
		pg := pagination.NewPagination(map[string][]string{}, &pagination.FilterDefinition{}, pagination.PaginationOptions{
			DefaultLimit: 10,
			MaxLimit:     100,
			DefaultOrder: "id asc",
		})

		count, err := repo.Count(ctx, pg)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})

	t.Run("count with filters", func(t *testing.T) {
		// Create filter definition for active field
		filterDef := pagination.NewFilterDefinition()
		filterDef.AddFilter("active", pagination.FilterConfig{
			Field:     "active",
			Type:      pagination.FilterTypeBool,
			Operators: []pagination.FilterOperator{pagination.OperatorEquals},
		})

		pg := pagination.NewPagination(map[string][]string{
			"active": {"true"},
		}, filterDef, pagination.PaginationOptions{
			DefaultLimit: 10,
			MaxLimit:     100,
			DefaultOrder: "id asc",
		})
		count, err := repo.Count(ctx, pg)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)
	})
}

func TestBaseRepository_FindByField(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBaseRepository[TestUser](db)
	ctx := context.Background()

	seedTestData(t, db)

	t.Run("find by email", func(t *testing.T) {
		users, err := repo.FindByField(ctx, "email", "john@example.com")
		assert.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, "John Doe", users[0].Name)
	})

	t.Run("find by active status", func(t *testing.T) {
		users, err := repo.FindByField(ctx, "active", true)
		assert.NoError(t, err)
		assert.Len(t, users, 2)
	})

	t.Run("find by non-existing value", func(t *testing.T) {
		users, err := repo.FindByField(ctx, "email", "nonexistent@example.com")
		assert.NoError(t, err)
		assert.Len(t, users, 0)
	})

	t.Run("find with preload", func(t *testing.T) {
		users, err := repo.FindByField(ctx, "email", "john@example.com", "Profile")
		assert.NoError(t, err)
		assert.Len(t, users, 1)
		assert.NotNil(t, users[0].Profile)
	})
}

func TestBaseRepository_FindOneByField(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBaseRepository[TestUser](db)
	ctx := context.Background()

	seedTestData(t, db)

	t.Run("find one by email", func(t *testing.T) {
		user, err := repo.FindOneByField(ctx, "email", "john@example.com")
		assert.NoError(t, err)
		assert.Equal(t, "John Doe", user.Name)
	})

	t.Run("find one by non-existing value", func(t *testing.T) {
		_, err := repo.FindOneByField(ctx, "email", "nonexistent@example.com")
		assert.Error(t, err)
	})

	t.Run("find one with preload", func(t *testing.T) {
		user, err := repo.FindOneByField(ctx, "email", "john@example.com", "Profile")
		assert.NoError(t, err)
		assert.Equal(t, "John Doe", user.Name)
		assert.NotNil(t, user.Profile)
	})
}

func TestBaseRepository_FindByFields(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBaseRepository[TestUser](db)
	ctx := context.Background()

	seedTestData(t, db)

	t.Run("find by multiple fields", func(t *testing.T) {
		conditions := map[string]any{
			"active": true,
			"age":    30,
		}
		users, err := repo.FindByFields(ctx, conditions)
		assert.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, "John Doe", users[0].Name)
	})

	t.Run("find by fields with no matches", func(t *testing.T) {
		conditions := map[string]any{
			"active": true,
			"age":    99,
		}
		users, err := repo.FindByFields(ctx, conditions)
		assert.NoError(t, err)
		assert.Len(t, users, 0)
	})

	t.Run("find by fields with preload", func(t *testing.T) {
		conditions := map[string]any{
			"email": "john@example.com",
		}
		users, err := repo.FindByFields(ctx, conditions, "Profile")
		assert.NoError(t, err)
		assert.Len(t, users, 1)
		assert.NotNil(t, users[0].Profile)
	})
}

func TestBaseRepository_FindOneByFields(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBaseRepository[TestUser](db)
	ctx := context.Background()

	seedTestData(t, db)

	t.Run("find one by multiple fields", func(t *testing.T) {
		conditions := map[string]any{
			"name":  "John Doe",
			"email": "john@example.com",
		}
		user, err := repo.FindOneByFields(ctx, conditions)
		assert.NoError(t, err)
		assert.Equal(t, "John Doe", user.Name)
	})

	t.Run("find one by fields with no match", func(t *testing.T) {
		conditions := map[string]any{
			"name": "Non Existent",
			"age":  99,
		}
		_, err := repo.FindOneByFields(ctx, conditions)
		assert.Error(t, err)
	})

	t.Run("find one by fields with preload", func(t *testing.T) {
		conditions := map[string]any{
			"email": "john@example.com",
		}
		user, err := repo.FindOneByFields(ctx, conditions, "Profile")
		assert.NoError(t, err)
		assert.Equal(t, "John Doe", user.Name)
		assert.NotNil(t, user.Profile)
	})
}

func TestBaseRepository_Exists(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewBaseRepository[TestUser](db)
	ctx := context.Background()

	seedTestData(t, db)

	t.Run("exists with matching conditions", func(t *testing.T) {
		conditions := map[string]any{
			"email": "john@example.com",
		}
		exists, err := repo.Exists(ctx, conditions)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("exists with non-matching conditions", func(t *testing.T) {
		conditions := map[string]any{
			"email": "nonexistent@example.com",
		}
		exists, err := repo.Exists(ctx, conditions)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("exists with multiple conditions", func(t *testing.T) {
		conditions := map[string]any{
			"name":   "John Doe",
			"active": true,
		}
		exists, err := repo.Exists(ctx, conditions)
		assert.NoError(t, err)
		assert.True(t, exists)
	})
}

// Benchmark tests
func BenchmarkBaseRepository_Create(b *testing.B) {
	db, cleanup := setupTestDB(&testing.T{})
	defer cleanup()

	repo := NewBaseRepository[TestUser](db)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := &TestUser{
			Name:   fmt.Sprintf("User %d", i),
			Email:  fmt.Sprintf("user%d@example.com", i),
			Age:    25 + (i % 50),
			Active: i%2 == 0,
		}
		_ = repo.Create(ctx, user, nil)
	}
}

func BenchmarkBaseRepository_FindByID(b *testing.B) {
	db, cleanup := setupTestDB(&testing.T{})
	defer cleanup()

	repo := NewBaseRepository[TestUser](db)
	ctx := context.Background()

	// Create test data
	for i := 1; i <= 100; i++ {
		user := &TestUser{
			Name:   fmt.Sprintf("User %d", i),
			Email:  fmt.Sprintf("user%d@example.com", i),
			Age:    25 + (i % 50),
			Active: i%2 == 0,
		}
		_ = repo.Create(ctx, user, nil)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := uint64((i % 100) + 1)
		_, _ = repo.FindByID(ctx, id)
	}
}
