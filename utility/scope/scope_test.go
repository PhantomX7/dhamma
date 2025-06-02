package scope

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Test model for testing scopes
type TestModel struct {
	ID     uint `gorm:"primarykey"`
	Name   string
	Age    int
	Active bool
	Email  *string
}

// stringPtr returns a pointer to the given string
func stringPtr(s string) *string {
	return &s
}

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the test model
	err = db.AutoMigrate(&TestModel{})
	if err != nil {
		t.Fatalf("Failed to migrate test model: %v", err)
	}

	// Insert test data
	testData := []TestModel{
		{Name: "Alice", Age: 25, Active: true, Email: stringPtr("alice@example.com")},
		{Name: "Bob", Age: 30, Active: false, Email: stringPtr("bob@example.com")},
		{Name: "Charlie", Age: 35, Active: true, Email: stringPtr("charlie@example.com")},
		{Name: "David", Age: 40, Active: false, Email: stringPtr("david@example.com")},
		{Name: "Eve", Age: 28, Active: true, Email: stringPtr("eve@example.com")},
	}

	for _, data := range testData {
		result := db.Create(&data)
		if result.Error != nil {
			t.Fatalf("Failed to insert test data: %v", result.Error)
		}
	}

	return db
}

func TestLimitScope(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name          string
		limit         int
		expectedCount int64
	}{
		{
			name:          "limit 2",
			limit:         2,
			expectedCount: 2,
		},
		{
			name:          "limit 0 (no limit)",
			limit:         0,
			expectedCount: 5, // All records
		},
		{
			name:          "limit greater than total records",
			limit:         10,
			expectedCount: 5, // All records
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []TestModel
			var count int64

			db.Scopes(LimitScope(tt.limit)).Find(&results).Count(&count)

			if int64(len(results)) != tt.expectedCount {
				t.Errorf("LimitScope() returned %d records, want %d", len(results), tt.expectedCount)
			}
		})
	}
}

func TestOffsetScope(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name            string
		offset          int
		expectedCount   int64
		expectedFirstID uint
	}{
		{
			name:            "offset 0",
			offset:          0,
			expectedCount:   5,
			expectedFirstID: 1,
		},
		{
			name:            "offset 2",
			offset:          2,
			expectedCount:   3,
			expectedFirstID: 3,
		},
		{
			name:          "offset greater than total records",
			offset:        10,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []TestModel

			db.Scopes(OffsetScope(tt.offset)).Find(&results)

			if int64(len(results)) != tt.expectedCount {
				t.Errorf("OffsetScope() returned %d records, want %d", len(results), tt.expectedCount)
			}

			if tt.expectedCount > 0 && len(results) > 0 {
				if results[0].ID != tt.expectedFirstID {
					t.Errorf("OffsetScope() first record ID = %d, want %d", results[0].ID, tt.expectedFirstID)
				}
			}
		})
	}
}

func TestOrderScope(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name                string
		orderBy             string
		expectedFirstRecord TestModel
		expectedLastRecord  TestModel
	}{
		{
			name:                "order by age ASC",
			orderBy:             "age ASC",
			expectedFirstRecord: TestModel{ID: 1, Name: "Alice", Age: 25},
			expectedLastRecord:  TestModel{ID: 4, Name: "David", Age: 40},
		},
		{
			name:                "order by age DESC",
			orderBy:             "age DESC",
			expectedFirstRecord: TestModel{ID: 4, Name: "David", Age: 40},
			expectedLastRecord:  TestModel{ID: 1, Name: "Alice", Age: 25},
		},
		{
			name:                "order by name ASC",
			orderBy:             "name ASC",
			expectedFirstRecord: TestModel{ID: 1, Name: "Alice"},
			expectedLastRecord:  TestModel{ID: 5, Name: "Eve"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []TestModel

			db.Scopes(OrderScope(tt.orderBy)).Find(&results)

			if len(results) == 0 {
				t.Error("OrderScope() returned no records")
				return
			}

			// Check first record
			if results[0].ID != tt.expectedFirstRecord.ID {
				t.Errorf("OrderScope() first record ID = %d, want %d", results[0].ID, tt.expectedFirstRecord.ID)
			}

			// Check last record
			lastIndex := len(results) - 1
			if results[lastIndex].ID != tt.expectedLastRecord.ID {
				t.Errorf("OrderScope() last record ID = %d, want %d", results[lastIndex].ID, tt.expectedLastRecord.ID)
			}
		})
	}
}

func TestWhereNotInScope(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name          string
		column        string
		values        []interface{}
		expectedCount int64
		expectedIDs   []uint
	}{
		{
			name:          "exclude specific IDs",
			column:        "id",
			values:        []interface{}{1, 3},
			expectedCount: 3,
			expectedIDs:   []uint{2, 4, 5},
		},
		{
			name:          "exclude specific ages",
			column:        "age",
			values:        []interface{}{25, 30},
			expectedCount: 3,
			expectedIDs:   []uint{3, 4, 5},
		},
		{
			name:          "empty values list",
			column:        "id",
			values:        []interface{}{},
			expectedCount: 5, // All records
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []TestModel

			db.Scopes(WhereNotInScope(tt.column, tt.values)).Find(&results)

			if int64(len(results)) != tt.expectedCount {
				t.Errorf("WhereNotInScope() returned %d records, want %d", len(results), tt.expectedCount)
			}

			if len(tt.expectedIDs) > 0 {
				actualIDs := make([]uint, len(results))
				for i, result := range results {
					actualIDs[i] = result.ID
				}

				for _, expectedID := range tt.expectedIDs {
					found := false
					for _, actualID := range actualIDs {
						if actualID == expectedID {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("WhereNotInScope() missing expected ID %d in results", expectedID)
					}
				}
			}
		})
	}
}

func TestWhereInScope(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name          string
		column        string
		values        []interface{}
		expectedCount int64
		expectedIDs   []uint
	}{
		{
			name:          "include specific IDs",
			column:        "id",
			values:        []interface{}{1, 3, 5},
			expectedCount: 3,
			expectedIDs:   []uint{1, 3, 5},
		},
		{
			name:          "include specific ages",
			column:        "age",
			values:        []interface{}{25, 30},
			expectedCount: 2,
			expectedIDs:   []uint{1, 2},
		},
		{
			name:          "empty values list",
			column:        "id",
			values:        []interface{}{},
			expectedCount: 0, // No records
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []TestModel

			db.Scopes(WhereInScope(tt.column, tt.values)).Find(&results)

			if int64(len(results)) != tt.expectedCount {
				t.Errorf("WhereInScope() returned %d records, want %d", len(results), tt.expectedCount)
			}

			if len(tt.expectedIDs) > 0 {
				actualIDs := make([]uint, len(results))
				for i, result := range results {
					actualIDs[i] = result.ID
				}

				for _, expectedID := range tt.expectedIDs {
					found := false
					for _, actualID := range actualIDs {
						if actualID == expectedID {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("WhereInScope() missing expected ID %d in results", expectedID)
					}
				}
			}
		})
	}
}

func TestWhereIsScope(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name          string
		column        string
		value         interface{}
		expectedCount int64
	}{
		{
			name:          "where active is true",
			column:        "active",
			value:         true,
			expectedCount: 3, // Alice, Charlie, Eve
		},
		{
			name:          "where active is false",
			column:        "active",
			value:         false,
			expectedCount: 2, // Bob, David
		},
		{
			name:          "where age is 25",
			column:        "age",
			value:         25,
			expectedCount: 1, // Alice
		},
		{
			name:          "where name is Alice",
			column:        "name",
			value:         "Alice",
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []TestModel

			db.Scopes(WhereIsScope(tt.column, tt.value)).Find(&results)

			if int64(len(results)) != tt.expectedCount {
				t.Errorf("WhereIsScope() returned %d records, want %d", len(results), tt.expectedCount)
			}
		})
	}
}

func TestWhereLikeScope(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name          string
		column        string
		value         string
		expectedCount int64
	}{
		{
			name:          "name contains 'a'",
			column:        "name",
			value:         "%a%",
			expectedCount: 3, // Alice, Charlie, David
		},
		{
			name:          "email ends with 'example.com'",
			column:        "email",
			value:         "%example.com",
			expectedCount: 5, // All emails
		},
		{
			name:          "name starts with 'A'",
			column:        "name",
			value:         "A%",
			expectedCount: 1, // Alice
		},
		{
			name:          "no matches",
			column:        "name",
			value:         "%xyz%",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []TestModel

			db.Scopes(WhereLikeScope(tt.column, tt.value)).Find(&results)

			if int64(len(results)) != tt.expectedCount {
				t.Errorf("WhereLikeScope() returned %d records, want %d", len(results), tt.expectedCount)
			}
		})
	}
}

func TestWhereBetweenScope(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name          string
		column        string
		min           interface{}
		max           interface{}
		expectedCount int64
	}{
		{
			name:          "age between 25 and 35",
			column:        "age",
			min:           25,
			max:           35,
			expectedCount: 4, // Alice(25), Bob(30), Charlie(35), Eve(28)
		},
		{
			name:          "age between 30 and 40",
			column:        "age",
			min:           30,
			max:           40,
			expectedCount: 3, // Bob(30), Charlie(35), David(40)
		},
		{
			name:          "no matches in range",
			column:        "age",
			min:           50,
			max:           60,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []TestModel

			db.Scopes(WhereBetweenScope(tt.column, tt.min, tt.max)).Find(&results)

			if int64(len(results)) != tt.expectedCount {
				t.Errorf("WhereBetweenScope() returned %d records, want %d", len(results), tt.expectedCount)
			}
		})
	}
}

func TestWhereIsNullScope(t *testing.T) {
	db := setupTestDB(t)

	// Add a record with null email for testing
	db.Create(&TestModel{Name: "Frank", Age: 45, Active: true})

	tests := []struct {
		name          string
		column        string
		expectedCount int64
	}{
		{
			name:          "email is null",
			column:        "email",
			expectedCount: 1, // Frank has null email
		},
		{
			name:          "name is null (none)",
			column:        "name",
			expectedCount: 0, // All have names
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []TestModel

			db.Scopes(WhereIsNullScope(tt.column)).Find(&results)

			if int64(len(results)) != tt.expectedCount {
				t.Errorf("WhereIsNullScope() returned %d records, want %d", len(results), tt.expectedCount)
			}
		})
	}
}

func TestWhereIsNotNullScope(t *testing.T) {
	db := setupTestDB(t)

	// Add a record with null email for testing
	db.Create(&TestModel{Name: "Frank", Age: 45, Active: true})

	tests := []struct {
		name          string
		column        string
		expectedCount int64
	}{
		{
			name:          "email is not null",
			column:        "email",
			expectedCount: 5, // All except Frank
		},
		{
			name:          "name is not null",
			column:        "name",
			expectedCount: 6, // All have names
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []TestModel

			db.Scopes(WhereIsNotNullScope(tt.column)).Find(&results)

			if int64(len(results)) != tt.expectedCount {
				t.Errorf("WhereIsNotNullScope() returned %d records, want %d", len(results), tt.expectedCount)
			}
		})
	}
}

func TestJoinScope(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name        string
		joinClause  string
		expectError bool
	}{
		{
			name:        "valid join clause",
			joinClause:  "LEFT JOIN test_models t2 ON t2.id = test_models.id",
			expectError: false,
		},
		{
			name:        "empty join clause",
			joinClause:  "",
			expectError: false, // Should not error, just no effect
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []TestModel

			// Test that the scope can be applied without error
			err := db.Scopes(JoinScope(tt.joinClause)).Find(&results).Error

			if tt.expectError && err == nil {
				t.Error("JoinScope() expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("JoinScope() unexpected error: %v", err)
			}
		})
	}
}

func TestPreloadScope(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name        string
		association string
		expectError bool
	}{
		{
			name:        "non-existent association",
			association: "Profile", // Non-existent association
			expectError: true,      // Should error for non-existent association
		},
		{
			name:        "empty association",
			association: "",
			expectError: true, // Should error for empty association
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []TestModel

			// Test that the scope can be applied
			err := db.Scopes(PreloadScope(tt.association)).Find(&results).Error

			if tt.expectError && err == nil {
				t.Error("PreloadScope() expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("PreloadScope() unexpected error: %v", err)
			}
		})
	}
}

func TestCombinedScopes(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name          string
		scopes        []func(*gorm.DB) *gorm.DB
		expectedCount int64
	}{
		{
			name: "limit and offset",
			scopes: []func(*gorm.DB) *gorm.DB{
				LimitScope(2),
				OffsetScope(1),
			},
			expectedCount: 2,
		},
		{
			name: "where and order",
			scopes: []func(*gorm.DB) *gorm.DB{
				WhereIsScope("active", true),
				OrderScope("age DESC"),
			},
			expectedCount: 3, // Alice, Charlie, Eve (active users)
		},
		{
			name: "complex combination",
			scopes: []func(*gorm.DB) *gorm.DB{
				WhereIsScope("active", true),
				WhereBetweenScope("age", 25, 35),
				OrderScope("age ASC"),
				LimitScope(2),
			},
			expectedCount: 2, // Alice(25), Eve(28) - first 2 active users aged 25-35
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []TestModel

			db.Scopes(tt.scopes...).Find(&results)

			if int64(len(results)) != tt.expectedCount {
				t.Errorf("Combined scopes returned %d records, want %d", len(results), tt.expectedCount)
			}
		})
	}
}

func TestScopeWithDryRun(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name  string
		scope func(*gorm.DB) *gorm.DB
	}{
		{
			name:  "LimitScope dry run",
			scope: LimitScope(10),
		},
		{
			name:  "WhereIsScope dry run",
			scope: WhereIsScope("active", true),
		},
		{
			name:  "OrderScope dry run",
			scope: OrderScope("age DESC"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test dry run to check SQL generation
			stmt := db.Session(&gorm.Session{DryRun: true}).Scopes(tt.scope).Find(&TestModel{}).Statement

			if stmt.SQL.String() == "" {
				t.Error("Scope should generate SQL in dry run mode")
			}

			// Verify that the scope doesn't cause errors
			if stmt.Error != nil {
				t.Errorf("Scope caused error in dry run: %v", stmt.Error)
			}
		})
	}
}
