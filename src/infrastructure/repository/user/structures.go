package user

import (
	"time"
)

type User struct {
	ID           int       `gorm:"primaryKey"`
	UserName     string    `gorm:"column:user_name;unique"`
	Email        string    `gorm:"unique"`
	FirstName    string    `gorm:"column:first_name"`
	LastName     string    `gorm:"column:last_name"`
	Status       bool      `gorm:"column:status"`
	HashPassword string    `gorm:"column:hash_password"`
	CreatedAt    time.Time `gorm:"autoCreateTime:mili"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime:mili"`
}
