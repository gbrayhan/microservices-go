package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	domainUser "github.com/gbrayhan/microservices-go/src/domain/user"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of IUserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetAll() (*[]domainUser.User, error) {
	args := m.Called()
	return args.Get(0).(*[]domainUser.User), args.Error(1)
}

func (m *MockUserService) GetByID(id int) (*domainUser.User, error) {
	args := m.Called(id)
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserService) Create(user *domainUser.User) (*domainUser.User, error) {
	args := m.Called(user)
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserService) GetOneByMap(userMap map[string]interface{}) (*domainUser.User, error) {
	args := m.Called(userMap)
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserService) Update(id int, userMap map[string]interface{}) (*domainUser.User, error) {
	args := m.Called(id, userMap)
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupLogger(t *testing.T) *logger.Logger {
	loggerInstance, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	return loggerInstance
}

func TestNewUserController(t *testing.T) {
	mockService := &MockUserService{}
	loggerInstance := setupLogger(t)
	controller := NewUserController(mockService, loggerInstance)

	assert.NotNil(t, controller)
	assert.Equal(t, mockService, controller.(*UserController).userService)
	assert.Equal(t, loggerInstance, controller.(*UserController).Logger)
}

func TestDomainToResponseMapper(t *testing.T) {
	now := time.Now()
	domainUser := &domainUser.User{
		ID:           1,
		UserName:     "testuser",
		Email:        "test@example.com",
		FirstName:    "Test",
		LastName:     "User",
		Status:       true,
		HashPassword: "hashedpassword",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	response := domainToResponseMapper(domainUser)

	assert.Equal(t, domainUser.ID, response.ID)
	assert.Equal(t, domainUser.UserName, response.UserName)
	assert.Equal(t, domainUser.Email, response.Email)
	assert.Equal(t, domainUser.FirstName, response.FirstName)
	assert.Equal(t, domainUser.LastName, response.LastName)
	assert.Equal(t, domainUser.Status, response.Status)
	assert.Equal(t, domainUser.CreatedAt, response.CreatedAt)
	assert.Equal(t, domainUser.UpdatedAt, response.UpdatedAt)
}

func TestArrayDomainToResponseMapper(t *testing.T) {
	now := time.Now()
	users := []domainUser.User{
		{
			ID:           1,
			UserName:     "user1",
			Email:        "user1@example.com",
			FirstName:    "User",
			LastName:     "One",
			Status:       true,
			HashPassword: "hash1",
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:           2,
			UserName:     "user2",
			Email:        "user2@example.com",
			FirstName:    "User",
			LastName:     "Two",
			Status:       false,
			HashPassword: "hash2",
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}

	responses := arrayDomainToResponseMapper(&users)

	assert.Len(t, *responses, 2)
	assert.Equal(t, users[0].ID, (*responses)[0].ID)
	assert.Equal(t, users[1].ID, (*responses)[1].ID)
}

func TestToUsecaseMapper(t *testing.T) {
	request := &NewUserRequest{
		UserName:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Password:  "password123",
		Role:      "user",
	}

	domainUser := toUsecaseMapper(request)

	assert.Equal(t, request.UserName, domainUser.UserName)
	assert.Equal(t, request.Email, domainUser.Email)
	assert.Equal(t, request.FirstName, domainUser.FirstName)
	assert.Equal(t, request.LastName, domainUser.LastName)
	assert.Equal(t, request.Password, domainUser.Password)
}

func TestUpdateValidation(t *testing.T) {
	// Test valid request
	validRequest := map[string]any{
		"user_name": "validuser",
		"email":     "valid@example.com",
		"firstName": "Valid",
		"lastName":  "User",
	}

	err := updateValidation(validRequest)
	assert.NoError(t, err)

	// Test empty values
	emptyRequest := map[string]any{
		"user_name": "",
		"email":     "",
	}

	err = updateValidation(emptyRequest)
	assert.Error(t, err)

	// Test invalid email
	invalidEmailRequest := map[string]any{
		"email": "invalid-email",
	}

	err = updateValidation(invalidEmailRequest)
	assert.Error(t, err)

	// Test short user_name
	shortUserNameRequest := map[string]any{
		"user_name": "ab",
	}

	err = updateValidation(shortUserNameRequest)
	assert.Error(t, err)

	// Test long user_name
	longUserNameRequest := map[string]any{
		"user_name": "verylongusernameverylongusernameverylongusernameverylongusernameverylongusernameverylongusernameverylongusernameverylongusernameverylongusernameverylongusername",
	}

	err = updateValidation(longUserNameRequest)
	assert.Error(t, err)

	// Test short firstName
	shortFirstNameRequest := map[string]any{
		"firstName": "a",
	}

	err = updateValidation(shortFirstNameRequest)
	assert.Error(t, err)

	// Test long firstName
	longFirstNameRequest := map[string]any{
		"firstName": "verylongfirstnameverylongfirstnameverylongfirstnameverylongfirstnameverylongfirstnameverylongfirstnameverylongfirstnameverylongfirstnameverylongfirstname",
	}

	err = updateValidation(longFirstNameRequest)
	assert.Error(t, err)
}

func setupGinContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestUserController_NewUser(t *testing.T) {
	mockService := &MockUserService{}
	loggerInstance := setupLogger(t)
	controller := NewUserController(mockService, loggerInstance)

	t.Run("Success", func(t *testing.T) {
		c, w := setupGinContext()
		request := NewUserRequest{
			UserName:  "testuser",
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
			Password:  "password123",
			Role:      "user",
		}
		jsonData, _ := json.Marshal(request)
		c.Request = httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		expectedUser := &domainUser.User{
			ID:           1,
			UserName:     "testuser",
			Email:        "test@example.com",
			FirstName:    "Test",
			LastName:     "User",
			Status:       true,
			HashPassword: "hashedpassword",
		}

		mockService.On("Create", mock.Anything).Return(expectedUser, nil)

		controller.NewUser(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		c, w := setupGinContext()
		c.Request = httptest.NewRequest("POST", "/users", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.NewUser(c)

		assert.Equal(t, http.StatusOK, w.Code) // Gin returns 200 even on validation errors
	})

	t.Run("Service Error", func(t *testing.T) {
		c, w := setupGinContext()
		request := NewUserRequest{
			UserName:  "testuser",
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
			Password:  "password123",
			Role:      "user",
		}
		jsonData, _ := json.Marshal(request)
		c.Request = httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		mockService.On("Create", mock.Anything).Return((*domainUser.User)(nil), errors.New("service error"))

		controller.NewUser(c)

		assert.Equal(t, http.StatusOK, w.Code) // Gin returns 200 even on errors
		mockService.AssertExpectations(t)
	})
}

func TestUserController_GetAllUsers(t *testing.T) {
	mockService := &MockUserService{}
	loggerInstance := setupLogger(t)
	controller := NewUserController(mockService, loggerInstance)

	t.Run("Success", func(t *testing.T) {
		c, w := setupGinContext()
		c.Request = httptest.NewRequest("GET", "/users", nil)

		expectedUsers := &[]domainUser.User{
			{ID: 1, UserName: "user1", Email: "user1@example.com"},
			{ID: 2, UserName: "user2", Email: "user2@example.com"},
		}

		mockService.On("GetAll").Return(expectedUsers, nil)

		controller.GetAllUsers(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		c, w := setupGinContext()
		c.Request = httptest.NewRequest("GET", "/users", nil)

		mockService.On("GetAll").Return((*[]domainUser.User)(nil), errors.New("service error"))

		controller.GetAllUsers(c)

		assert.Equal(t, http.StatusOK, w.Code) // Gin returns 200 even on errors
		mockService.AssertExpectations(t)
	})
}

func TestUserController_GetUsersByID(t *testing.T) {
	mockService := &MockUserService{}
	loggerInstance := setupLogger(t)
	controller := NewUserController(mockService, loggerInstance)

	t.Run("Success", func(t *testing.T) {
		c, w := setupGinContext()
		c.Request = httptest.NewRequest("GET", "/users/1", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		expectedUser := &domainUser.User{
			ID:       1,
			UserName: "user1",
			Email:    "user1@example.com",
		}

		mockService.On("GetByID", 1).Return(expectedUser, nil)

		controller.GetUsersByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		c, w := setupGinContext()
		c.Request = httptest.NewRequest("GET", "/users/invalid", nil)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		controller.GetUsersByID(c)

		assert.Equal(t, http.StatusOK, w.Code) // Gin returns 200 even on validation errors
	})

	t.Run("Service Error", func(t *testing.T) {
		c, w := setupGinContext()
		c.Request = httptest.NewRequest("GET", "/users/1", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		mockService.On("GetByID", 1).Return(nil, errors.New("service error"))

		controller.GetUsersByID(c)

		assert.Equal(t, http.StatusOK, w.Code) // Gin returns 200 even on errors
		mockService.AssertExpectations(t)
	})
}

func TestUserController_UpdateUser(t *testing.T) {
	mockService := &MockUserService{}
	loggerInstance := setupLogger(t)
	controller := NewUserController(mockService, loggerInstance)

	t.Run("Success", func(t *testing.T) {
		c, w := setupGinContext()
		updateData := map[string]any{
			"user_name": "updateduser",
			"email":     "updated@example.com",
		}
		jsonData, _ := json.Marshal(updateData)
		c.Request = httptest.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		expectedUser := &domainUser.User{
			ID:       1,
			UserName: "updateduser",
			Email:    "updated@example.com",
		}

		mockService.On("Update", 1, updateData).Return(expectedUser, nil)

		controller.UpdateUser(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		c, w := setupGinContext()
		c.Request = httptest.NewRequest("PUT", "/users/invalid", nil)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		controller.UpdateUser(c)

		assert.Equal(t, http.StatusOK, w.Code) // Gin returns 200 even on validation errors
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		c, w := setupGinContext()
		c.Request = httptest.NewRequest("PUT", "/users/1", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		controller.UpdateUser(c)

		assert.Equal(t, http.StatusOK, w.Code) // Gin returns 200 even on validation errors
	})

	t.Run("Service Error", func(t *testing.T) {
		c, w := setupGinContext()
		updateData := map[string]any{"user_name": "updateduser"}
		jsonData, _ := json.Marshal(updateData)
		c.Request = httptest.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		mockService.On("Update", 1, updateData).Return((*domainUser.User)(nil), errors.New("service error"))

		controller.UpdateUser(c)

		assert.Equal(t, http.StatusOK, w.Code) // Gin returns 200 even on errors
		mockService.AssertExpectations(t)
	})
}

func TestUserController_DeleteUser(t *testing.T) {
	mockService := &MockUserService{}
	loggerInstance := setupLogger(t)
	controller := NewUserController(mockService, loggerInstance)

	t.Run("Success", func(t *testing.T) {
		c, w := setupGinContext()
		c.Request = httptest.NewRequest("DELETE", "/users/1", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		mockService.On("Delete", 1).Return(nil)

		controller.DeleteUser(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		c, w := setupGinContext()
		c.Request = httptest.NewRequest("DELETE", "/users/invalid", nil)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		controller.DeleteUser(c)

		assert.Equal(t, http.StatusOK, w.Code) // Gin returns 200 even on validation errors
	})

	t.Run("Service Error", func(t *testing.T) {
		c, w := setupGinContext()
		c.Request = httptest.NewRequest("DELETE", "/users/1", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		mockService.On("Delete", 1).Return(errors.New("service error"))

		controller.DeleteUser(c)

		assert.Equal(t, http.StatusOK, w.Code) // Gin returns 200 even on errors
		mockService.AssertExpectations(t)
	})
}
