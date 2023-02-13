package adapter

import (
	authService "github.com/gbrayhan/microservices-go/application/usecases/auth"
	userRepository "github.com/gbrayhan/microservices-go/infrastructure/repository/user"
	authController "github.com/gbrayhan/microservices-go/infrastructure/rest/controllers/auth"
	"gorm.io/gorm"
)

func AuthAdapter(db *gorm.DB) *authController.Controller {
	uRepository := userRepository.Repository{DB: db}
	service := authService.Service{UserRepository: uRepository}
	return &authController.Controller{AuthService: service}
}
