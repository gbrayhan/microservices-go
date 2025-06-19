package repository

import (
	"time"
)

// User model for GORM
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

func (User) TableName() string {
	return "users"
}

// Medicine model for GORM
type Medicine struct {
	ID          int       `gorm:"primaryKey"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	EANCode     string    `gorm:"column:ean_code"`
	Laboratory  string    `gorm:"column:laboratory"`
	CreatedAt   time.Time `gorm:"autoCreateTime:mili"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:mili"`
}

func (Medicine) TableName() string {
	return "medicines"
}
