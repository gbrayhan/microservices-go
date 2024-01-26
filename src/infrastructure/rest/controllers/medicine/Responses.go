// Package medicine contains the medicine controller
package medicine

import "time"

// MessageResponse is a struct that contains the response body for the message
type MessageResponse struct {
	Message string `json:"message"`
}

type ResponseMedicine struct {
	ID          int       `json:"id" example:"1099"`
	Name        string    `json:"name" example:"Aspirina"`
	Description string    `json:"description" example:"Some Description"`
	Laboratory  string    `json:"laboratory" example:"Some Laboratory"`
	EanCode     string    `json:"eanCode" example:"Some EanCode"`
	CreatedAt   time.Time `json:"createdAt,omitempty" example:"2021-02-24 20:19:39" gorm:"autoCreateTime:mili"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" example:"2021-02-24 20:19:39" gorm:"autoUpdateTime:mili"`
}

// PaginationResultMedicine is the structure for pagination result of client
type PaginationResultMedicine struct {
	Data       *[]ResponseMedicine `json:"data"`
	Total      int64               `json:"total"`
	Limit      int64               `json:"limit"`
	Current    int64               `json:"current"`
	NextCursor int64               `json:"nextCursor"`
	PrevCursor int64               `json:"prevCursor"`
	NumPages   int64               `json:"numPages"`
}
