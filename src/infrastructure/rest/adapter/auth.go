// Package adapter is a layer that connects the infrastructure with the application layer
package adapter

import (
	authService "github.com/gbrayhan/microservices-go/src/application/usecases/auth"
	userRepository "github.com/gbrayhan/microservices-go/src/infrastructure/repository/user"
	authController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/auth"
	"gorm.io/gorm"
)

// AuthAdapter is a function that returns a auth controller
func AuthAdapter(db *gorm.DB) *authController.Controller {
	uRepository := userRepository.Repository{DB: db}
	service := authService.Service{UserRepository: uRepository}
	return &authController.Controller{AuthService: service}
}
