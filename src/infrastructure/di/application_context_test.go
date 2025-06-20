package di

import (
	"os"
	"testing"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainMedicine "github.com/gbrayhan/microservices-go/src/domain/medicine"
	domainUser "github.com/gbrayhan/microservices-go/src/domain/user"
	"github.com/gbrayhan/microservices-go/src/infrastructure/security"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories and services
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetAll() (*[]domainUser.User, error) {
	args := m.Called()
	return args.Get(0).(*[]domainUser.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(id int) (*domainUser.User, error) {
	args := m.Called(id)
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *domainUser.User) (*domainUser.User, error) {
	args := m.Called(user)
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserRepository) GetOneByMap(userMap map[string]interface{}) (*domainUser.User, error) {
	args := m.Called(userMap)
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

func (m *MockMedicineRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMedicineRepository) Update(id int, medicineMap map[string]any) (*domainMedicine.Medicine, error) {
	args := m.Called(id, medicineMap)
	return args.Get(0).(*domainMedicine.Medicine), args.Error(1)
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateJWTToken(userID int, tokenType string) (*security.AppToken, error) {
	args := m.Called(userID, tokenType)
	return args.Get(0).(*security.AppToken), args.Error(1)
}

func (m *MockJWTService) GetClaimsAndVerifyToken(tokenString string, tokenType string) (jwt.MapClaims, error) {
	args := m.Called(tokenString, tokenType)
	return args.Get(0).(jwt.MapClaims), args.Error(1)
}

func TestNewTestApplicationContext(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockMedicineRepo := &MockMedicineRepository{}
	mockJWTService := &MockJWTService{}

	appContext := NewTestApplicationContext(mockUserRepo, mockMedicineRepo, mockJWTService)

	assert.NotNil(t, appContext)
	assert.Equal(t, mockUserRepo, appContext.UserRepository)
	assert.Equal(t, mockMedicineRepo, appContext.MedicineRepository)
	assert.Equal(t, mockJWTService, appContext.JWTService)

	// Test that controllers are created
	assert.NotNil(t, appContext.AuthController)
	assert.NotNil(t, appContext.UserController)
	assert.NotNil(t, appContext.MedicineController)

	// Test that use cases are created
	assert.NotNil(t, appContext.AuthUseCase)
	assert.NotNil(t, appContext.UserUseCase)
	assert.NotNil(t, appContext.MedicineUseCase)
}

func TestSetupDependencies(t *testing.T) {
	// This test will fail in CI/CD without a real database connection
	// We'll test the error path by setting invalid environment variables
	originalHost := os.Getenv("DB_HOST")
	os.Setenv("DB_HOST", "")
	defer os.Setenv("DB_HOST", originalHost)

	appContext, err := SetupDependencies()

	assert.Error(t, err)
	assert.Nil(t, appContext)
}

func TestApplicationContextStructure(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockMedicineRepo := &MockMedicineRepository{}
	mockJWTService := &MockJWTService{}

	appContext := NewTestApplicationContext(mockUserRepo, mockMedicineRepo, mockJWTService)

	// Test that all fields are properly set
	assert.NotNil(t, appContext.AuthController)
	assert.NotNil(t, appContext.UserController)
	assert.NotNil(t, appContext.MedicineController)
	assert.NotNil(t, appContext.JWTService)
	assert.NotNil(t, appContext.UserRepository)
	assert.NotNil(t, appContext.MedicineRepository)
	assert.NotNil(t, appContext.AuthUseCase)
	assert.NotNil(t, appContext.UserUseCase)
	assert.NotNil(t, appContext.MedicineUseCase)

	// Test that DB is nil in test context (as expected)
	assert.Nil(t, appContext.DB)
}
