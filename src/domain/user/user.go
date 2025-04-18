package user

import "time"

type User struct {
	ID           int
	UserName     string
	Email        string
	FirstName    string
	LastName     string
	Status       bool
	HashPassword string
	Password     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type IUserService interface {
	GetAll() (*[]User, error)
	GetByID(id int) (*User, error)
	Create(newUser *User) (*User, error)
	GetOneByMap(userMap map[string]interface{}) (*User, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*User, error)
}
