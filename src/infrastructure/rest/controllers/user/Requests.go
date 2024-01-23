// Package user contains the user controller
package user

import "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"

// NewUserRequest is a struct that contains the request body for the new user
type NewUserRequest struct {
	UserName  string `json:"user" example:"someUser" gorm:"unique" binding:"required"`
	Email     string `json:"email" example:"mail@mail.com" gorm:"unique" binding:"required"`
	FirstName string `json:"firstName" example:"John" binding:"required"`
	LastName  string `json:"lastName" example:"Doe" binding:"required"`
	Password  string `json:"password" example:"Password123" binding:"required"`
	Role      string `json:"role" example:"admin" binding:"required"`
}

// DataUserRequest is a struct that contains the request body for get users
type DataUserRequest struct {
	Limit           int64                                   `json:"limit" example:"10"`
	Page            int64                                   `json:"page" example:"1"`
	GlobalSearch    string                                  `json:"globalSearch" example:"John"`
	Filters         map[string][]string                     `json:"filters"`
	SorBy           controllers.SortByDataRequest           `json:"sortBy"`
	FieldsDateRange []controllers.FieldDateRangeDataRequest `json:"fieldsDateRange"`
}
