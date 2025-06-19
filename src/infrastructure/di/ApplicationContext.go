package di

import (
	authUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/auth"
	medicineUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/medicine"
	userUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/user"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository"
	medicineRepository "github.com/gbrayhan/microservices-go/src/infrastructure/repository/medicine"
	userRepository "github.com/gbrayhan/microservices-go/src/infrastructure/repository/user"
	authController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/auth"
	medicineController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/medicine"
	userController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/user"
	"github.com/gbrayhan/microservices-go/src/infrastructure/security"
	"gorm.io/gorm"
)

// ApplicationContext holds all application dependencies and services
type ApplicationContext struct {
	DB                 *gorm.DB
	AuthController     authController.IAuthController
	UserController     userController.IUserController
	MedicineController medicineController.IMedicineController
	JWTService         security.IJWTService
}

// SetupDependencies creates a new application context with all dependencies
func SetupDependencies() (*ApplicationContext, error) {
	// Initialize database
	db, err := repository.InitDB()
	if err != nil {
		return nil, err
	}

	// Initialize JWT service
	jwtService := security.NewJWTService()

	// Initialize repositories
	userRepo := userRepository.NewUserRepository(db)
	medicineRepo := medicineRepository.NewMedicineRepository(db)

	// Initialize use cases
	authUC := authUseCase.NewAuthUseCase(userRepo, jwtService)
	userUC := userUseCase.NewUserUseCase(userRepo)
	medicineUC := medicineUseCase.NewMedicineUseCase(medicineRepo)

	// Initialize controllers
	authController := authController.NewAuthController(authUC)
	userController := userController.NewUserController(userUC)
	medicineController := medicineController.NewMedicineController(medicineUC)

	return &ApplicationContext{
		DB:                 db,
		AuthController:     authController,
		UserController:     userController,
		MedicineController: medicineController,
		JWTService:         jwtService,
	}, nil
}

// NewTestApplicationContext creates an application context for testing with mocked dependencies
func NewTestApplicationContext(
	mockUserRepo userRepository.IUserRepository,
	mockMedicineRepo medicineRepository.IMedicineRepository,
	mockJWTService security.IJWTService,
) *ApplicationContext {
	// Initialize use cases with mocked repositories
	authUC := authUseCase.NewAuthUseCase(mockUserRepo, mockJWTService)
	userUC := userUseCase.NewUserUseCase(mockUserRepo)
	medicineUC := medicineUseCase.NewMedicineUseCase(mockMedicineRepo)

	// Initialize controllers
	authController := authController.NewAuthController(authUC)
	userController := userController.NewUserController(userUC)
	medicineController := medicineController.NewMedicineController(medicineUC)

	return &ApplicationContext{
		AuthController:     authController,
		UserController:     userController,
		MedicineController: medicineController,
		JWTService:         mockJWTService,
	}
}
