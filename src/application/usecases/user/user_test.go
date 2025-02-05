package user

import (
	"errors"
	"reflect"
	"testing"

	userDomain "github.com/gbrayhan/microservices-go/src/domain/user"
)

type mockUserService struct {
	getAllFn      func() (*[]userDomain.User, error)
	getByIDFn     func(id int) (*userDomain.User, error)
	createFn      func(u *userDomain.User) (*userDomain.User, error)
	getOneByMapFn func(m map[string]interface{}) (*userDomain.User, error)
	deleteFn      func(id int) error
	updateFn      func(id int, m map[string]interface{}) (*userDomain.User, error)
}

func (m *mockUserService) GetAll() (*[]userDomain.User, error) {
	return m.getAllFn()
}
func (m *mockUserService) GetByID(id int) (*userDomain.User, error) {
	return m.getByIDFn(id)
}
func (m *mockUserService) Create(newUser *userDomain.User) (*userDomain.User, error) {
	return m.createFn(newUser)
}
func (m *mockUserService) GetOneByMap(userMap map[string]interface{}) (*userDomain.User, error) {
	return m.getOneByMapFn(userMap)
}
func (m *mockUserService) Delete(id int) error {
	return m.deleteFn(id)
}
func (m *mockUserService) Update(id int, userMap map[string]interface{}) (*userDomain.User, error) {
	return m.updateFn(id, userMap)
}

func TestUserUseCase(t *testing.T) {

	mockRepo := &mockUserService{}
	useCase := NewUserUseCase(mockRepo)

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

	t.Run("Test GetOneByMap", func(t *testing.T) {
		mockRepo.getOneByMapFn = func(m map[string]interface{}) (*userDomain.User, error) {
			if m["email"] == "no_exist" {
				return nil, errors.New("not found")
			}
			return &userDomain.User{ID: 777}, nil
		}
		_, err := useCase.GetOneByMap(map[string]interface{}{"email": "no_exist"})
		if err == nil {
			t.Error("expected error for not found")
		}
		single, err := useCase.GetOneByMap(map[string]interface{}{"email": "yes_exist"})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if single.ID != 777 {
			t.Error("expected ID=777")
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
	uc := NewUserUseCase(mockRepo)
	if reflect.TypeOf(uc).String() != "*user.UserUseCase" {
		t.Error("expected *user.UserUseCase type")
	}
}
