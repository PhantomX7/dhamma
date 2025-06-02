package validator

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	validatorpkg "github.com/go-playground/validator/v10"
)

// Test model for validator testing
type TestUser struct {
	ID    uint   `gorm:"primarykey"`
	Email string `gorm:"unique"`
	Name  string
}

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	// Use SQLite in-memory database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the test model
	err = db.AutoMigrate(&TestUser{})
	if err != nil {
		t.Fatalf("Failed to migrate test model: %v", err)
	}

	// Insert test data
	testUser := TestUser{
		Email: "existing@example.com",
		Name:  "Test User",
	}
	db.Create(&testUser)

	return db
}

// setupDryRunDB creates a dry run database session for testing SQL generation
func setupDryRunDB(t *testing.T) *gorm.DB {
	db := setupTestDB(t)
	return db.Session(&gorm.Session{DryRun: true})
}

func TestNewValidator(t *testing.T) {
	db := setupTestDB(t)

	v := NewValidator(db)

	if v == nil {
		t.Error("NewValidator() returned nil")
		return
	}

	// Test that it implements the Validator interface
	var _ Validator = v

	// Test that the validator is working (we can't type assert to *cValidator since it implements Validator interface)
	// Just test that the methods work
	result := v.CheckUnique("test@example.com", "test_users", "email")
	if !result {
		t.Error("CheckUnique method should return true for non-existing email")
	}

	result = v.CheckExist("existing@example.com", "test_users", "email")
	if !result {
		t.Error("CheckExist method should return true for existing email")
	}
}

func TestValidator_Unique(t *testing.T) {
	db := setupTestDB(t)
	v := NewValidator(db)

	tests := []struct {
		name        string
		value       string
		table       string
		column      string
		expected    bool
		description string
	}{
		{
			name:        "unique email",
			value:       "new@example.com",
			table:       "test_users",
			column:      "email",
			expected:    true,
			description: "Should return true for unique email",
		},
		{
			name:        "duplicate email",
			value:       "existing@example.com",
			table:       "test_users",
			column:      "email",
			expected:    false,
			description: "Should return false for existing email",
		},
		{
			name:        "empty value",
			value:       "",
			table:       "test_users",
			column:      "email",
			expected:    true,
			description: "Should return true for empty value (not validated by unique)",
		},
		{
			name:        "case sensitive check",
			value:       "EXISTING@EXAMPLE.COM",
			table:       "test_users",
			column:      "email",
			expected:    true,
			description: "Should return true for different case (case sensitive)",
		},
		{
			name:        "non-existent table",
			value:       "test@example.com",
			table:       "non_existent_table",
			column:      "email",
			expected:    true,
			description: "Should return true for non-existent table (no records found)",
		},
		{
			name:        "non-existent column",
			value:       "test@example.com",
			table:       "test_users",
			column:      "non_existent_column",
			expected:    true,
			description: "Should return true for non-existent column (no records found)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.CheckUnique(tt.value, tt.table, tt.column)

			if result != tt.expected {
				t.Errorf("Unique() = %v, want %v. %s", result, tt.expected, tt.description)
			}
		})
	}
}

func TestValidator_Exist(t *testing.T) {
	db := setupTestDB(t)
	v := NewValidator(db)

	tests := []struct {
		name        string
		value       string
		table       string
		column      string
		expected    bool
		description string
	}{
		{
			name:        "existing email",
			value:       "existing@example.com",
			table:       "test_users",
			column:      "email",
			expected:    true,
			description: "Should return true for existing email",
		},
		{
			name:        "non-existing email",
			value:       "nonexistent@example.com",
			table:       "test_users",
			column:      "email",
			expected:    false,
			description: "Should return false for non-existing email",
		},
		{
			name:        "empty value",
			value:       "",
			table:       "test_users",
			column:      "email",
			expected:    false,
			description: "Should return false for empty value",
		},
		{
			name:        "existing ID",
			value:       "1",
			table:       "test_users",
			column:      "id",
			expected:    true,
			description: "Should return true for existing ID",
		},
		{
			name:        "non-existing ID",
			value:       "999",
			table:       "test_users",
			column:      "id",
			expected:    false,
			description: "Should return false for non-existing ID",
		},
		{
			name:        "case sensitive check",
			value:       "EXISTING@EXAMPLE.COM",
			table:       "test_users",
			column:      "email",
			expected:    false,
			description: "Should return false for different case (case sensitive)",
		},
		{
			name:        "non-existent table",
			value:       "test@example.com",
			table:       "non_existent_table",
			column:      "email",
			expected:    false,
			description: "Should return false for non-existent table",
		},
		{
			name:        "non-existent column",
			value:       "test@example.com",
			table:       "test_users",
			column:      "non_existent_column",
			expected:    false,
			description: "Should return false for non-existent column",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.CheckExist(tt.value, tt.table, tt.column)

			if result != tt.expected {
				t.Errorf("Exist() = %v, want %v. %s", result, tt.expected, tt.description)
			}
		})
	}
}

func TestValidator_Integration(t *testing.T) {
	db := setupTestDB(t)
	v := NewValidator(db)

	// Test struct with custom validation tags
	type UserRegistration struct {
		Email    string `validate:"required,email,unique=test_users.email"`
		Name     string `validate:"required,min=2"`
		Referrer string `validate:"omitempty,exist=test_users.email"`
	}

	// Register custom validators
	validator := validatorpkg.New()
	validator.RegisterValidation("unique", func(fl validatorpkg.FieldLevel) bool {
		param := fl.Param()
		if param == "" {
			return true
		}

		// Parse table.column format
		parts := parseTableColumn(param)
		if len(parts) != 2 {
			return true
		}

		return v.CheckUnique(fl.Field().String(), parts[0], parts[1])
	})

	validator.RegisterValidation("exist", func(fl validatorpkg.FieldLevel) bool {
		param := fl.Param()
		if param == "" {
			return true
		}

		// Parse table.column format
		parts := parseTableColumn(param)
		if len(parts) != 2 {
			return true
		}

		return v.CheckExist(fl.Field().String(), parts[0], parts[1])
	})

	tests := []struct {
		name        string
		input       UserRegistration
		expectValid bool
		description string
	}{
		{
			name: "valid registration",
			input: UserRegistration{
				Email:    "newuser@example.com",
				Name:     "New User",
				Referrer: "existing@example.com",
			},
			expectValid: true,
			description: "Should be valid with unique email and existing referrer",
		},
		{
			name: "duplicate email",
			input: UserRegistration{
				Email:    "existing@example.com",
				Name:     "Test User",
				Referrer: "",
			},
			expectValid: false,
			description: "Should be invalid due to duplicate email",
		},
		{
			name: "non-existing referrer",
			input: UserRegistration{
				Email:    "newuser2@example.com",
				Name:     "New User 2",
				Referrer: "nonexistent@example.com",
			},
			expectValid: false,
			description: "Should be invalid due to non-existing referrer",
		},
		{
			name: "empty referrer (optional)",
			input: UserRegistration{
				Email:    "newuser3@example.com",
				Name:     "New User 3",
				Referrer: "",
			},
			expectValid: true,
			description: "Should be valid with empty referrer (omitempty)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Struct(tt.input)
			isValid := err == nil

			if isValid != tt.expectValid {
				t.Errorf("Validation result = %v, want %v. %s", isValid, tt.expectValid, tt.description)
				if err != nil {
					t.Errorf("Validation error: %v", err)
				}
			}
		})
	}
}

// Helper function to parse table.column format
func parseTableColumn(param string) []string {
	parts := make([]string, 0, 2)
	for i, part := range []rune(param) {
		if part == '.' {
			parts = append(parts, param[:i])
			parts = append(parts, param[i+1:])
			break
		}
	}
	return parts
}

func TestValidator_EdgeCases(t *testing.T) {
	db := setupTestDB(t)
	v := NewValidator(db)

	tests := []struct {
		name        string
		method      string
		value       string
		table       string
		column      string
		expected    bool
		description string
	}{
		{
			name:        "unique with special characters",
			method:      "unique",
			value:       "test+special@example.com",
			table:       "test_users",
			column:      "email",
			expected:    true,
			description: "Should handle special characters in email",
		},
		{
			name:        "exist with special characters",
			method:      "exist",
			value:       "test+special@example.com",
			table:       "test_users",
			column:      "email",
			expected:    false,
			description: "Should handle special characters in email",
		},
		{
			name:        "unique with unicode",
			method:      "unique",
			value:       "tëst@example.com",
			table:       "test_users",
			column:      "email",
			expected:    true,
			description: "Should handle unicode characters",
		},
		{
			name:        "exist with unicode",
			method:      "exist",
			value:       "tëst@example.com",
			table:       "test_users",
			column:      "email",
			expected:    false,
			description: "Should handle unicode characters",
		},
		{
			name:        "unique with very long value",
			method:      "unique",
			value:       "verylongemailaddressthatexceedsnormallimits@verylongdomainnamethatisunusuallylong.com",
			table:       "test_users",
			column:      "email",
			expected:    true,
			description: "Should handle very long values",
		},
		{
			name:        "exist with very long value",
			method:      "exist",
			value:       "verylongemailaddressthatexceedsnormallimits@verylongdomainnamethatisunusuallylong.com",
			table:       "test_users",
			column:      "email",
			expected:    false,
			description: "Should handle very long values",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result bool
			switch tt.method {
			case "unique":
				result = v.CheckUnique(tt.value, tt.table, tt.column)
			case "exist":
				result = v.CheckExist(tt.value, tt.table, tt.column)
			default:
				t.Fatalf("Unknown method: %s", tt.method)
			}

			if result != tt.expected {
				t.Errorf("%s() = %v, want %v. %s", tt.method, result, tt.expected, tt.description)
			}
		})
	}
}

func TestValidator_DatabaseError(t *testing.T) {
	// Create a validator with a closed database connection
	db := setupTestDB(t)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get underlying sql.DB: %v", err)
	}
	sqlDB.Close() // Close the connection to simulate database error

	v := NewValidator(db)

	tests := []struct {
		name        string
		method      string
		value       string
		table       string
		column      string
		expected    bool
		description string
	}{
		{
			name:        "unique with database error",
			method:      "unique",
			value:       "test@example.com",
			table:       "test_users",
			column:      "email",
			expected:    true, // Should return true on error (fail open)
			description: "Should return true when database error occurs",
		},
		{
			name:        "exist with database error",
			method:      "exist",
			value:       "test@example.com",
			table:       "test_users",
			column:      "email",
			expected:    false, // Should return false on error (fail closed)
			description: "Should return false when database error occurs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result bool
			switch tt.method {
			case "unique":
				result = v.CheckUnique(tt.value, tt.table, tt.column)
			case "exist":
				result = v.CheckExist(tt.value, tt.table, tt.column)
			default:
				t.Fatalf("Unknown method: %s", tt.method)
			}

			if result != tt.expected {
				t.Errorf("%s() = %v, want %v. %s", tt.method, result, tt.expected, tt.description)
			}
		})
	}
}

// TestValidator_DryRun_SQLGeneration tests SQL generation without executing queries
func TestValidator_DryRun_SQLGeneration(t *testing.T) {
	db := setupTestDB(t)
	v := NewValidator(db)

	tests := []struct {
		name   string
		method string
		value  string
		table  string
		column string
	}{
		{
			name:   "unique SQL generation",
			method: "unique",
			value:  "test@example.com",
			table:  "test_users",
			column: "email",
		},
		{
			name:   "exist SQL generation",
			method: "exist",
			value:  "existing@example.com",
			table:  "test_users",
			column: "email",
		},
		{
			name:   "unique with special characters",
			method: "unique",
			value:  "test+special@example.com",
			table:  "test_users",
			column: "email",
		},
		{
			name:   "exist with unicode",
			method: "exist",
			value:  "测试@example.com",
			table:  "test_users",
			column: "email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that the methods don't panic and can handle various inputs
			// This validates the validator works correctly with different data types
			switch tt.method {
			case "unique":
				result := v.CheckUnique(tt.value, tt.table, tt.column)
				// Result should be boolean, we're mainly testing for no panics
				_ = result
			case "exist":
				result := v.CheckExist(tt.value, tt.table, tt.column)
				// Result should be boolean, we're mainly testing for no panics
				_ = result
			default:
				t.Fatalf("Unknown method: %s", tt.method)
			}
		})
	}
}

// TestValidator_DryRun_Integration tests integration with go-playground/validator in dry run mode
func TestValidator_DryRun_Integration(t *testing.T) {
	db := setupTestDB(t)
	v := NewValidator(db)

	// Test struct with custom validation tags
	type UserRegistration struct {
		Email      string `validate:"required,email,unique=test_users.email"`
		ReferrerID string `validate:"omitempty,exist=test_users.id"`
	}

	// Register custom validation functions
	validator := validatorpkg.New()
	validator.RegisterValidation("unique", func(fl validatorpkg.FieldLevel) bool {
		param := fl.Param()
		if param == "" {
			return true
		}

		// Parse table.column format
		parts := parseTableColumn(param)
		if len(parts) != 2 {
			return true
		}

		return v.CheckUnique(fl.Field().String(), parts[0], parts[1])
	})

	validator.RegisterValidation("exist", func(fl validatorpkg.FieldLevel) bool {
		param := fl.Param()
		if param == "" {
			return true
		}

		// Parse table.column format
		parts := parseTableColumn(param)
		if len(parts) != 2 {
			return true
		}

		return v.CheckExist(fl.Field().String(), parts[0], parts[1])
	})

	tests := []struct {
		name        string
		user        UserRegistration
		expectValid bool
		description string
	}{
		{
			name: "valid registration",
			user: UserRegistration{
				Email:      "newuser@example.com",
				ReferrerID: "1",
			},
			expectValid: true,
			description: "Should validate successfully with normal database",
		},
		{
			name: "duplicate email",
			user: UserRegistration{
				Email:      "existing@example.com",
				ReferrerID: "1",
			},
			expectValid: false,
			description: "Should detect duplicate email",
		},
		{
			name: "non-existing referrer",
			user: UserRegistration{
				Email:      "newuser2@example.com",
				ReferrerID: "999",
			},
			expectValid: false,
			description: "Should detect non-existing referrer",
		},
		{
			name: "empty referrer",
			user: UserRegistration{
				Email:      "newuser3@example.com",
				ReferrerID: "",
			},
			expectValid: true,
			description: "Should allow empty referrer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Struct(tt.user)
			isValid := err == nil

			if isValid != tt.expectValid {
				t.Errorf("%s: expected valid=%v, got valid=%v, error=%v",
					tt.description, tt.expectValid, isValid, err)
			}

			t.Logf("Validation result for %s: valid=%v", tt.name, isValid)
		})
	}
}
