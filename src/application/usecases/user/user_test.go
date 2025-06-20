package user

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gbrayhan/microservices-go/src/domain"
	userDomain "github.com/gbrayhan/microservices-go/src/domain/user"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
)

type mockUserService struct {
	getAllFn     func() (*[]userDomain.User, error)
	getByIDFn    func(id int) (*userDomain.User, error)
	getByEmailFn func(email string) (*userDomain.User, error)
	createFn     func(u *userDomain.User) (*userDomain.User, error)
	deleteFn     func(id int) error
	updateFn     func(id int, m map[string]interface{}) (*userDomain.User, error)
}

func (m *mockUserService) GetAll() (*[]userDomain.User, error) {
	return m.getAllFn()
}
func (m *mockUserService) GetByID(id int) (*userDomain.User, error) {
	return m.getByIDFn(id)
}
func (m *mockUserService) GetByEmail(email string) (*userDomain.User, error) {
	return m.getByEmailFn(email)
}
func (m *mockUserService) Create(newUser *userDomain.User) (*userDomain.User, error) {
	return m.createFn(newUser)
}
func (m *mockUserService) Delete(id int) error {
	return m.deleteFn(id)
}
func (m *mockUserService) Update(id int, userMap map[string]interface{}) (*userDomain.User, error) {
	return m.updateFn(id, userMap)
}
func (m *mockUserService) SearchPaginated(filters domain.DataFilters) (*userDomain.SearchResultUser, error) {
	return nil, nil
}
func (m *mockUserService) SearchByProperty(property string, searchText string) (*[]string, error) {
	return nil, nil
}

func setupLogger(t *testing.T) *logger.Logger {
	loggerInstance, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	return loggerInstance
}

func TestUserUseCase(t *testing.T) {

	mockRepo := &mockUserService{}
	logger := setupLogger(t)
	useCase := NewUserUseCase(mockRepo, logger)

	t.Run("Test GetAll", func(t *testing.T) {
		mockRepo.getAllFn = func() (*[]userDomain.User, error) {
			return &[]userDomain.User{{ID: 1}}, nil
		}
		us, err := useCase.GetAll()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(*us) != 1 {
			t.Error("expected 1 user from GetAll")
		}
	})

	t.Run("Test GetByID", func(t *testing.T) {
		mockRepo.getByIDFn = func(id int) (*userDomain.User, error) {
			if id == 999 {
				return nil, errors.New("not found")
			}
			return &userDomain.User{ID: id}, nil
		}
		_, err := useCase.GetByID(999)
		if err == nil {
			t.Error("expected error, got nil")
		}
		u, err := useCase.GetByID(10)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if u.ID != 10 {
			t.Errorf("expected user ID=10, got %d", u.ID)
		}
	})

	t.Run("Test GetByEmail", func(t *testing.T) {
		mockRepo.getByEmailFn = func(email string) (*userDomain.User, error) {
			if email == "notfound@example.com" {
				return nil, errors.New("not found")
			}
			return &userDomain.User{ID: 123, Email: email}, nil
		}
		_, err := useCase.GetByEmail("notfound@example.com")
		if err == nil {
			t.Error("expected error, got nil")
		}
		u, err := useCase.GetByEmail("test@example.com")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if u.ID != 123 {
			t.Errorf("expected user ID=123, got %d", u.ID)
		}
		if u.Email != "test@example.com" {
			t.Errorf("expected email=test@example.com, got %s", u.Email)
		}
	})

	t.Run("Test Create (OK)", func(t *testing.T) {
		mockRepo.createFn = func(newU *userDomain.User) (*userDomain.User, error) {
			if !newU.Status {
				t.Error("expected user.Status to be true")
			}
			if newU.HashPassword == "" {
				t.Error("expected user.HashPassword to be set")
			}
			if newU.Email == "" {
				return nil, errors.New("bad data")
			}
			newU.ID = 555
			return newU, nil
		}
		created, err := useCase.Create(&userDomain.User{Email: "test@mail.com", Password: "abc"})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if created.ID != 555 {
			t.Error("expected ID=555 after create")
		}
	})

	t.Run("Test Create (Error por email vacío)", func(t *testing.T) {
		_, err := useCase.Create(&userDomain.User{Email: "", Password: "abc"})
		if err == nil {
			t.Error("expected error on create user with empty email")
		}
	})

	t.Run("Test Delete", func(t *testing.T) {
		mockRepo.deleteFn = func(id int) error {
			if id == 101 {
				return nil
			}
			return errors.New("cannot delete")
		}
		err := useCase.Delete(999)
		if err == nil {
			t.Error("expected error for cannot delete")
		}
		err = useCase.Delete(101)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Test Update", func(t *testing.T) {
		mockRepo.updateFn = func(id int, m map[string]interface{}) (*userDomain.User, error) {
			if id != 1001 {
				return nil, errors.New("not found")
			}
			return &userDomain.User{ID: id, UserName: "Updated"}, nil
		}
		_, err := useCase.Update(999, map[string]interface{}{"userName": "any"})
		if err == nil {
			t.Error("expected error, got nil")
		}
		updated, err := useCase.Update(1001, map[string]interface{}{"userName": "whatever"})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if updated.UserName != "Updated" {
			t.Error("expected userName=Updated")
		}
	})
}

func TestNewUserUseCase(t *testing.T) {
	mockRepo := &mockUserService{}
	loggerInstance := setupLogger(t)
	useCase := NewUserUseCase(mockRepo, loggerInstance)
	if reflect.TypeOf(useCase).String() != "*user.UserUseCase" {
		t.Error("expected *user.UserUseCase type")
	}
}
