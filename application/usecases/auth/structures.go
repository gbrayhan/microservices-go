package auth

import "time"

type LoginUser struct {
	Email    string
	Password string
}

type DataUserAuthenticated struct {
	UserName  string `json:"userName" example:"UserName" gorm:"unique"`
	Email     string `json:"email" example:"some@mail.com" gorm:"unique"`
	FirstName string `json:"firstName" example:"John"`
	LastName  string `json:"lastName" example:"Doe"`
	Status    bool   `json:"status" example:"1"`
	Role      string `json:"role" example:"admin"`
	ID        int    `json:"id" example:"123"`
}

type DataSecurityAuthenticated struct {
	JWTAccessToken        string    `json:"jwtAccessToken" example:"SomeAccessToken"`
	JWTRefreshToken       string    `json:"jwtRefreshToken" example:"SomeRefreshToken"`
	ExpirationAccessTime  time.Time `json:"expirationAccessTime" example:"2023-02-02T21:03:53.196419-06:00"`
	ExpirationRefreshTime time.Time `json:"expirationRefreshTime" example:"2023-02-03T06:53:53.196419-06:00"`
}

type SecurityAuthenticatedUser struct {
	Data     DataUserAuthenticated     `json:"data"`
	Security DataSecurityAuthenticated `json:"security"`
}
