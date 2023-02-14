// Package medicine contains the medicine controller
package medicine

// NewMedicineRequest is a struct that contains the new medicine request information
type NewMedicineRequest struct {
	Name        string `json:"name" example:"Paracetamol" gorm:"unique" binding:"required"`
	Description string `json:"description" example:"Something" binding:"required"`
	Laboratory  string `json:"laboratory" example:"Roche" binding:"required"`
	EanCode     string `json:"ean_code" example:"122000000021" gorm:"unique" binding:"required"`
}
