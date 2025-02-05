package medicine

import "time"

type ResponseMedicine struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	EanCode     string    `json:"eanCode"`
	Laboratory  string    `json:"laboratory"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}

type PaginationResultMedicine struct {
	Data       *[]ResponseMedicine `json:"data"`
	Total      int64               `json:"total"`
	Limit      int64               `json:"limit"`
	Current    int64               `json:"current"`
	NextCursor int64               `json:"nextCursor"`
	PrevCursor int64               `json:"prevCursor"`
	NumPages   int64               `json:"numPages"`
}
