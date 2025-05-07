package seed

import (
	"errors"
	"log"
	"os"
	"strings"

	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/entity"

	"github.com/go-faker/faker/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserSeeder struct {
	db *gorm.DB
}

func NewUserSeeder(db *gorm.DB) *UserSeeder {
	return &UserSeeder{db: db}
}

func (s *UserSeeder) GenerateRootUser() (err error) {
	rootUser := []entity.User{
		{
			Username:     os.Getenv("ADMIN_USERNAME"),
			Password:     os.Getenv("ADMIN_PASSWORD"),
			IsSuperAdmin: true,
			IsActive:     true,
		},
	}

	log.Print("seeding root user")
	for _, user := range rootUser {
		if !errors.Is(s.db.First(&entity.User{}, entity.User{
			Username: user.Username,
		}).Error, gorm.ErrRecordNotFound) {
			continue
		}

		var password []byte
		password, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(password)

		err = s.db.Create(&user).Error
		if err != nil {
			return err
		}
	}

	return err
}

func (s *UserSeeder) GenerateUsers(count int, opts ...UserOption) error {

	log.Printf("seeding %d users", count)
	// Default options
	options := &UserOptions{
		password:     "password123", // default password
		isActive:     true,          // default active status
		isSuperAdmin: false,         // default not super admin
	}

	// Apply provided options
	for _, opt := range opts {
		opt(options)
	}

	users := make([]entity.User, count)

	for i := 0; i < count; i++ {
		// Generate fake data
		username := strings.ToLower(faker.FirstName())

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(options.password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		users[i] = entity.User{
			Username:     username, // using email as username
			Password:     string(hashedPassword),
			IsActive:     options.isActive,
			IsSuperAdmin: options.isSuperAdmin,
		}
	}

	// Batch insert users
	return s.db.CreateInBatches(users, 100).Error
}

type UserOptions struct {
	password     string
	isActive     bool
	isSuperAdmin bool
}

type UserOption func(*UserOptions)

func WithPassword(password string) UserOption {
	return func(o *UserOptions) {
		o.password = password
	}
}

func WithActiveStatus(isActive bool) UserOption {
	return func(o *UserOptions) {
		o.isActive = isActive
	}
}

func WithSuperAdmin(isSuperAdmin bool) UserOption {
	return func(o *UserOptions) {
		o.isSuperAdmin = isSuperAdmin
	}
}
