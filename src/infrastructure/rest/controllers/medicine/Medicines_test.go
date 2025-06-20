package medicine

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gbrayhan/microservices-go/src/domain"
	medicineDomain "github.com/gbrayhan/microservices-go/src/domain/medicine"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

// MockMedicineService implements IMedicineService for testing
type MockMedicineService struct {
	createFunc  func(*medicineDomain.Medicine) (*medicineDomain.Medicine, error)
	getAllFunc  func() ([]*medicineDomain.Medicine, error)
	getDataFunc func(int64, int64, string, string, map[string][]string, string, []domain.DateRangeFilter) (*medicineDomain.DataMedicine, error)
	getByIDFunc func(int) (*medicineDomain.Medicine, error)
	updateFunc  func(int, map[string]any) (*medicineDomain.Medicine, error)
	deleteFunc  func(int) error
}

func (m *MockMedicineService) Create(medicine *medicineDomain.Medicine) (*medicineDomain.Medicine, error) {
	if m.createFunc != nil {
		return m.createFunc(medicine)
	}
	return nil, nil
}

func (m *MockMedicineService) GetAll() (*[]medicineDomain.Medicine, error) {
	if m.getAllFunc != nil {
		result, err := m.getAllFunc()
		if result == nil {
			return nil, err
		}
		// Convert []*medicineDomain.Medicine to []medicineDomain.Medicine
		slice := make([]medicineDomain.Medicine, len(result))
		for i, v := range result {
			slice[i] = *v
		}
		return &slice, err
	}
	return nil, nil
}

func (m *MockMedicineService) GetData(page, limit int64, sortBy, sortDirection string, filters map[string][]string, globalSearch string, dateRangeFilters []domain.DateRangeFilter) (*medicineDomain.DataMedicine, error) {
	if m.getDataFunc != nil {
		return m.getDataFunc(page, limit, sortBy, sortDirection, filters, globalSearch, dateRangeFilters)
	}
	return nil, nil
}

func (m *MockMedicineService) GetByID(id int) (*medicineDomain.Medicine, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id)
	}
	return nil, nil
}

func (m *MockMedicineService) Update(id int, updates map[string]any) (*medicineDomain.Medicine, error) {
	if m.updateFunc != nil {
		return m.updateFunc(id, updates)
	}
	return nil, nil
}

func (m *MockMedicineService) Delete(id int) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	return nil
}

func (m *MockMedicineService) GetByMap(medicineMap map[string]any) (*medicineDomain.Medicine, error) {
	return nil, nil
}

func setupLogger(t *testing.T) *logger.Logger {
	loggerInstance, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	return loggerInstance
}

func TestNewMedicineController(t *testing.T) {
	mockService := &MockMedicineService{}
	logger := setupLogger(t)
	controller := NewMedicineController(mockService, logger)

	if controller == nil {
		t.Error("Expected NewMedicineController to return a non-nil controller")
	}
}

func TestController_NewMedicine_Success(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock service
	mockService := &MockMedicineService{
		createFunc: func(medicine *medicineDomain.Medicine) (*medicineDomain.Medicine, error) {
			medicine.ID = 1
			return medicine, nil
		},
	}

	// Create controller
	logger := setupLogger(t)
	controller := NewMedicineController(mockService, logger)

	// Create test request
	request := NewMedicineRequest{
		Name:        "Test Medicine",
		Description: "Test Description",
		Laboratory:  "Test Lab",
		EanCode:     "1234567890123",
	}

	requestBody, _ := json.Marshal(request)

	// Create HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/medicines", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the method
	controller.NewMedicine(c)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestController_NewMedicine_InvalidRequest(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock service
	mockService := &MockMedicineService{}

	// Create controller
	logger := setupLogger(t)
	controller := NewMedicineController(mockService, logger)

	// Create invalid request
	requestBody := []byte(`{"name": "Test"}`) // Missing required fields

	// Create HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/medicines", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the method
	controller.NewMedicine(c)

	// Check that an error was added to the context
	if len(c.Errors) == 0 {
		t.Error("Expected error to be added to context")
	}
}

func TestController_GetAllMedicines_Success(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock service
	mockService := &MockMedicineService{
		getAllFunc: func() ([]*medicineDomain.Medicine, error) {
			return []*medicineDomain.Medicine{
				{
					ID:          1,
					Name:        "Medicine 1",
					Description: "Description 1",
					Laboratory:  "Lab 1",
					EanCode:     "1234567890123",
				},
			}, nil
		},
	}

	// Create controller
	logger := setupLogger(t)
	controller := NewMedicineController(mockService, logger)

	// Create HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/medicines", nil)

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the method
	controller.GetAllMedicines(c)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestController_GetMedicinesByID_Success(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock service
	mockService := &MockMedicineService{
		getByIDFunc: func(id int) (*medicineDomain.Medicine, error) {
			return &medicineDomain.Medicine{
				ID:          id,
				Name:        "Test Medicine",
				Description: "Test Description",
				Laboratory:  "Test Lab",
				EanCode:     "1234567890123",
			}, nil
		},
	}

	// Create controller
	logger := setupLogger(t)
	controller := NewMedicineController(mockService, logger)

	// Create HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/medicines/1", nil)

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Call the method
	controller.GetMedicinesByID(c)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestController_GetMedicinesByID_InvalidID(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock service
	mockService := &MockMedicineService{}

	// Create controller
	logger := setupLogger(t)
	controller := NewMedicineController(mockService, logger)

	// Create HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/medicines/invalid", nil)

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	// Call the method
	controller.GetMedicinesByID(c)

	// Check that an error was added to the context
	if len(c.Errors) == 0 {
		t.Error("Expected error to be added to context")
	}
}

func TestController_UpdateMedicine_Success(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock service
	mockService := &MockMedicineService{
		updateFunc: func(id int, updates map[string]any) (*medicineDomain.Medicine, error) {
			return &medicineDomain.Medicine{
				ID:          id,
				Name:        "Updated Medicine",
				Description: "Updated Description",
				Laboratory:  "Updated Lab",
				EanCode:     "1234567890123",
			}, nil
		},
	}

	// Create controller
	logger := setupLogger(t)
	controller := NewMedicineController(mockService, logger)

	// Create test request
	request := map[string]any{
		"name":        "Updated Medicine",
		"description": "Updated Description",
	}

	requestBody, _ := json.Marshal(request)

	// Create HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/medicines/1", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Call the method
	controller.UpdateMedicine(c)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestController_DeleteMedicine_Success(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock service
	mockService := &MockMedicineService{
		deleteFunc: func(id int) error {
			return nil
		},
	}

	// Create controller
	logger := setupLogger(t)
	controller := NewMedicineController(mockService, logger)

	// Create HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/medicines/1", nil)

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Call the method
	controller.DeleteMedicine(c)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestController_DeleteMedicine_InvalidID(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock service
	mockService := &MockMedicineService{}

	// Create controller
	logger := setupLogger(t)
	controller := NewMedicineController(mockService, logger)

	// Create HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/medicines/invalid", nil)

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	// Call the method
	controller.DeleteMedicine(c)

	// Check that an error was added to the context
	if len(c.Errors) == 0 {
		t.Error("Expected error to be added to context")
	}
}
