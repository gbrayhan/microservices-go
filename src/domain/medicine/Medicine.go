package medicine

import (
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
)

type Medicine struct {
	ID          int
	Name        string
	Description string
	EanCode     string
	Laboratory  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type DataMedicine struct {
	Data  *[]Medicine
	Total int64
}

type SearchResultMedicine struct {
	Data       *[]Medicine
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

type IMedicineService interface {
	GetAll() (*[]Medicine, error)
	GetByID(id int) (*Medicine, error)
	Create(medicine *Medicine) (*Medicine, error)
	Delete(id int) error
	Update(id int, medicineMap map[string]any) (*Medicine, error)
	SearchPaginated(filters domain.DataFilters) (*SearchResultMedicine, error)
	SearchByProperty(property string, searchText string) (*[]string, error)
}
