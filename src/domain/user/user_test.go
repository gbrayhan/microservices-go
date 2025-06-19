package user

import (
	"testing"
	"time"
)

func TestUser_Fields(t *testing.T) {
	user := User{
		ID:           1,
		UserName:     "testuser",
		Email:        "test@example.com",
		FirstName:    "Test",
		LastName:     "User",
		Status:       true,
		HashPassword: "hashedpassword",
		Password:     "password",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if user.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", user.ID)
	}

	if user.UserName != "testuser" {
		t.Errorf("Expected UserName to be 'testuser', got %s", user.UserName)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected Email to be 'test@example.com', got %s", user.Email)
	}

	if user.FirstName != "Test" {
		t.Errorf("Expected FirstName to be 'Test', got %s", user.FirstName)
	}

	if user.LastName != "User" {
		t.Errorf("Expected LastName to be 'User', got %s", user.LastName)
	}

	if !user.Status {
		t.Errorf("Expected Status to be true, got %t", user.Status)
	}

	if user.HashPassword != "hashedpassword" {
		t.Errorf("Expected HashPassword to be 'hashedpassword', got %s", user.HashPassword)
	}

	if user.Password != "password" {
		t.Errorf("Expected Password to be 'password', got %s", user.Password)
	}
}

func TestUser_TimeFields(t *testing.T) {
	now := time.Now()
	user := User{
		CreatedAt: now,
		UpdatedAt: now,
	}

	if !user.CreatedAt.Equal(now) {
		t.Errorf("Expected CreatedAt to be %v, got %v", now, user.CreatedAt)
	}

	if !user.UpdatedAt.Equal(now) {
		t.Errorf("Expected UpdatedAt to be %v, got %v", now, user.UpdatedAt)
	}
}

func TestUser_ZeroValues(t *testing.T) {
	user := User{}

	if user.ID != 0 {
		t.Errorf("Expected ID to be 0, got %d", user.ID)
	}

	if user.UserName != "" {
		t.Errorf("Expected UserName to be empty, got %s", user.UserName)
	}

	if user.Email != "" {
		t.Errorf("Expected Email to be empty, got %s", user.Email)
	}

	if user.FirstName != "" {
		t.Errorf("Expected FirstName to be empty, got %s", user.FirstName)
	}

	if user.LastName != "" {
		t.Errorf("Expected LastName to be empty, got %s", user.LastName)
	}

	if user.Status {
		t.Errorf("Expected Status to be false, got %t", user.Status)
	}

	if user.HashPassword != "" {
		t.Errorf("Expected HashPassword to be empty, got %s", user.HashPassword)
	}

	if user.Password != "" {
		t.Errorf("Expected Password to be empty, got %s", user.Password)
	}

	if !user.CreatedAt.IsZero() {
		t.Errorf("Expected CreatedAt to be zero, got %v", user.CreatedAt)
	}

	if !user.UpdatedAt.IsZero() {
		t.Errorf("Expected UpdatedAt to be zero, got %v", user.UpdatedAt)
	}
}
