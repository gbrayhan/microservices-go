package adapter

import (
	authUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/auth"
	userRepository "github.com/gbrayhan/microservices-go/src/infrastructure/repository/user"
	authController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/auth"
	"gorm.io/gorm"
)

func AuthAdapter(db *gorm.DB) authController.IAuthController {
	userRepository := userRepository.NewUserRepository(db)
	authUseCase := authUseCase.NewAuthUseCase(userRepository)
	authController := authController.NewAuthController(authUseCase)
	return authController
}
