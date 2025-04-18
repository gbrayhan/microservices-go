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

type IMedicineService interface {
	GetAll() (*[]Medicine, error)
	GetData(page int64, limit int64, sortBy string, sortDirection string, filters map[string][]string, searchText string, dateRangeFilters []domain.DateRangeFilter) (*DataMedicine, error)
	GetByID(id int) (*Medicine, error)
	Create(medicine *Medicine) (*Medicine, error)
	GetByMap(medicineMap map[string]any) (*Medicine, error)
	Delete(id int) error
	Update(id int, medicineMap map[string]any) (*Medicine, error)
}
