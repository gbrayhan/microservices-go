// Package medicine contains the repository implementation for the medicine entity
package medicine

import (
	domainMedicine "github.com/gbrayhan/microservices-go/domain/medicine"
	"time"
)

type Medicine struct {
	ID          int       `json:"id" example:"123" gorm:"primaryKey"`
	Name        string    `json:"name" example:"Paracetamol" gorm:"unique"`
	Description string    `json:"description" example:"Some Description"`
	EANCode     string    `json:"ean_code" example:"9900000124" gorm:"unique"`
	Laboratory  string    `json:"laboratory" example:"Roche"`
	CreatedAt   time.Time `json:"created_at,omitempty" example:"2021-02-24 20:19:39" gorm:"autoCreateTime:mili"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" example:"2021-02-24 20:19:39" gorm:"autoUpdateTime:mili"`
}

// TableName overrides the table name used by User to `users`
func (*Medicine) TableName() string {
	return "medicines"
}

// PaginationResultMedicine is a struct that contains the pagination result for medicine
type PaginationResultMedicine struct {
	Data       *[]domainMedicine.Medicine
	Total      int64
	Limit      int64
	Current    int64
	NextCursor uint
	PrevCursor uint
	NumPages   int64
}
