package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Test with default values
	config := LoadConfig()

	if config.Server.Port == "" {
		t.Error("Expected Server.Port to have a default value")
	}

	if config.Database.Host == "" {
		t.Error("Expected Database.Host to have a default value")
	}

	if config.Database.Port == "" {
		t.Error("Expected Database.Port to have a default value")
	}

	if config.Database.DBName == "" {
		t.Error("Expected Database.DBName to have a default value")
	}

	if config.Database.User == "" {
		t.Error("Expected Database.User to have a default value")
	}

	if config.Database.Password == "" {
		t.Error("Expected Database.Password to have a default value")
	}

	if config.JWT.AccessSecret == "" {
		t.Error("Expected JWT.AccessSecret to have a default value")
	}

	if config.JWT.RefreshSecret == "" {
		t.Error("Expected JWT.RefreshSecret to have a default value")
	}
}

func TestLoadConfig_WithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("JWT_ACCESS_SECRET", "test_access_secret")
	os.Setenv("JWT_REFRESH_SECRET", "test_refresh_secret")

	// Clean up after test
	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("JWT_ACCESS_SECRET")
		os.Unsetenv("JWT_REFRESH_SECRET")
	}()

	config := LoadConfig()

	if config.Server.Port != "8080" {
		t.Errorf("Expected Server.Port to be '8080', got %s", config.Server.Port)
	}

	if config.Database.Host != "localhost" {
		t.Errorf("Expected Database.Host to be 'localhost', got %s", config.Database.Host)
	}

	if config.Database.Port != "5432" {
		t.Errorf("Expected Database.Port to be '5432', got %s", config.Database.Port)
	}

	if config.Database.DBName != "testdb" {
		t.Errorf("Expected Database.DBName to be 'testdb', got %s", config.Database.DBName)
	}

	if config.Database.User != "testuser" {
		t.Errorf("Expected Database.User to be 'testuser', got %s", config.Database.User)
	}

	if config.Database.Password != "testpass" {
		t.Errorf("Expected Database.Password to be 'testpass', got %s", config.Database.Password)
	}

	if config.JWT.AccessSecret != "test_access_secret" {
		t.Errorf("Expected JWT.AccessSecret to be 'test_access_secret', got %s", config.JWT.AccessSecret)
	}

	if config.JWT.RefreshSecret != "test_refresh_secret" {
		t.Errorf("Expected JWT.RefreshSecret to be 'test_refresh_secret', got %s", config.JWT.RefreshSecret)
	}
}

func TestLoadTestConfig(t *testing.T) {
	config := LoadTestConfig()

	if config.Server.Port != "8081" {
		t.Errorf("Expected Server.Port to be '8081', got %s", config.Server.Port)
	}

	if config.Database.Host != "localhost" {
		t.Errorf("Expected Database.Host to be 'localhost', got %s", config.Database.Host)
	}

	if config.Database.DBName != "test_db" {
		t.Errorf("Expected Database.DBName to be 'test_db', got %s", config.Database.DBName)
	}

	if config.Database.User != "test_user" {
		t.Errorf("Expected Database.User to be 'test_user', got %s", config.Database.User)
	}

	if config.Database.Password != "test_password" {
		t.Errorf("Expected Database.Password to be 'test_password', got %s", config.Database.Password)
	}

	if config.JWT.AccessSecret != "test_access_secret" {
		t.Errorf("Expected JWT.AccessSecret to be 'test_access_secret', got %s", config.JWT.AccessSecret)
	}

	if config.JWT.RefreshSecret != "test_refresh_secret" {
		t.Errorf("Expected JWT.RefreshSecret to be 'test_refresh_secret', got %s", config.JWT.RefreshSecret)
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	// Test with environment variable set
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")

	result := getEnvOrDefault("TEST_KEY", "default_value")
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got %s", result)
	}

	// Test with environment variable not set
	result = getEnvOrDefault("NON_EXISTENT_KEY", "default_value")
	if result != "default_value" {
		t.Errorf("Expected 'default_value', got %s", result)
	}
}

func TestGetEnvAsInt64OrDefault(t *testing.T) {
	// Test with valid integer environment variable
	os.Setenv("TEST_INT", "123")
	defer os.Unsetenv("TEST_INT")

	result := getEnvAsInt64OrDefault("TEST_INT", 456)
	if result != 123 {
		t.Errorf("Expected 123, got %d", result)
	}

	// Test with invalid integer environment variable
	os.Setenv("TEST_INVALID", "not_a_number")
	defer os.Unsetenv("TEST_INVALID")

	result = getEnvAsInt64OrDefault("TEST_INVALID", 456)
	if result != 456 {
		t.Errorf("Expected 456, got %d", result)
	}

	// Test with environment variable not set
	result = getEnvAsInt64OrDefault("NON_EXISTENT_INT", 789)
	if result != 789 {
		t.Errorf("Expected 789, got %d", result)
	}
}
