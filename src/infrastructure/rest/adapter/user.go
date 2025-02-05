package adapter

import (
	userUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/user"
	userRepository "github.com/gbrayhan/microservices-go/src/infrastructure/repository/user"
	userController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/user"
	"gorm.io/gorm"
)

func UserAdapter(db *gorm.DB) userController.IUserController {
	repository := userRepository.NewUserRepository(db)
	service := userUseCase.NewUserUseCase(repository)
	controller := userController.NewUserController(service)
	return controller
}
