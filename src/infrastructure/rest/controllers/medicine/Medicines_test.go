package medicine

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainMedicine "github.com/gbrayhan/microservices-go/src/domain/medicine"
	"github.com/gin-gonic/gin"
)

// MockMedicineService implements IMedicineService for testing
type MockMedicineService struct {
	createFunc  func(*domainMedicine.Medicine) (*domainMedicine.Medicine, error)
	getAllFunc  func() ([]*domainMedicine.Medicine, error)
	getDataFunc func(int64, int64, string, string, map[string][]string, string, []domain.DateRangeFilter) (*domainMedicine.DataMedicine, error)
	getByIDFunc func(int) (*domainMedicine.Medicine, error)
	updateFunc  func(int, map[string]any) (*domainMedicine.Medicine, error)
	deleteFunc  func(int) error
}

func (m *MockMedicineService) Create(medicine *domainMedicine.Medicine) (*domainMedicine.Medicine, error) {
	if m.createFunc != nil {
		return m.createFunc(medicine)
	}
	return nil, nil
}

func (m *MockMedicineService) GetAll() (*[]domainMedicine.Medicine, error) {
	if m.getAllFunc != nil {
		result, err := m.getAllFunc()
		if result == nil {
			return nil, err
		}
		// Convert []*domainMedicine.Medicine to []domainMedicine.Medicine
		slice := make([]domainMedicine.Medicine, len(result))
		for i, v := range result {
			slice[i] = *v
		}
		return &slice, err
	}
	return nil, nil
}

func (m *MockMedicineService) GetData(page, limit int64, sortBy, sortDirection string, filters map[string][]string, globalSearch string, dateRangeFilters []domain.DateRangeFilter) (*domainMedicine.DataMedicine, error) {
	if m.getDataFunc != nil {
		return m.getDataFunc(page, limit, sortBy, sortDirection, filters, globalSearch, dateRangeFilters)
	}
	return nil, nil
}

func (m *MockMedicineService) GetByID(id int) (*domainMedicine.Medicine, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id)
	}
	return nil, nil
}

func (m *MockMedicineService) Update(id int, updates map[string]any) (*domainMedicine.Medicine, error) {
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

func (m *MockMedicineService) GetByMap(medicineMap map[string]any) (*domainMedicine.Medicine, error) {
	return nil, nil
}

func TestNewMedicineController(t *testing.T) {
	mockService := &MockMedicineService{}
	controller := NewMedicineController(mockService)

	if controller == nil {
		t.Error("Expected NewMedicineController to return a non-nil controller")
	}
}

func TestController_NewMedicine_Success(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock service
	mockService := &MockMedicineService{
		createFunc: func(medicine *domainMedicine.Medicine) (*domainMedicine.Medicine, error) {
			medicine.ID = 1
			return medicine, nil
		},
	}

	// Create controller
	controller := NewMedicineController(mockService)

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
	controller := NewMedicineController(mockService)

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
		getAllFunc: func() ([]*domainMedicine.Medicine, error) {
			return []*domainMedicine.Medicine{
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
	controller := NewMedicineController(mockService)

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
		getByIDFunc: func(id int) (*domainMedicine.Medicine, error) {
			return &domainMedicine.Medicine{
				ID:          id,
				Name:        "Test Medicine",
				Description: "Test Description",
				Laboratory:  "Test Lab",
				EanCode:     "1234567890123",
			}, nil
		},
	}

	// Create controller
	controller := NewMedicineController(mockService)

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
	controller := NewMedicineController(mockService)

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
		updateFunc: func(id int, updates map[string]any) (*domainMedicine.Medicine, error) {
			return &domainMedicine.Medicine{
				ID:          id,
				Name:        "Updated Medicine",
				Description: "Updated Description",
				Laboratory:  "Updated Lab",
				EanCode:     "1234567890123",
			}, nil
		},
	}

	// Create controller
	controller := NewMedicineController(mockService)

	// Create test request
	requestBody := []byte(`{"name": "Updated Medicine"}`)

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
	controller := NewMedicineController(mockService)

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
	controller := NewMedicineController(mockService)

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
