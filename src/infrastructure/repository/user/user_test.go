package user

import (
	"testing"

	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainUser "github.com/gbrayhan/microservices-go/src/domain/user"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDatabaseInterface is a mock implementation of DatabaseInterface for testing
type MockDatabaseInterface struct {
	mock.Mock
}

func (m *MockDatabaseInterface) Find(dest interface{}, conds ...interface{}) repository.DatabaseInterface {
	args := m.Called(dest, conds)
	return args.Get(0).(repository.DatabaseInterface)
}

func (m *MockDatabaseInterface) Create(value interface{}) repository.DatabaseInterface {
	args := m.Called(value)
	return args.Get(0).(repository.DatabaseInterface)
}

func (m *MockDatabaseInterface) Limit(limit int) repository.DatabaseInterface {
	args := m.Called(limit)
	return args.Get(0).(repository.DatabaseInterface)
}

func (m *MockDatabaseInterface) Where(query interface{}, args ...interface{}) repository.DatabaseInterface {
	mockArgs := m.Called(query, args)
	return mockArgs.Get(0).(repository.DatabaseInterface)
}

func (m *MockDatabaseInterface) First(dest interface{}, conds ...interface{}) repository.DatabaseInterface {
	args := m.Called(dest, conds)
	return args.Get(0).(repository.DatabaseInterface)
}

func (m *MockDatabaseInterface) Model(value interface{}) repository.DatabaseInterface {
	args := m.Called(value)
	return args.Get(0).(repository.DatabaseInterface)
}

func (m *MockDatabaseInterface) Select(query interface{}, args ...interface{}) repository.DatabaseInterface {
	mockArgs := m.Called(query, args)
	return mockArgs.Get(0).(repository.DatabaseInterface)
}

func (m *MockDatabaseInterface) Updates(values interface{}) repository.DatabaseInterface {
	args := m.Called(values)
	return args.Get(0).(repository.DatabaseInterface)
}

func (m *MockDatabaseInterface) Delete(value interface{}, conds ...interface{}) repository.DatabaseInterface {
	args := m.Called(value, conds)
	return args.Get(0).(repository.DatabaseInterface)
}

func (m *MockDatabaseInterface) Error() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDatabaseInterface) RowsAffected() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func TestNewUserRepository(t *testing.T) {
	mockDB := &MockDatabaseInterface{}
	repo := &Repository{DB: mockDB}

	assert.NotNil(t, repo)
	assert.IsType(t, &Repository{}, repo)
}

func TestRepository_GetAll_Success(t *testing.T) {
	mockDB := &MockDatabaseInterface{}
	repo := &Repository{DB: mockDB}

	// Mock successful database query
	mockDB.On("Find", mock.AnythingOfType("*[]user.User"), mock.Anything).Return(mockDB)
	mockDB.On("Error").Return(nil)

	result, err := repo.GetAll()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockDB.AssertExpectations(t)
}

func TestRepository_GetAll_Error(t *testing.T) {
	mockDB := &MockDatabaseInterface{}
	repo := &Repository{DB: mockDB}

	// Mock database error
	mockDB.On("Find", mock.AnythingOfType("*[]user.User"), mock.Anything).Return(mockDB)
	mockDB.On("Error").Return(gorm.ErrInvalidDB)

	result, err := repo.GetAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	mockDB.AssertExpectations(t)
}

func TestRepository_Create_Success(t *testing.T) {
	mockDB := &MockDatabaseInterface{}
	repo := &Repository{DB: mockDB}

	userDomain := &domainUser.User{
		UserName:     "testuser",
		Email:        "test@example.com",
		FirstName:    "Test",
		LastName:     "User",
		HashPassword: "hashedpassword",
	}

	// Mock successful database creation
	mockDB.On("Create", mock.AnythingOfType("*user.User")).Return(mockDB)
	mockDB.On("Error").Return(nil)

	result, err := repo.Create(userDomain)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockDB.AssertExpectations(t)
}

func TestRepository_Create_DuplicateError(t *testing.T) {
	mockDB := &MockDatabaseInterface{}
	repo := &Repository{DB: mockDB}

	userDomain := &domainUser.User{
		UserName:     "testuser",
		Email:        "test@example.com",
		FirstName:    "Test",
		LastName:     "User",
		HashPassword: "hashedpassword",
	}

	// Mock duplicate key error (MySQL error 1062)
	mockDB.On("Create", mock.AnythingOfType("*user.User")).Return(mockDB)
	mockDB.On("Error").Return(gorm.ErrDuplicatedKey)

	result, err := repo.Create(userDomain)

	assert.Error(t, err)
	assert.NotNil(t, result)
	// Should return ResourceAlreadyExists error
	assert.IsType(t, &domainErrors.AppError{}, err)
	mockDB.AssertExpectations(t)
}

func TestRepository_GetOneByMap_Success(t *testing.T) {
	mockDB := &MockDatabaseInterface{}
	repo := &Repository{DB: mockDB}

	userMap := map[string]interface{}{"email": "test@example.com"}

	// Mock successful database query
	mockDB.On("Limit", 1).Return(mockDB)
	mockDB.On("Where", "email = ?", mock.Anything).Return(mockDB)
	mockDB.On("Find", mock.Anything, mock.Anything).Return(mockDB)
	mockDB.On("Error").Return(nil)

	result, err := repo.GetOneByMap(userMap)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockDB.AssertExpectations(t)
}

func TestRepository_GetOneByMap_Error(t *testing.T) {
	mockDB := &MockDatabaseInterface{}
	repo := &Repository{DB: mockDB}

	userMap := map[string]interface{}{"email": "test@example.com"}

	// Mock database error
	mockDB.On("Limit", 1).Return(mockDB)
	mockDB.On("Where", "email = ?", mock.Anything).Return(mockDB)
	mockDB.On("Find", mock.Anything, mock.Anything).Return(mockDB)
	mockDB.On("Error").Return(gorm.ErrInvalidDB)

	result, err := repo.GetOneByMap(userMap)

	assert.Error(t, err)
	assert.NotNil(t, result)
	mockDB.AssertExpectations(t)
}

func TestRepository_GetByID_Success(t *testing.T) {
	mockDB := &MockDatabaseInterface{}
	repo := &Repository{DB: mockDB}

	// Mock successful database query
	mockDB.On("Where", "id = ?", mock.Anything).Return(mockDB)
	mockDB.On("First", mock.Anything, mock.Anything).Return(mockDB)
	mockDB.On("Error").Return(nil)

	result, err := repo.GetByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockDB.AssertExpectations(t)
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	mockDB := &MockDatabaseInterface{}
	repo := &Repository{DB: mockDB}

	// Mock record not found error
	mockDB.On("Where", "id = ?", mock.Anything).Return(mockDB)
	mockDB.On("First", mock.Anything, mock.Anything).Return(mockDB)
	mockDB.On("Error").Return(gorm.ErrRecordNotFound)

	result, err := repo.GetByID(1)

	assert.Error(t, err)
	assert.NotNil(t, result)
	// Should return NotFound error
	assert.IsType(t, &domainErrors.AppError{}, err)
	mockDB.AssertExpectations(t)
}

func TestRepository_GetByID_OtherError(t *testing.T) {
	mockDB := &MockDatabaseInterface{}
	repo := &Repository{DB: mockDB}

	// Mock other database error
	mockDB.On("Where", "id = ?", mock.Anything).Return(mockDB)
	mockDB.On("First", mock.Anything, mock.Anything).Return(mockDB)
	mockDB.On("Error").Return(gorm.ErrInvalidDB)

	result, err := repo.GetByID(1)

	assert.Error(t, err)
	assert.NotNil(t, result)
	// Should return UnknownError
	assert.IsType(t, &domainErrors.AppError{}, err)
	mockDB.AssertExpectations(t)
}

func TestRepository_Update_Success(t *testing.T) {
	mockDB := &MockDatabaseInterface{}
	repo := &Repository{DB: mockDB}

	userMap := map[string]interface{}{"first_name": "Updated"}

	// Mock successful database update
	mockDB.On("Model", mock.AnythingOfType("*user.User")).Return(mockDB)
	mockDB.On("Select", mock.Anything, mock.Anything).Return(mockDB)
	mockDB.On("Updates", userMap).Return(mockDB)
	mockDB.On("Error").Return(nil)
	mockDB.On("Where", "id = ?", mock.Anything).Return(mockDB)
	mockDB.On("First", mock.Anything, mock.Anything).Return(mockDB)

	result, err := repo.Update(1, userMap)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockDB.AssertExpectations(t)
}

func TestRepository_Update_DuplicateError(t *testing.T) {
	mockDB := &MockDatabaseInterface{}
	repo := &Repository{DB: mockDB}

	userMap := map[string]interface{}{"first_name": "Updated"}

	// Mock duplicate key error
	mockDB.On("Model", mock.AnythingOfType("*user.User")).Return(mockDB)
	mockDB.On("Select", mock.Anything, mock.Anything).Return(mockDB)
	mockDB.On("Updates", userMap).Return(mockDB)
	mockDB.On("Error").Return(gorm.ErrDuplicatedKey)

	result, err := repo.Update(1, userMap)

	assert.Error(t, err)
	assert.NotNil(t, result)
	// Should return ResourceAlreadyExists error
	assert.IsType(t, &domainErrors.AppError{}, err)
	mockDB.AssertExpectations(t)
}

func TestRepository_Delete_Success(t *testing.T) {
	mockDB := &MockDatabaseInterface{}
	repo := &Repository{DB: mockDB}

	// Mock successful deletion
	mockDB.On("Delete", mock.Anything, mock.Anything).Return(mockDB)
	mockDB.On("Error").Return(nil)
	mockDB.On("RowsAffected").Return(int64(1))

	err := repo.Delete(1)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestRepository_Delete_NotFound(t *testing.T) {
	mockDB := &MockDatabaseInterface{}
	repo := &Repository{DB: mockDB}

	// Mock deletion with no rows affected
	mockDB.On("Delete", mock.Anything, mock.Anything).Return(mockDB)
	mockDB.On("Error").Return(nil)
	mockDB.On("RowsAffected").Return(int64(0))

	err := repo.Delete(1)

	assert.Error(t, err)
	// Should return NotFound error
	assert.IsType(t, &domainErrors.AppError{}, err)
	mockDB.AssertExpectations(t)
}

func TestRepository_Delete_Error(t *testing.T) {
	mockDB := &MockDatabaseInterface{}
	repo := &Repository{DB: mockDB}

	// Mock database error
	mockDB.On("Delete", mock.Anything, mock.Anything).Return(mockDB)
	mockDB.On("Error").Return(gorm.ErrInvalidDB)

	err := repo.Delete(1)

	assert.Error(t, err)
	// Should return UnknownError
	assert.IsType(t, &domainErrors.AppError{}, err)
	mockDB.AssertExpectations(t)
}
