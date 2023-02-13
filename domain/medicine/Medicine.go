// Package medicine contains the business logic for the medicine entity
package medicine

import (
	"time"
)

type Medicine struct {
	ID          int       `json:"id" example:"123"`
	Name        string    `json:"name" example:"Paracetamol"`
	Description string    `json:"description" example:"Some Description"`
	EANCode     string    `json:"ean_code" example:"9900000124"`
	Laboratory  string    `json:"laboratory" example:"Roche"`
	CreatedAt   time.Time `json:"created_at,omitempty" `
	UpdatedAt   time.Time `json:"updated_at,omitempty" example:"2021-02-24 20:19:39"`
}

type Service interface {
	Get(int) (*Medicine, error)
	GetAll() ([]*Medicine, error)
	Create(*Medicine) error
	GetByMap(map[string]interface{}) map[string]interface{}
	GetByID(int) (*Medicine, error)
	Delete(int) error
	Update(int, map[string]interface{}) (*Medicine, error)
}
