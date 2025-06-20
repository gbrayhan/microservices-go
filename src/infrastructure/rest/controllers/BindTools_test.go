package controllers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupGinContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestBindJSON(t *testing.T) {
	// Test valid JSON
	validJSON := `{"name": "test", "email": "test@example.com"}`

	c, _ := setupGinContext()
	c.Request = httptest.NewRequest("POST", "/test", bytes.NewBufferString(validJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	var request map[string]string
	err := BindJSON(c, &request)

	assert.NoError(t, err)
	assert.Equal(t, "test", request["name"])
	assert.Equal(t, "test@example.com", request["email"])

	// Test invalid JSON
	invalidJSON := `{"name": "test", "email": "test@example.com"`

	c, _ = setupGinContext()
	c.Request = httptest.NewRequest("POST", "/test", bytes.NewBufferString(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	err = BindJSON(c, &request)

	assert.Error(t, err)

	// Test empty body
	c, _ = setupGinContext()
	c.Request = httptest.NewRequest("POST", "/test", bytes.NewBufferString(""))
	c.Request.Header.Set("Content-Type", "application/json")

	err = BindJSON(c, &request)

	assert.Error(t, err)
}

func TestBindJSONMap(t *testing.T) {
	// Test valid JSON
	validJSON := `{"name": "test", "email": "test@example.com", "age": 25}`

	c, _ := setupGinContext()
	c.Request = httptest.NewRequest("POST", "/test", bytes.NewBufferString(validJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	var request map[string]any
	err := BindJSONMap(c, &request)

	assert.NoError(t, err)
	assert.Equal(t, "test", request["name"])
	assert.Equal(t, "test@example.com", request["email"])
	assert.Equal(t, float64(25), request["age"]) // JSON numbers are unmarshaled as float64

	// Test invalid JSON
	invalidJSON := `{"name": "test", "email": "test@example.com"`

	c, _ = setupGinContext()
	c.Request = httptest.NewRequest("POST", "/test", bytes.NewBufferString(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	err = BindJSONMap(c, &request)

	assert.Error(t, err)

	// Test empty body
	c, _ = setupGinContext()
	c.Request = httptest.NewRequest("POST", "/test", bytes.NewBufferString(""))
	c.Request.Header.Set("Content-Type", "application/json")

	err = BindJSONMap(c, &request)

	assert.Error(t, err)
}

func TestPaginationValues(t *testing.T) {
	// Test case 1: Normal pagination
	numPages, nextCursor, prevCursor := PaginationValues(10, 2, 25)
	assert.Equal(t, int64(3), numPages)   // (25 + 10 - 1) / 10 = 34 / 10 = 3
	assert.Equal(t, int64(3), nextCursor) // 2 + 1 = 3
	assert.Equal(t, int64(1), prevCursor) // 2 - 1 = 1

	// Test case 2: First page
	numPages, nextCursor, prevCursor = PaginationValues(10, 1, 25)
	assert.Equal(t, int64(3), numPages)   // (25 + 10 - 1) / 10 = 34 / 10 = 3
	assert.Equal(t, int64(2), nextCursor) // 1 + 1 = 2
	assert.Equal(t, int64(0), prevCursor) // 1 - 1 = 0 (but should be 0 for first page)

	// Test case 3: Last page
	numPages, nextCursor, prevCursor = PaginationValues(10, 3, 25)
	assert.Equal(t, int64(3), numPages)   // (25 + 10 - 1) / 10 = 34 / 10 = 3
	assert.Equal(t, int64(0), nextCursor) // 3 >= 3, so no next page
	assert.Equal(t, int64(2), prevCursor) // 3 - 1 = 2

	// Test case 4: Single page
	numPages, nextCursor, prevCursor = PaginationValues(10, 1, 5)
	assert.Equal(t, int64(1), numPages)   // (5 + 10 - 1) / 10 = 14 / 10 = 1
	assert.Equal(t, int64(0), nextCursor) // 1 >= 1, so no next page
	assert.Equal(t, int64(0), prevCursor) // 1 - 1 = 0

	// Test case 5: Empty result
	numPages, nextCursor, prevCursor = PaginationValues(10, 1, 0)
	assert.Equal(t, int64(0), numPages)   // (0 + 10 - 1) / 10 = 9 / 10 = 0
	assert.Equal(t, int64(0), nextCursor) // 1 >= 0, so no next page
	assert.Equal(t, int64(0), prevCursor) // 1 - 1 = 0

	// Test case 6: Large numbers
	numPages, nextCursor, prevCursor = PaginationValues(100, 5, 1000)
	assert.Equal(t, int64(10), numPages)  // (1000 + 100 - 1) / 100 = 1099 / 100 = 10
	assert.Equal(t, int64(6), nextCursor) // 5 + 1 = 6
	assert.Equal(t, int64(4), prevCursor) // 5 - 1 = 4
}

func TestMessageResponse(t *testing.T) {
	message := MessageResponse{
		Message: "Test message",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(message)
	assert.NoError(t, err)

	var unmarshaled MessageResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, message.Message, unmarshaled.Message)
}

func TestSortByDataRequest(t *testing.T) {
	sortRequest := SortByDataRequest{
		Field:     "name",
		Direction: "asc",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(sortRequest)
	assert.NoError(t, err)

	var unmarshaled SortByDataRequest
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, sortRequest.Field, unmarshaled.Field)
	assert.Equal(t, sortRequest.Direction, unmarshaled.Direction)
}

func TestFieldDateRangeDataRequest(t *testing.T) {
	dateRangeRequest := FieldDateRangeDataRequest{
		Field:     "created_at",
		StartDate: "2023-01-01",
		EndDate:   "2023-12-31",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(dateRangeRequest)
	assert.NoError(t, err)

	var unmarshaled FieldDateRangeDataRequest
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, dateRangeRequest.Field, unmarshaled.Field)
	assert.Equal(t, dateRangeRequest.StartDate, unmarshaled.StartDate)
	assert.Equal(t, dateRangeRequest.EndDate, unmarshaled.EndDate)
}
