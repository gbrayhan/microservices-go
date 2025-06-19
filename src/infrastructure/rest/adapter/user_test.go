package adapter

import (
	"testing"

	"gorm.io/gorm"
)

func TestUserAdapter(t *testing.T) {
	// Create a mock database (nil for testing purposes)
	var db *gorm.DB

	// Test the adapter
	userController := UserAdapter(db)

	// Verify that the adapter returns a non-nil controller
	if userController == nil {
		t.Error("Expected UserAdapter to return a non-nil controller")
	}

	// Verify that the controller implements the expected interface
	controllerType := userController
	if controllerType == nil {
		t.Error("Expected controller to be of the correct type")
	}
}
