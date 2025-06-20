package user

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	userDomain "github.com/gbrayhan/microservices-go/src/domain/user"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type IUserUseCase interface {
	GetAll() (*[]userDomain.User, error)
	GetByID(id int) (*userDomain.User, error)
	GetByEmail(email string) (*userDomain.User, error)
	Create(newUser *userDomain.User) (*userDomain.User, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*userDomain.User, error)
	SearchPaginated(filters domain.DataFilters) (*userDomain.SearchResultUser, error)
	SearchByProperty(property string, searchText string) (*[]string, error)
}

type UserUseCase struct {
	userRepository user.UserRepositoryInterface
	Logger         *logger.Logger
}

func NewUserUseCase(userRepository user.UserRepositoryInterface, logger *logger.Logger) IUserUseCase {
	return &UserUseCase{
		userRepository: userRepository,
		Logger:         logger,
	}
}

func (s *UserUseCase) GetAll() (*[]userDomain.User, error) {
	s.Logger.Info("Getting all users")
	return s.userRepository.GetAll()
}

func (s *UserUseCase) GetByID(id int) (*userDomain.User, error) {
	s.Logger.Info("Getting user by ID", zap.Int("id", id))
	return s.userRepository.GetByID(id)
}

func (s *UserUseCase) GetByEmail(email string) (*userDomain.User, error) {
	s.Logger.Info("Getting user by email", zap.String("email", email))
	return s.userRepository.GetByEmail(email)
}

func (s *UserUseCase) Create(newUser *userDomain.User) (*userDomain.User, error) {
	s.Logger.Info("Creating new user", zap.String("email", newUser.Email))
	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Logger.Error("Error hashing password", zap.Error(err))
		return &userDomain.User{}, err
	}
	newUser.HashPassword = string(hash)
	newUser.Status = true

	return s.userRepository.Create(newUser)
}

func (s *UserUseCase) Delete(id int) error {
	s.Logger.Info("Deleting user", zap.Int("id", id))
	return s.userRepository.Delete(id)
}

func (s *UserUseCase) Update(id int, userMap map[string]interface{}) (*userDomain.User, error) {
	s.Logger.Info("Updating user", zap.Int("id", id))
	return s.userRepository.Update(id, userMap)
}

func (s *UserUseCase) SearchPaginated(filters domain.DataFilters) (*userDomain.SearchResultUser, error) {
	s.Logger.Info("Searching users with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.userRepository.SearchPaginated(filters)
}

func (s *UserUseCase) SearchByProperty(property string, searchText string) (*[]string, error) {
	s.Logger.Info("Searching users by property",
		zap.String("property", property),
		zap.String("searchText", searchText))
	return s.userRepository.SearchByProperty(property, searchText)
}
