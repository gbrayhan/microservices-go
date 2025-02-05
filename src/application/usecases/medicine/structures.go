package medicine

import domainMedicine "github.com/gbrayhan/microservices-go/src/domain/medicine"

type NewMedicine struct {
	Name        string `json:"name" example:"Paracetamol"`
	Description string `json:"description" example:"Some Description"`
	EANCode     string `json:"ean_code" example:"9900000124"`
	Laboratory  string `json:"laboratory" example:"Roche"`
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
