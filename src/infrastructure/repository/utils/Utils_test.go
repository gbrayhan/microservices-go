package utils

import (
	"reflect"
	"testing"

	"github.com/gbrayhan/microservices-go/src/domain"
)

func TestComplementSearch_NilDB(t *testing.T) {
	query, err := ComplementSearch(nil, "", "", 0, 0, nil, nil, "", nil, nil)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if query != nil {
		t.Errorf("Expected nil query, got %v", query)
	}
}

func TestUpdateFilterKeys(t *testing.T) {
	filters := map[string][]string{
		"name": {"test1", "test2"},
		"age":  {"25", "30"},
	}

	columnMapping := map[string]string{
		"name": "user_name",
		"age":  "user_age",
	}

	result := UpdateFilterKeys(filters, columnMapping)

	expected := map[string][]string{
		"user_name": {"test1", "test2"},
		"user_age":  {"25", "30"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestUpdateFilterKeys_NoMapping(t *testing.T) {
	filters := map[string][]string{
		"name": {"test1", "test2"},
	}

	columnMapping := map[string]string{}

	result := UpdateFilterKeys(filters, columnMapping)

	if !reflect.DeepEqual(result, filters) {
		t.Errorf("Expected %v, got %v", filters, result)
	}
}

func TestIsZeroValue_String(t *testing.T) {
	if !IsZeroValue("") {
		t.Error("Expected empty string to be zero value")
	}

	if IsZeroValue("test") {
		t.Error("Expected non-empty string to not be zero value")
	}
}

func TestIsZeroValue_Int(t *testing.T) {
	if !IsZeroValue(0) {
		t.Error("Expected 0 to be zero value")
	}

	if IsZeroValue(42) {
		t.Error("Expected non-zero int to not be zero value")
	}
}

func TestIsZeroValue_Bool(t *testing.T) {
	if !IsZeroValue(false) {
		t.Error("Expected false to be zero value")
	}

	if IsZeroValue(true) {
		t.Error("Expected true to not be zero value")
	}
}

func TestIsZeroValue_Slice(t *testing.T) {
	// Test with nil slice
	if !IsZeroValue([]string(nil)) {
		t.Error("Expected nil slice to be zero value")
	}

	// Test with empty slice - this is the tricky case
	// An empty slice is not a zero value in Go, it's a valid slice with length 0
	if IsZeroValue([]string{}) {
		t.Error("Expected empty slice to not be zero value")
	}

	if IsZeroValue([]string{"test"}) {
		t.Error("Expected non-empty slice to not be zero value")
	}
}

func TestIsZeroValue_Struct(t *testing.T) {
	type TestStruct struct {
		Name string
		Age  int
	}

	if !IsZeroValue(TestStruct{}) {
		t.Error("Expected empty struct to be zero value")
	}

	if IsZeroValue(TestStruct{Name: "test", Age: 25}) {
		t.Error("Expected non-empty struct to not be zero value")
	}
}

func TestApplyFilters_WithFilters(t *testing.T) {
	columnMapping := map[string]string{
		"name": "user_name",
	}

	filters := map[string][]string{
		"name": {"test1", "test2"},
	}

	applyFunc := ApplyFilters(columnMapping, filters, nil, "", nil)

	// This is a basic test to ensure the function returns a function
	if applyFunc == nil {
		t.Error("Expected ApplyFilters to return a function")
	}
}

func TestApplyFilters_WithDateRangeFilters(t *testing.T) {
	columnMapping := map[string]string{
		"created_at": "created_date",
	}

	dateRangeFilters := []domain.DateRangeFilter{
		{
			Field: "created_at",
			Start: "2023-01-01",
			End:   "2023-12-31",
		},
	}

	applyFunc := ApplyFilters(columnMapping, nil, dateRangeFilters, "", nil)

	if applyFunc == nil {
		t.Error("Expected ApplyFilters to return a function")
	}
}

func TestApplyFilters_WithSearchText(t *testing.T) {
	searchColumns := []string{"name", "description"}
	searchText := "test"

	applyFunc := ApplyFilters(nil, nil, nil, searchText, searchColumns)

	if applyFunc == nil {
		t.Error("Expected ApplyFilters to return a function")
	}
}
