package validator

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Test model for unique validator testing
type UniqueTestModel struct {
	ID    uint   `gorm:"primarykey"`
	Email string `gorm:"unique"`
	Name  string
	Age   int
}

// setupUniqueTestDB creates an in-memory SQLite database for unique testing
func setupUniqueTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the test model
	err = db.AutoMigrate(&UniqueTestModel{})
	if err != nil {
		t.Fatalf("Failed to migrate test model: %v", err)
	}

	// Insert test data
	testData := []UniqueTestModel{
		{ID: 1, Email: "existing1@example.com", Name: "User One", Age: 25},
		{ID: 2, Email: "existing2@example.com", Name: "User Two", Age: 30},
		{ID: 3, Email: "existing3@example.com", Name: "User Three", Age: 35},
	}

	for _, data := range testData {
		db.Create(&data)
	}

	return db
}

func TestUnique_ValidatorFunction(t *testing.T) {
	db := setupUniqueTestDB(t)
	v := validator.New()

	// Register the unique validator
	customValidator := New(db)
	v.RegisterValidation("unique", customValidator.Unique())

	// Test struct with unique validation
	type TestStruct struct {
		UserEmail string `validate:"unique=unique_test_models.email"`
		UserName  string `validate:"unique=unique_test_models.name"`
		Optional  string `validate:"omitempty,unique=unique_test_models.email"`
	}

	tests := []struct {
		name        string
		input       TestStruct
		expectValid bool
		description string
	}{
		{
			name: "valid unique email",
			input: TestStruct{
				UserEmail: "new@example.com",
				UserName:  "New User",
				Optional:  "",
			},
			expectValid: true,
			description: "Should be valid with unique values",
		},
		{
			name: "valid unique name",
			input: TestStruct{
				UserEmail: "another@example.com",
				UserName:  "Another User",
				Optional:  "unique@example.com",
			},
			expectValid: true,
			description: "Should be valid with unique values and optional field",
		},
		{
			name: "invalid duplicate email",
			input: TestStruct{
				UserEmail: "existing1@example.com",
				UserName:  "New User",
				Optional:  "",
			},
			expectValid: false,
			description: "Should be invalid with duplicate email",
		},
		{
			name: "invalid duplicate name",
			input: TestStruct{
				UserEmail: "new@example.com",
				UserName:  "User One",
				Optional:  "",
			},
			expectValid: false,
			description: "Should be invalid with duplicate name",
		},
		{
			name: "invalid duplicate optional field",
			input: TestStruct{
				UserEmail: "new@example.com",
				UserName:  "New User",
				Optional:  "existing2@example.com",
			},
			expectValid: false,
			description: "Should be invalid with duplicate optional field",
		},
		{
			name: "valid empty optional field",
			input: TestStruct{
				UserEmail: "unique@example.com",
				UserName:  "Unique User",
				Optional:  "",
			},
			expectValid: true,
			description: "Should be valid with empty optional field (omitempty)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.input)
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

func TestUnique_InvalidTableColumnFormat(t *testing.T) {
	db := setupUniqueTestDB(t)
	v := validator.New()

	// Register the unique validator
	customValidator := New(db)
	v.RegisterValidation("unique", customValidator.Unique())

	// Test struct with invalid table.column format
	type TestStruct struct {
		InvalidFormat1 string `validate:"unique=invalid_format"`
		InvalidFormat2 string `validate:"unique=table.column.extra"`
		InvalidFormat3 string `validate:"unique="`
		ValidFormat    string `validate:"unique=unique_test_models.email"`
	}

	tests := []struct {
		name        string
		input       TestStruct
		expectValid bool
		description string
	}{
		{
			name: "invalid format without dot",
			input: TestStruct{
				InvalidFormat1: "test",
				InvalidFormat2: "test",
				InvalidFormat3: "test",
				ValidFormat:    "new@example.com",
			},
			expectValid: true, // Invalid formats should pass validation (fail open)
			description: "Should pass validation for invalid table.column formats",
		},
		{
			name: "valid format with unique value",
			input: TestStruct{
				InvalidFormat1: "anything",
				InvalidFormat2: "anything",
				InvalidFormat3: "anything",
				ValidFormat:    "unique@example.com",
			},
			expectValid: true,
			description: "Should pass validation with valid format and unique value",
		},
		{
			name: "valid format with duplicate value",
			input: TestStruct{
				InvalidFormat1: "anything",
				InvalidFormat2: "anything",
				InvalidFormat3: "anything",
				ValidFormat:    "existing1@example.com",
			},
			expectValid: false,
			description: "Should fail validation with valid format but duplicate value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.input)
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

func TestUnique_DatabaseError(t *testing.T) {
	// Create a validator with a closed database connection
	db := setupUniqueTestDB(t)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get underlying sql.DB: %v", err)
	}
	sqlDB.Close() // Close the connection to simulate database error

	v := validator.New()
	customValidator := New(db)
	v.RegisterValidation("unique", customValidator.Unique())

	type TestStruct struct {
		UserEmail string `validate:"unique=unique_test_models.email"`
	}

	tests := []struct {
		name        string
		input       TestStruct
		expectValid bool
		description string
	}{
		{
			name: "database error should pass validation",
			input: TestStruct{
				UserEmail: "test@example.com",
			},
			expectValid: true, // Should pass when database error occurs (fail open)
			description: "Should pass validation when database error occurs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.input)
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

func TestUnique_EdgeCases(t *testing.T) {
	db := setupUniqueTestDB(t)
	v := validator.New()
	customValidator := New(db)
	v.RegisterValidation("unique", customValidator.Unique())

	type TestStruct struct {
		Value string `validate:"unique=unique_test_models.email"`
	}

	tests := []struct {
		name        string
		input       TestStruct
		expectValid bool
		description string
	}{
		{
			name: "empty string",
			input: TestStruct{
				Value: "",
			},
			expectValid: true,
			description: "Should pass validation for empty string (unique constraint doesn't apply)",
		},
		{
			name: "whitespace only",
			input: TestStruct{
				Value: "   ",
			},
			expectValid: true,
			description: "Should pass validation for whitespace only (unique)",
		},
		{
			name: "case sensitive check",
			input: TestStruct{
				Value: "EXISTING1@EXAMPLE.COM",
			},
			expectValid: true,
			description: "Should pass validation for different case (case sensitive)",
		},
		{
			name: "special characters",
			input: TestStruct{
				Value: "test+special@example.com",
			},
			expectValid: true,
			description: "Should pass validation for unique value with special characters",
		},
		{
			name: "unicode characters",
			input: TestStruct{
				Value: "tëst@example.com",
			},
			expectValid: true,
			description: "Should pass validation for unique value with unicode characters",
		},
		{
			name: "very long string",
			input: TestStruct{
				Value: "verylongemailaddressthatisunique@verylongdomainnamethatisunusuallylong.com",
			},
			expectValid: true,
			description: "Should pass validation for very long unique value",
		},
		{
			name: "sql injection attempt",
			input: TestStruct{
				Value: "test'; DROP TABLE unique_test_models; --@example.com",
			},
			expectValid: true,
			description: "Should safely handle SQL injection attempts",
		},
		{
			name: "exact duplicate",
			input: TestStruct{
				Value: "existing1@example.com",
			},
			expectValid: false,
			description: "Should fail validation for exact duplicate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.input)
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

func TestUnique_NonExistentTable(t *testing.T) {
	db := setupUniqueTestDB(t)
	v := validator.New()
	customValidator := New(db)
	v.RegisterValidation("unique", customValidator.Unique())

	type TestStruct struct {
		Value string `validate:"unique=non_existent_table.email"`
	}

	tests := []struct {
		name        string
		input       TestStruct
		expectValid bool
		description string
	}{
		{
			name: "non-existent table",
			input: TestStruct{
				Value: "test@example.com",
			},
			expectValid: true,
			description: "Should pass validation for non-existent table (no records found = unique)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.input)
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

func TestUnique_NonExistentColumn(t *testing.T) {
	db := setupUniqueTestDB(t)
	v := validator.New()
	customValidator := New(db)
	v.RegisterValidation("unique", customValidator.Unique())

	type TestStruct struct {
		Value string `validate:"unique=unique_test_models.non_existent_column"`
	}

	tests := []struct {
		name        string
		input       TestStruct
		expectValid bool
		description string
	}{
		{
			name: "non-existent column",
			input: TestStruct{
				Value: "test@example.com",
			},
			expectValid: true,
			description: "Should pass validation for non-existent column (no records found = unique)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.input)
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

func TestUnique_MultipleValues(t *testing.T) {
	db := setupUniqueTestDB(t)
	v := validator.New()
	customValidator := New(db)
	v.RegisterValidation("unique", customValidator.Unique())

	// Add more test data with duplicate names but unique emails
	db.Create(&UniqueTestModel{ID: 4, Email: "user4@example.com", Name: "User One", Age: 40}) // Duplicate name
	db.Create(&UniqueTestModel{ID: 5, Email: "user5@example.com", Name: "User Two", Age: 45}) // Duplicate name

	type TestStruct struct {
		UserEmail string `validate:"unique=unique_test_models.email"`
		UserName  string `validate:"unique=unique_test_models.name"`
	}

	tests := []struct {
		name        string
		input       TestStruct
		expectValid bool
		description string
	}{
		{
			name: "unique email, duplicate name",
			input: TestStruct{
				UserEmail: "unique@example.com",
				UserName:  "User One", // This name already exists multiple times
			},
			expectValid: false,
			description: "Should fail validation due to duplicate name",
		},
		{
			name: "duplicate email, unique name",
			input: TestStruct{
				UserEmail: "existing1@example.com", // This email already exists
				UserName:  "Unique Name",
			},
			expectValid: false,
			description: "Should fail validation due to duplicate email",
		},
		{
			name: "both unique",
			input: TestStruct{
				UserEmail: "totallynew@example.com",
				UserName:  "Totally New User",
			},
			expectValid: true,
			description: "Should pass validation with both values unique",
		},
		{
			name: "both duplicate",
			input: TestStruct{
				UserEmail: "existing2@example.com",
				UserName:  "User Two",
			},
			expectValid: false,
			description: "Should fail validation due to both values being duplicate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.input)
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
