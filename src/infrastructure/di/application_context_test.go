package di

import (
	"testing"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainMedicine "github.com/gbrayhan/microservices-go/src/domain/medicine"
	domainUser "github.com/gbrayhan/microservices-go/src/domain/user"
	"github.com/gbrayhan/microservices-go/src/infrastructure/security"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations for testing
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetAll() (*[]domainUser.User, error) {
	args := m.Called()
	return args.Get(0).(*[]domainUser.User), args.Error(1)
}

func (m *MockUserRepository) Create(userDomain *domainUser.User) (*domainUser.User, error) {
	args := m.Called(userDomain)
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserRepository) GetOneByMap(userMap map[string]interface{}) (*domainUser.User, error) {
	args := m.Called(userMap)
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(id int) (*domainUser.User, error) {
	args := m.Called(id)
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserRepository) Update(id int, userMap map[string]interface{}) (*domainUser.User, error) {
	args := m.Called(id, userMap)
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockMedicineRepository struct {
	mock.Mock
}

func (m *MockMedicineRepository) GetAll() (*[]domainMedicine.Medicine, error) {
	args := m.Called()
	return args.Get(0).(*[]domainMedicine.Medicine), args.Error(1)
}

func (m *MockMedicineRepository) GetData(page int64, limit int64, sortBy string, sortDirection string, filters map[string][]string, searchText string, dateRangeFilters []domain.DateRangeFilter) (*domainMedicine.DataMedicine, error) {
	args := m.Called(page, limit, sortBy, sortDirection, filters, searchText, dateRangeFilters)
	return args.Get(0).(*domainMedicine.DataMedicine), args.Error(1)
}

func (m *MockMedicineRepository) GetByID(id int) (*domainMedicine.Medicine, error) {
	args := m.Called(id)
	return args.Get(0).(*domainMedicine.Medicine), args.Error(1)
}

func (m *MockMedicineRepository) Create(medicine *domainMedicine.Medicine) (*domainMedicine.Medicine, error) {
	args := m.Called(medicine)
	return args.Get(0).(*domainMedicine.Medicine), args.Error(1)
}

func (m *MockMedicineRepository) GetByMap(medicineMap map[string]any) (*domainMedicine.Medicine, error) {
	args := m.Called(medicineMap)
	return args.Get(0).(*domainMedicine.Medicine), args.Error(1)
}

func (m *MockMedicineRepository) Update(id int, medicineMap map[string]any) (*domainMedicine.Medicine, error) {
	args := m.Called(id, medicineMap)
	return args.Get(0).(*domainMedicine.Medicine), args.Error(1)
}

func (m *MockMedicineRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateJWTToken(userID int, tokenType string) (*security.AppToken, error) {
	args := m.Called(userID, tokenType)
	return args.Get(0).(*security.AppToken), args.Error(1)
}

func (m *MockJWTService) GetClaimsAndVerifyToken(token string, tokenType string) (jwt.MapClaims, error) {
	args := m.Called(token, tokenType)
	return args.Get(0).(jwt.MapClaims), args.Error(1)
}

func TestNewTestApplicationContext(t *testing.T) {
	// Create mock dependencies
	mockUserRepo := &MockUserRepository{}
	mockMedicineRepo := &MockMedicineRepository{}
	mockJWTService := &MockJWTService{}

	// Create test application context
	appContext := NewTestApplicationContext(mockUserRepo, mockMedicineRepo, mockJWTService)

	// Assert that all dependencies are properly set
	assert.NotNil(t, appContext)
	assert.Equal(t, mockJWTService, appContext.JWTService)
	assert.NotNil(t, appContext.AuthController)
	assert.NotNil(t, appContext.UserController)
	assert.NotNil(t, appContext.MedicineController)
}

func TestNewTestApplicationContext_ControllerDependencies(t *testing.T) {
	// Create mock dependencies
	mockUserRepo := &MockUserRepository{}
	mockMedicineRepo := &MockMedicineRepository{}
	mockJWTService := &MockJWTService{}

	// Create test application context
	appContext := NewTestApplicationContext(mockUserRepo, mockMedicineRepo, mockJWTService)

	// Test that controllers are properly initialized
	// This is an indirect test - we're verifying the structure is correct
	assert.NotNil(t, appContext.AuthController)
	assert.NotNil(t, appContext.UserController)
	assert.NotNil(t, appContext.MedicineController)
}

func TestNewTestApplicationContext_JWTServiceDependency(t *testing.T) {
	// Create mock dependencies
	mockUserRepo := &MockUserRepository{}
	mockMedicineRepo := &MockMedicineRepository{}
	mockJWTService := &MockJWTService{}

	// Create test application context
	appContext := NewTestApplicationContext(mockUserRepo, mockMedicineRepo, mockJWTService)

	// Test that JWT service is properly set
	assert.Equal(t, mockJWTService, appContext.JWTService)
}
