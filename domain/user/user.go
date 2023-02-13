// Package user contains the business logic for the user entity
package user

import "time"

// User is a struct that contains the user information
type User struct {
	ID           int
	UserName     string
	Email        string
	FirstName    string
	LastName     string
	Status       bool
	Role         string
	HashPassword string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Service is the interface that provides user methods
type Service interface {
	Get(int) (*User, error)
	GetAll() ([]*User, error)
	Create(*User) error
	GetByMap(map[string]interface{}) map[string]interface{}
	GetByID(int) (*User, error)
	Delete(int) error
	Update(int, map[string]interface{}) (*User, error)
}
