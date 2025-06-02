package utility

import (
	"testing"
)

func TestPointOf(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "string value",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "integer value",
			input:    42,
			expected: 42,
		},
		{
			name:     "boolean value",
			input:    true,
			expected: true,
		},
		{
			name:     "float value",
			input:    3.14,
			expected: 3.14,
		},
		{
			name:     "zero value string",
			input:    "",
			expected: "",
		},
		{
			name:     "zero value int",
			input:    0,
			expected: 0,
		},
		{
			name:     "zero value bool",
			input:    false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.input.(type) {
			case string:
				result := PointOf(v)
				if result == nil {
					t.Errorf("PointOf() returned nil for input %v", v)
					return
				}
				if *result != tt.expected {
					t.Errorf("PointOf() = %v, want %v", *result, tt.expected)
				}
			case int:
				result := PointOf(v)
				if result == nil {
					t.Errorf("PointOf() returned nil for input %v", v)
					return
				}
				if *result != tt.expected {
					t.Errorf("PointOf() = %v, want %v", *result, tt.expected)
				}
			case bool:
				result := PointOf(v)
				if result == nil {
					t.Errorf("PointOf() returned nil for input %v", v)
					return
				}
				if *result != tt.expected {
					t.Errorf("PointOf() = %v, want %v", *result, tt.expected)
				}
			case float64:
				result := PointOf(v)
				if result == nil {
					t.Errorf("PointOf() returned nil for input %v", v)
					return
				}
				if *result != tt.expected {
					t.Errorf("PointOf() = %v, want %v", *result, tt.expected)
				}
			}
		})
	}
}

func TestPointOf_PointerBehavior(t *testing.T) {
	// Test that the function returns a pointer to the value
	value := "test"
	ptr := PointOf(value)

	if ptr == nil {
		t.Error("PointOf() returned nil")
		return
	}

	// Verify it's actually a pointer
	if *ptr != value {
		t.Errorf("PointOf() = %v, want %v", *ptr, value)
	}

	// Verify modifying the original value doesn't affect the pointer
	originalValue := value
	value = "modified"
	if *ptr != originalValue {
		t.Errorf("Pointer value changed when original was modified: got %v, want %v", *ptr, originalValue)
	}
}

func TestPointOf_DifferentTypes(t *testing.T) {
	// Test with custom struct
	type TestStruct struct {
		Name string
		Age  int
	}

	testStruct := TestStruct{Name: "John", Age: 30}
	ptr := PointOf(testStruct)

	if ptr == nil {
		t.Error("PointOf() returned nil for struct")
		return
	}

	if ptr.Name != testStruct.Name || ptr.Age != testStruct.Age {
		t.Errorf("PointOf() struct = %+v, want %+v", *ptr, testStruct)
	}

	// Test with slice
	slice := []int{1, 2, 3}
	slicePtr := PointOf(slice)

	if slicePtr == nil {
		t.Error("PointOf() returned nil for slice")
		return
	}

	if len(*slicePtr) != len(slice) {
		t.Errorf("PointOf() slice length = %d, want %d", len(*slicePtr), len(slice))
	}

	for i, v := range slice {
		if (*slicePtr)[i] != v {
			t.Errorf("PointOf() slice[%d] = %d, want %d", i, (*slicePtr)[i], v)
		}
	}
}
