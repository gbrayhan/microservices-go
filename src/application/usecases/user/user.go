package user

import (
	userDomain "github.com/gbrayhan/microservices-go/src/domain/user"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	"golang.org/x/crypto/bcrypt"
)

type IUserUseCase interface {
	GetAll() (*[]userDomain.User, error)
	GetByID(id int) (*userDomain.User, error)
	Create(newUser *userDomain.User) (*userDomain.User, error)
	GetOneByMap(userMap map[string]interface{}) (*userDomain.User, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*userDomain.User, error)
}

type UserUseCase struct {
	userRepository user.UserRepositoryInterface
}

func NewUserUseCase(userRepository user.UserRepositoryInterface) IUserUseCase {
	return &UserUseCase{
		userRepository: userRepository,
	}
}

func (s *UserUseCase) GetAll() (*[]userDomain.User, error) {
	return s.userRepository.GetAll()
}

func (s *UserUseCase) GetByID(id int) (*userDomain.User, error) {
	return s.userRepository.GetByID(id)
}

func (s *UserUseCase) Create(newUser *userDomain.User) (*userDomain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return &userDomain.User{}, err
	}
	newUser.HashPassword = string(hash)
	newUser.Status = true

	return s.userRepository.Create(newUser)
}

func (s *UserUseCase) GetOneByMap(userMap map[string]interface{}) (*userDomain.User, error) {
	return s.userRepository.GetOneByMap(userMap)
}

func (s *UserUseCase) Delete(id int) error {
	return s.userRepository.Delete(id)
}

func (s *UserUseCase) Update(id int, userMap map[string]interface{}) (*userDomain.User, error) {
	return s.userRepository.Update(id, userMap)
}
