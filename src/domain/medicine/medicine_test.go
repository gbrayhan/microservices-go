package medicine

import (
	"testing"
	"time"
)

func TestMedicine_Fields(t *testing.T) {
	medicine := Medicine{
		ID:          1,
		Name:        "Test Medicine",
		Description: "Test Description",
		EanCode:     "1234567890123",
		Laboratory:  "Test Lab",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if medicine.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", medicine.ID)
	}

	if medicine.Name != "Test Medicine" {
		t.Errorf("Expected Name to be 'Test Medicine', got %s", medicine.Name)
	}

	if medicine.Description != "Test Description" {
		t.Errorf("Expected Description to be 'Test Description', got %s", medicine.Description)
	}

	if medicine.EanCode != "1234567890123" {
		t.Errorf("Expected EanCode to be '1234567890123', got %s", medicine.EanCode)
	}

	if medicine.Laboratory != "Test Lab" {
		t.Errorf("Expected Laboratory to be 'Test Lab', got %s", medicine.Laboratory)
	}
}

func TestMedicine_TimeFields(t *testing.T) {
	now := time.Now()
	medicine := Medicine{
		CreatedAt: now,
		UpdatedAt: now,
	}

	if !medicine.CreatedAt.Equal(now) {
		t.Errorf("Expected CreatedAt to be %v, got %v", now, medicine.CreatedAt)
	}

	if !medicine.UpdatedAt.Equal(now) {
		t.Errorf("Expected UpdatedAt to be %v, got %v", now, medicine.UpdatedAt)
	}
}

func TestMedicine_ZeroValues(t *testing.T) {
	medicine := Medicine{}

	if medicine.ID != 0 {
		t.Errorf("Expected ID to be 0, got %d", medicine.ID)
	}

	if medicine.Name != "" {
		t.Errorf("Expected Name to be empty, got %s", medicine.Name)
	}

	if medicine.Description != "" {
		t.Errorf("Expected Description to be empty, got %s", medicine.Description)
	}

	if medicine.EanCode != "" {
		t.Errorf("Expected EanCode to be empty, got %s", medicine.EanCode)
	}

	if medicine.Laboratory != "" {
		t.Errorf("Expected Laboratory to be empty, got %s", medicine.Laboratory)
	}

	if !medicine.CreatedAt.IsZero() {
		t.Errorf("Expected CreatedAt to be zero, got %v", medicine.CreatedAt)
	}

	if !medicine.UpdatedAt.IsZero() {
		t.Errorf("Expected UpdatedAt to be zero, got %v", medicine.UpdatedAt)
	}
}

func TestDataMedicine_Fields(t *testing.T) {
	medicines := []Medicine{
		{ID: 1, Name: "Medicine 1"},
		{ID: 2, Name: "Medicine 2"},
	}

	dataMedicine := DataMedicine{
		Data:  &medicines,
		Total: 2,
	}

	if len(*dataMedicine.Data) != 2 {
		t.Errorf("Expected Data length to be 2, got %d", len(*dataMedicine.Data))
	}

	if dataMedicine.Total != 2 {
		t.Errorf("Expected Total to be 2, got %d", dataMedicine.Total)
	}

	if (*dataMedicine.Data)[0].ID != 1 {
		t.Errorf("Expected first medicine ID to be 1, got %d", (*dataMedicine.Data)[0].ID)
	}

	if (*dataMedicine.Data)[1].ID != 2 {
		t.Errorf("Expected second medicine ID to be 2, got %d", (*dataMedicine.Data)[1].ID)
	}
}

func TestDataMedicine_ZeroValues(t *testing.T) {
	dataMedicine := DataMedicine{}

	if dataMedicine.Data != nil {
		t.Errorf("Expected Data to be nil, got %v", dataMedicine.Data)
	}

	if dataMedicine.Total != 0 {
		t.Errorf("Expected Total to be 0, got %d", dataMedicine.Total)
	}
}
