// Package medicine contains the medicine controller
package medicine

import "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"

// NewMedicineRequest is a struct that contains the new medicine request information
type NewMedicineRequest struct {
	Name        string `json:"name" example:"Paracetamol" gorm:"unique" binding:"required"`
	Description string `json:"description" example:"Something" binding:"required"`
	Laboratory  string `json:"laboratory" example:"Roche" binding:"required"`
	EanCode     string `json:"eanCode" example:"122000000021" gorm:"unique" binding:"required"`
}

// DataMedicineRequest is a struct that contains the request body for get medicines
type DataMedicineRequest struct {
	Limit           int64                                   `json:"limit" example:"10"`
	Page            int64                                   `json:"page" example:"1"`
	GlobalSearch    string                                  `json:"globalSearch" example:"John"`
	Filters         map[string][]string                     `json:"filters"`
	SorBy           controllers.SortByDataRequest           `json:"sortBy"`
	FieldsDateRange []controllers.FieldDateRangeDataRequest `json:"fieldsDateRange"`
}
