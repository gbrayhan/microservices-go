package user

import (
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
)

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

type SearchResultUser struct {
	Data       *[]User
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

type IUserService interface {
	GetAll() (*[]User, error)
	GetByID(id int) (*User, error)
	Create(newUser *User) (*User, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*User, error)
	SearchPaginated(filters domain.DataFilters) (*SearchResultUser, error)
	SearchByProperty(property string, searchText string) (*[]string, error)
}
