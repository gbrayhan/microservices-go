package repository

import "gorm.io/gorm"

// Repository is a struct that contains the database implementation for user entity
type Repository struct {
	DB *gorm.DB
}
