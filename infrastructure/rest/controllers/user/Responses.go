// Package user contains the user controller
package user

import "time"

// MessageResponse is a struct that contains the response body for the message
type MessageResponse struct {
	Message string `json:"message"`
}

// ResponseUser is a struct that contains the response body for the user
type ResponseUser struct {
	ID        int       `json:"id" example:"1099"`
	UserName  string    `json:"user" example:"BossonH"`
	Email     string    `json:"email" example:"some@mail.com"`
	FirstName string    `json:"firstName" example:"John"`
	LastName  string    `json:"lastName" example:"Doe"`
	Status    bool      `json:"status" example:"false"`
	CreatedAt time.Time `json:"createdAt,omitempty" example:"2021-02-24 20:19:39" gorm:"autoCreateTime:mili"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" example:"2021-02-24 20:19:39" gorm:"autoUpdateTime:mili"`
}
