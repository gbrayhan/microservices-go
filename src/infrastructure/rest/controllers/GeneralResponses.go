// Package controllers contains the common functions and structures for the controllers
package controllers

// JSONSwagger is a struct that contains the swagger documentation
type JSONSwagger struct {
}

// MessageResponse is a struct that contains the response body for the message
type MessageResponse struct {
	Message string `json:"message"`
}

type SortByDataRequest struct {
	Field     string `json:"field" example:"user_name"`
	Direction string `json:"direction" example:"asc"`
}

type FieldDateRangeDataRequest struct {
	Field     string `json:"field" example:"createdAt"`
	StartDate string `json:"startDate" example:"2021-01-01"`
	EndDate   string `json:"endDate" example:"2021-01-01"`
}
