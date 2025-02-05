package medicine

import (
	"time"

	domainMedicine "github.com/gbrayhan/microservices-go/src/domain/medicine"
)

type Medicine struct {
	ID          int    `gorm:"primaryKey"`
	Name        string `gorm:"unique"`
	Description string
	EANCode     string `gorm:"unique"`
	Laboratory  string
	CreatedAt   time.Time `gorm:"autoCreateTime:milli"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:milli"`
}

type PaginationResultMedicine struct {
	Data       *[]domainMedicine.Medicine
	Total      int64
	Limit      int64
	Current    int64
	NextCursor uint
	PrevCursor uint
	NumPages   int64
}
