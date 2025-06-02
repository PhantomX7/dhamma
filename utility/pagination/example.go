package pagination

// Example usage:
func ExampleUsage() {
	// # String filter with like operator
	// GET /users?name=like:john

	// # Number filter with between operator
	// GET /users?age=between:20,30

	// # Enum filter with in operator
	// GET /users?status=in:active,inactive

	// # Combined filters
	// GET /users?name=like:john&age=gt:25&status=eq:active

	// Define filterable fields
	filterDef := NewFilterDefinition().
		AddFilter("name", FilterConfig{
			Field:     "name",
			Type:      FilterTypeString,
			Operators: []FilterOperator{OperatorLike, OperatorEquals},
		}).
		AddFilter("user_id", FilterConfig{
			Field:     "id",
			TableName: "User",
			Type:      FilterTypeID,
			Operators: []FilterOperator{OperatorIn, OperatorEquals},
		}).
		AddFilter("created_at", FilterConfig{
			Field:     "created_at",
			Type:      FilterTypeDateTime,
			Operators: []FilterOperator{OperatorBetween},
		}).
		AddFilter("status", FilterConfig{
			Field:      "status",
			Type:       FilterTypeEnum,
			EnumValues: []string{"active", "inactive", "pending"},
		})

	// Request conditions
	conditions := map[string][]string{
		"limit":      {"20"},
		"offset":     {"0"},
		"sort":       {"created_at desc"},
		"name":       {"John"},
		"user_id":    {"1,2,3"},
		"created_at": {"2023-01-01,2023-12-31"},
		"status":     {"active"},
	}

	// Create pagination
	pagination := NewPagination(conditions, filterDef, PaginationOptions{
		DefaultLimit: 20,
		MaxLimit:     100,
		DefaultOrder: "id desc",
	})

	// Build scopes
	scopeBuilder := NewScopeBuilder(pagination)
	filterScopes, metaScopes := scopeBuilder.Build()

	type YourModel struct {
		ID   int
		Name string
	}
	// Use in GORM query (example - requires actual database connection)
	// var db *gorm.DB
	// var results []YourModel
	// db.Scopes(metaScopes...).Scopes(filterScopes...).Find(&results)

	// For demonstration purposes, just verify scopes are created
	_ = filterScopes
	_ = metaScopes
}
