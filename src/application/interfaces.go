package application

import (
	"github.com/gbrayhan/microservices-go/src/application/usecases/auth"
	"github.com/gbrayhan/microservices-go/src/domain"
	medicineDomain "github.com/gbrayhan/microservices-go/src/domain/medicine"
	userDomain "github.com/gbrayhan/microservices-go/src/domain/user"
)

// UserUseCaseInterface defines the interface for user use cases
type UserUseCaseInterface interface {
	GetAll() (*[]userDomain.User, error)
	GetByID(id int) (*userDomain.User, error)
	Create(newUser *userDomain.User) (*userDomain.User, error)
	GetOneByMap(userMap map[string]interface{}) (*userDomain.User, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*userDomain.User, error)
}

// MedicineUseCaseInterface defines the interface for medicine use cases
type MedicineUseCaseInterface interface {
	GetData(page int64, limit int64, sortBy string, sortDirection string, filters map[string][]string, searchText string, dateRangeFilters []domain.DateRangeFilter) (*medicineDomain.DataMedicine, error)
	GetByID(id int) (*medicineDomain.Medicine, error)
	Create(medicine *medicineDomain.Medicine) (*medicineDomain.Medicine, error)
	GetByMap(medicineMap map[string]any) (*medicineDomain.Medicine, error)
	Delete(id int) error
	Update(id int, medicineMap map[string]any) (*medicineDomain.Medicine, error)
	GetAll() (*[]medicineDomain.Medicine, error)
}

// AuthUseCaseInterface defines the interface for auth use cases
type AuthUseCaseInterface interface {
	Login(user auth.LoginUser) (*auth.SecurityAuthenticatedUser, error)
	AccessTokenByRefreshToken(refreshToken string) (*auth.SecurityAuthenticatedUser, error)
}
