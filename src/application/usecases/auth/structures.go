// Package auth provides the use case for authentication
package auth

import "time"

// LoginUser is a struct that contains the request body for the login user
type LoginUser struct {
	Email    string
	Password string
}

// DataUserAuthenticated is a struct that contains the data for the authenticated user
type DataUserAuthenticated struct {
	UserName  string `json:"userName" example:"UserName" gorm:"unique"`
	Email     string `json:"email" example:"some@mail.com" gorm:"unique"`
	FirstName string `json:"firstName" example:"John"`
	LastName  string `json:"lastName" example:"Doe"`
	Status    bool   `json:"status" example:"1"`
	Role      string `json:"role" example:"admin"`
	ID        int    `json:"id" example:"123"`
}

// DataSecurityAuthenticated is a struct that contains the security data for the authenticated user
type DataSecurityAuthenticated struct {
	JWTAccessToken            string    `json:"jwtAccessToken" example:"SomeAccessToken"`
	JWTRefreshToken           string    `json:"jwtRefreshToken" example:"SomeRefreshToken"`
	ExpirationAccessDateTime  time.Time `json:"expirationAccessDateTime" example:"2023-02-02T21:03:53.196419-06:00"`
	ExpirationRefreshDateTime time.Time `json:"expirationRefreshDateTime" example:"2023-02-03T06:53:53.196419-06:00"`
}

// SecurityAuthenticatedUser is a struct that contains the data for the authenticated user
type SecurityAuthenticatedUser struct {
	Data     DataUserAuthenticated     `json:"data"`
	Security DataSecurityAuthenticated `json:"security"`
}
