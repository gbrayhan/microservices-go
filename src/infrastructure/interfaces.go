package infrastructure

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	medicineDomain "github.com/gbrayhan/microservices-go/src/domain/medicine"
	userDomain "github.com/gbrayhan/microservices-go/src/domain/user"
)

// UserRepositoryInterface defines the interface for user repository operations
type UserRepositoryInterface interface {
	GetAll() (*[]userDomain.User, error)
	Create(userDomain *userDomain.User) (*userDomain.User, error)
	GetOneByMap(userMap map[string]interface{}) (*userDomain.User, error)
	GetByID(id int) (*userDomain.User, error)
	Update(id int, userMap map[string]interface{}) (*userDomain.User, error)
	Delete(id int) error
}

// MedicineRepositoryInterface defines the interface for medicine repository operations
type MedicineRepositoryInterface interface {
	GetAll() (*[]medicineDomain.Medicine, error)
	GetData(page int64, limit int64, sortBy string, sortDirection string, filters map[string][]string, searchText string, dateRangeFilters []domain.DateRangeFilter) (*medicineDomain.DataMedicine, error)
	GetByID(id int) (*medicineDomain.Medicine, error)
	Create(medicine *medicineDomain.Medicine) (*medicineDomain.Medicine, error)
	GetByMap(medicineMap map[string]any) (*medicineDomain.Medicine, error)
	Delete(id int) error
	Update(id int, medicineMap map[string]any) (*medicineDomain.Medicine, error)
}

// DatabaseInterface defines the interface for database operations
type DatabaseInterface interface {
	Connect() error
	GetDB() interface{}
	Close() error
}
