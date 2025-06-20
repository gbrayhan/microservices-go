package medicine

import (
	"time"
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
	GetByID(id int) (*Medicine, error)
	Create(medicine *Medicine) (*Medicine, error)
	Delete(id int) error
	Update(id int, medicineMap map[string]any) (*Medicine, error)
}
