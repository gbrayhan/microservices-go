// Package medicine contains the business logic for the medicine entity
package medicine

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	"time"
)

// Medicine is a struct that contains the medicine information
type Medicine struct {
	ID          int
	Name        string
	Description string
	EanCode     string
	Laboratory  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type NewMedicine struct {
	Name        string
	Description string
	EANCode     string
	Laboratory  string
}

// ToDomainMapper is a function that maps the NewMedicine struct to Medicine struct
func (newMedicine *NewMedicine) ToDomainMapper() *Medicine {
	return &Medicine{
		Name:        newMedicine.Name,
		Description: newMedicine.Description,
		EanCode:     newMedicine.EANCode,
		Laboratory:  newMedicine.Laboratory,
	}
}

// DataMedicine is a struct that contains the medicine data and the total of records
type DataMedicine struct {
	Data  *[]Medicine
	Total int64
}

// Service is a interface that contains the methods for the medicine service
type Service interface {
	GetAll() (*[]Medicine, error)
	GetData(page int64, limit int64, sortBy string, sortDirection string, filters map[string][]string, searchText string, dateRangeFilters []domain.DateRangeFilter) (*DataMedicine, error)
	GetByID(id int) (*Medicine, error)
	Create(medicine *NewMedicine) (*Medicine, error)
	GetByMap(medicineMap map[string]any) (*Medicine, error)
	Delete(id int) error
	Update(id int, medicineMap map[string]any) (*Medicine, error)
}
