// Package medicine contains the business logic for the medicine entity
package medicine

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	"time"
)

// Medicine is a struct that contains the medicine information
type Medicine struct {
	ID          int       `json:"id" example:"123"`
	Name        string    `json:"name" example:"Paracetamol"`
	Description string    `json:"description" example:"Some Description"`
	EanCode     string    `json:"ean_code" example:"9900000124"`
	Laboratory  string    `json:"laboratory" example:"Roche"`
	CreatedAt   time.Time `json:"created_at,omitempty" `
	UpdatedAt   time.Time `json:"updated_at,omitempty" example:"2021-02-24 20:19:39"`
}

type NewMedicine struct {
	Name        string `json:"name" example:"Paracetamol"`
	Description string `json:"description" example:"Some Description"`
	EANCode     string `json:"ean_code" example:"9900000124"`
	Laboratory  string `json:"laboratory" example:"Roche"`
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
	GetData(page int64, limit int64, sortBy string, sortDirection string, filters map[string][]string, searchText string, dateRangeFilters []domain.DateRangeFilter) (*DataMedicine, error)
	GetByID(id int) (*Medicine, error)
	Create(medicine *NewMedicine) (*Medicine, error)
	GetByMap(medicineMap map[string]any) (*Medicine, error)
	Delete(id int) error
	Update(id int, medicineMap map[string]any) (*Medicine, error)
}
