// Package auth contains the auth controller
package auth

// LoginRequest is a struct that contains the login request information
type LoginRequest struct {
	Email    string `json:"email" example:"mail@mail.com" gorm:"unique" binding:"required"`
	Password string `json:"password" example:"Password123" binding:"required"`
}
