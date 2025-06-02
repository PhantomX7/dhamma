package utility

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildResponseSuccess(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		data     any
		expected Response
	}{
		{
			name:    "success response with string data",
			message: "Operation successful",
			data:    "test data",
			expected: Response{
				Status:  true,
				Message: "Operation successful",
				Error:   nil,
				Data:    "test data",
			},
		},
		{
			name:    "success response with map data",
			message: "Data retrieved",
			data:    map[string]interface{}{"id": 1, "name": "test"},
			expected: Response{
				Status:  true,
				Message: "Data retrieved",
				Error:   nil,
				Data:    map[string]interface{}{"id": 1, "name": "test"},
			},
		},
		{
			name:    "success response with nil data",
			message: "Success",
			data:    nil,
			expected: Response{
				Status:  true,
				Message: "Success",
				Error:   nil,
				Data:    nil,
			},
		},
		{
			name:    "success response with empty message",
			message: "",
			data:    "data",
			expected: Response{
				Status:  true,
				Message: "",
				Error:   nil,
				Data:    "data",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildResponseSuccess(tt.message, tt.data)
			assert.Equal(t, tt.expected, result)
			assert.True(t, result.Status)
			assert.Equal(t, tt.message, result.Message)
			assert.Equal(t, tt.data, result.Data)
			assert.Nil(t, result.Error)
		})
	}
}

func TestBuildResponseFailed(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		err      any
		expected Response
	}{
		{
			name:    "failed response with string error",
			message: "Operation failed",
			err:     "validation error",
			expected: Response{
				Status:  false,
				Message: "Operation failed",
				Error:   "validation error",
				Data:    nil,
			},
		},
		{
			name:    "failed response with map error",
			message: "Validation failed",
			err:     map[string]string{"field": "required"},
			expected: Response{
				Status:  false,
				Message: "Validation failed",
				Error:   map[string]string{"field": "required"},
				Data:    nil,
			},
		},
		{
			name:    "failed response with nil error",
			message: "Unknown error",
			err:     nil,
			expected: Response{
				Status:  false,
				Message: "Unknown error",
				Error:   nil,
				Data:    nil,
			},
		},
		{
			name:    "failed response with empty message",
			message: "",
			err:     "error occurred",
			expected: Response{
				Status:  false,
				Message: "",
				Error:   "error occurred",
				Data:    nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildResponseFailed(tt.message, tt.err)
			assert.Equal(t, tt.expected, result)
			assert.False(t, result.Status)
			assert.Equal(t, tt.message, result.Message)
			assert.Equal(t, tt.err, result.Error)
			assert.Nil(t, result.Data)
		})
	}
}

func TestBuildPaginationResponseSuccess(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		data     any
		meta     PaginationMeta
		expected Response
	}{
		{
			name:    "pagination response with data",
			message: "Data retrieved successfully",
			data:    []map[string]interface{}{{"id": 1}, {"id": 2}},
			meta: PaginationMeta{
				Limit:  10,
				Offset: 0,
				Total:  100,
			},
			expected: Response{
				Status:  true,
				Message: "Data retrieved successfully",
				Error:   nil,
				Data:    []map[string]interface{}{{"id": 1}, {"id": 2}},
				Meta: PaginationMeta{
					Limit:  10,
					Offset: 0,
					Total:  100,
				},
			},
		},
		{
			name:    "pagination response with empty data",
			message: "No data found",
			data:    []interface{}{},
			meta: PaginationMeta{
				Limit:  10,
				Offset: 20,
				Total:  0,
			},
			expected: Response{
				Status:  true,
				Message: "No data found",
				Error:   nil,
				Data:    []interface{}{},
				Meta: PaginationMeta{
					Limit:  10,
					Offset: 20,
					Total:  0,
				},
			},
		},
		{
			name:    "pagination response with nil data",
			message: "Success",
			data:    nil,
			meta: PaginationMeta{
				Limit:  5,
				Offset: 0,
				Total:  50,
			},
			expected: Response{
				Status:  true,
				Message: "Success",
				Error:   nil,
				Data:    nil,
				Meta: PaginationMeta{
					Limit:  5,
					Offset: 0,
					Total:  50,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildPaginationResponseSuccess(tt.message, tt.data, tt.meta)
			assert.Equal(t, tt.expected, result)
			assert.True(t, result.Status)
			assert.Equal(t, tt.message, result.Message)
			assert.Equal(t, tt.data, result.Data)
			assert.Equal(t, tt.meta, result.Meta)
			assert.Nil(t, result.Error)
		})
	}
}

func TestPaginationMeta(t *testing.T) {
	tests := []struct {
		name     string
		meta     PaginationMeta
		expected map[string]interface{}
	}{
		{
			name: "valid pagination meta",
			meta: PaginationMeta{
				Limit:  10,
				Offset: 0,
				Total:  100,
			},
			expected: map[string]interface{}{
				"limit":  float64(10),
				"offset": float64(0),
				"total":  float64(100),
			},
		},
		{
			name: "zero values pagination meta",
			meta: PaginationMeta{
				Limit:  0,
				Offset: 0,
				Total:  0,
			},
			expected: map[string]interface{}{
				"limit":  float64(0),
				"offset": float64(0),
				"total":  float64(0),
			},
		},
		{
			name: "large values pagination meta",
			meta: PaginationMeta{
				Limit:  1000,
				Offset: 5000,
				Total:  999999,
			},
			expected: map[string]interface{}{
				"limit":  float64(1000),
				"offset": float64(5000),
				"total":  float64(999999),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.meta)
			require.NoError(t, err)

			// Test JSON unmarshaling
			var result map[string]interface{}
			err = json.Unmarshal(jsonData, &result)
			require.NoError(t, err)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestResponse(t *testing.T) {
	tests := []struct {
		name     string
		response Response
		expected map[string]interface{}
	}{
		{
			name: "complete response with all fields",
			response: Response{
				Status:  true,
				Message: "Success",
				Error:   nil,
				Data:    map[string]string{"key": "value"},
				Meta:    PaginationMeta{Limit: 10, Offset: 0, Total: 1},
			},
			expected: map[string]interface{}{
				"status":  true,
				"message": "Success",
				"data":    map[string]interface{}{"key": "value"},
				"meta": map[string]interface{}{
					"limit":  float64(10),
					"offset": float64(0),
					"total":  float64(1),
				},
			},
		},
		{
			name: "error response",
			response: Response{
				Status:  false,
				Message: "Error occurred",
				Error:   "validation failed",
				Data:    nil,
				Meta:    nil,
			},
			expected: map[string]interface{}{
				"status":  false,
				"message": "Error occurred",
				"error":   "validation failed",
			},
		},
		{
			name: "minimal success response",
			response: Response{
				Status:  true,
				Message: "OK",
			},
			expected: map[string]interface{}{
				"status":  true,
				"message": "OK",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.response)
			require.NoError(t, err)

			// Test JSON unmarshaling
			var result map[string]interface{}
			err = json.Unmarshal(jsonData, &result)
			require.NoError(t, err)

			// Check that omitempty works correctly
			for key, expectedValue := range tt.expected {
				assert.Equal(t, expectedValue, result[key], "Field %s should match", key)
			}

			// Check that omitted fields are not present
			if tt.response.Error == nil {
				_, exists := result["error"]
				assert.False(t, exists, "Error field should be omitted when nil")
			}
			if tt.response.Data == nil {
				_, exists := result["data"]
				assert.False(t, exists, "Data field should be omitted when nil")
			}
			if tt.response.Meta == nil {
				_, exists := result["meta"]
				assert.False(t, exists, "Meta field should be omitted when nil")
			}
		})
	}
}

func TestResponseJSONSerialization(t *testing.T) {
	t.Run("success response serialization", func(t *testing.T) {
		response := BuildResponseSuccess("Test message", map[string]string{"test": "data"})
		jsonData, err := json.Marshal(response)
		require.NoError(t, err)

		expectedJSON := `{"status":true,"message":"Test message","data":{"test":"data"}}`
		assert.JSONEq(t, expectedJSON, string(jsonData))
	})

	t.Run("failed response serialization", func(t *testing.T) {
		response := BuildResponseFailed("Error message", "error details")
		jsonData, err := json.Marshal(response)
		require.NoError(t, err)

		expectedJSON := `{"status":false,"message":"Error message","error":"error details"}`
		assert.JSONEq(t, expectedJSON, string(jsonData))
	})

	t.Run("pagination response serialization", func(t *testing.T) {
		meta := PaginationMeta{Limit: 10, Offset: 0, Total: 100}
		response := BuildPaginationResponseSuccess("Paginated data", []string{"item1", "item2"}, meta)
		jsonData, err := json.Marshal(response)
		require.NoError(t, err)

		expectedJSON := `{"status":true,"message":"Paginated data","data":["item1","item2"],"meta":{"limit":10,"offset":0,"total":100}}`
		assert.JSONEq(t, expectedJSON, string(jsonData))
	})
}

// Benchmark tests
func BenchmarkBuildResponseSuccess(b *testing.B) {
	data := map[string]interface{}{"id": 1, "name": "test"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BuildResponseSuccess("Success", data)
	}
}

func BenchmarkBuildResponseFailed(b *testing.B) {
	err := "validation error"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BuildResponseFailed("Failed", err)
	}
}

func BenchmarkBuildPaginationResponseSuccess(b *testing.B) {
	data := []map[string]interface{}{{"id": 1}, {"id": 2}}
	meta := PaginationMeta{Limit: 10, Offset: 0, Total: 100}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BuildPaginationResponseSuccess("Success", data, meta)
	}
}

func BenchmarkResponseJSONMarshal(b *testing.B) {
	response := BuildResponseSuccess("Test", map[string]string{"key": "value"})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(response)
	}
}
