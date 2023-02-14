// Package adapter is a layer that connects the infrastructure with the application layer
package adapter

import (
	userService "github.com/gbrayhan/microservices-go/application/usecases/user"
	userRepository "github.com/gbrayhan/microservices-go/infrastructure/repository/user"
	userController "github.com/gbrayhan/microservices-go/infrastructure/rest/controllers/user"
	"gorm.io/gorm"
)

// UserAdapter is a function that returns a user controller
func UserAdapter(db *gorm.DB) *userController.Controller {
	uRepository := userRepository.Repository{DB: db}
	service := userService.Service{UserRepository: uRepository}
	return &userController.Controller{UserService: service}
}
