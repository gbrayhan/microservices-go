package user

import "time"

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

type Service interface {
	Get(int) (*User, error)
	GetAll() ([]*User, error)
	Create(*User) error
	GetByMap(map[string]interface{}) map[string]interface{}
	GetById(int) (*User, error)
	Delete(int) error
	Update(int, map[string]interface{}) (*User, error)
}
