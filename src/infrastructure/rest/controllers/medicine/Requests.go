package medicine

import "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"

type NewMedicineRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Laboratory  string `json:"laboratory" binding:"required"`
	EanCode     string `json:"eanCode" binding:"required"`
}

type DataMedicineRequest struct {
	Limit           int64                                   `json:"limit" example:"10"`
	Page            int64                                   `json:"page" example:"1"`
	GlobalSearch    string                                  `json:"globalSearch" example:"John"`
	Filters         map[string][]string                     `json:"filters"`
	SorBy           controllers.SortByDataRequest           `json:"sortBy"`
	FieldsDateRange []controllers.FieldDateRangeDataRequest `json:"fieldsDateRange"`
}
