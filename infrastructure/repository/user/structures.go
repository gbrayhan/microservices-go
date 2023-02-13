// Package user contains the business logic for the user entity
package user

import (
	"time"
)

// User is a struct that contains the user information
type User struct {
	ID           int       `json:"id" example:"1099" gorm:"primaryKey"`
	UserName     string    `json:"userName" example:"UserName" gorm:"unique"`
	Email        string    `json:"email" example:"some@mail.com" gorm:"unique"`
	FirstName    string    `json:"first_name" example:"John"`
	LastName     string    `json:"last_name" example:"Doe"`
	Status       bool      `json:"status" example:"false"`
	HashPassword string    `json:"hash_password" example:"SomeHashPass"`
	CreatedAt    time.Time `json:"created_at,omitempty" example:"2021-02-24 20:19:39" gorm:"autoCreateTime:mili"`
	UpdatedAt    time.Time `json:"updated_at,omitempty" example:"2021-02-24 20:19:39" gorm:"autoUpdateTime:mili"`
}

// TableName overrides the table name used by User to `users`
func (*User) TableName() string {
	return "users"
}

type PaginationResultUser struct {
	Data       []User
	Total      int
	Limit      int
	Current    int
	NextCursor uint
	PrevCursor uint
	NumPages   int
}
