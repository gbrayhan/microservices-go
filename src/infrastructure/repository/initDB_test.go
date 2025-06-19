package repository

import (
	"os"
	"testing"
)

func TestGetEnvAsInt_WithValidEnvVar(t *testing.T) {
	// Set environment variable
	os.Setenv("TEST_INT_VAR", "42")
	defer os.Unsetenv("TEST_INT_VAR")

	result := getEnvAsInt("TEST_INT_VAR", 10)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestGetEnvAsInt_WithInvalidEnvVar(t *testing.T) {
	// Set invalid environment variable
	os.Setenv("TEST_INVALID_VAR", "not_a_number")
	defer os.Unsetenv("TEST_INVALID_VAR")

	result := getEnvAsInt("TEST_INVALID_VAR", 10)
	if result != 10 {
		t.Errorf("Expected default value 10, got %d", result)
	}
}

func TestGetEnvAsInt_WithMissingEnvVar(t *testing.T) {
	// Ensure environment variable is not set
	os.Unsetenv("TEST_MISSING_VAR")

	result := getEnvAsInt("TEST_MISSING_VAR", 25)
	if result != 25 {
		t.Errorf("Expected default value 25, got %d", result)
	}
}

func TestGetEnvAsInt_WithEmptyEnvVar(t *testing.T) {
	// Set empty environment variable
	os.Setenv("TEST_EMPTY_VAR", "")
	defer os.Unsetenv("TEST_EMPTY_VAR")

	result := getEnvAsInt("TEST_EMPTY_VAR", 15)
	if result != 15 {
		t.Errorf("Expected default value 15, got %d", result)
	}
}

func TestGetEnvAsInt_WithZeroValue(t *testing.T) {
	// Set environment variable to zero
	os.Setenv("TEST_ZERO_VAR", "0")
	defer os.Unsetenv("TEST_ZERO_VAR")

	result := getEnvAsInt("TEST_ZERO_VAR", 100)
	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

func TestGetEnvAsInt_WithNegativeValue(t *testing.T) {
	// Set environment variable to negative value
	os.Setenv("TEST_NEGATIVE_VAR", "-5")
	defer os.Unsetenv("TEST_NEGATIVE_VAR")

	result := getEnvAsInt("TEST_NEGATIVE_VAR", 100)
	if result != -5 {
		t.Errorf("Expected -5, got %d", result)
	}
}
