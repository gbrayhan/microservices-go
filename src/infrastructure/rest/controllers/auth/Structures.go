// Package auth contains the auth controller
package auth

// LoginRequest is a struct that contains the login request information
type LoginRequest struct {
	Email    string `json:"email" example:"gbrayhan@gmail.com" gorm:"unique" binding:"required"`
	Password string `json:"password" example:"Password123" binding:"required"`
}

// AccessTokenRequest is a struct that contains the login request information
type AccessTokenRequest struct {
	RefreshToken string `json:"refreshToken" example:"badbunybabybebe" binding:"required"`
}
