// Package user contains the user controller
package user

// NewUserRequest is a struct that contains the request body for the new user
type NewUserRequest struct {
	UserName  string `json:"user" example:"someUser" gorm:"unique" binding:"required"`
	Email     string `json:"email" example:"mail@mail.com" gorm:"unique" binding:"required"`
	FirstName string `json:"firstName" example:"John" binding:"required"`
	LastName  string `json:"lastName" example:"Doe" binding:"required"`
	Password  string `json:"password" example:"Password123" binding:"required"`
	Role      string `json:"role" example:"admin" binding:"required"`
}
