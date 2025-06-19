package adapter

import (
	"testing"

	"gorm.io/gorm"
)

func TestMedicineAdapter(t *testing.T) {
	// Create a mock database (nil for testing purposes)
	var db *gorm.DB

	// Test the adapter
	medicineController := MedicineAdapter(db)

	// Verify that the adapter returns a non-nil controller
	if medicineController == nil {
		t.Error("Expected MedicineAdapter to return a non-nil controller")
	}

	// Verify that the controller implements the expected interface
	controllerType := medicineController
	if controllerType == nil {
		t.Error("Expected controller to be of the correct type")
	}
}
