package di

import (
	authUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/auth"
	medicineUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/medicine"
	userUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/user"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/medicine"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
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
	UserRepository     user.UserRepositoryInterface
	MedicineRepository medicine.MedicineRepositoryInterface
	AuthUseCase        authUseCase.IAuthUseCase
	UserUseCase        userUseCase.IUserUseCase
	MedicineUseCase    medicineUseCase.IMedicineUseCase
}

// SetupDependencies creates a new application context with all dependencies
func SetupDependencies() (*ApplicationContext, error) {
	// Initialize database
	db, err := psql.InitPSQLDB()
	if err != nil {
		return nil, err
	}

	// Initialize JWT service (manages its own configuration)
	jwtService := security.NewJWTService()

	// Initialize repositories
	userRepo := user.NewUserRepository(db)
	medicineRepo := medicine.NewMedicineRepository(db)

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
		UserRepository:     userRepo,
		MedicineRepository: medicineRepo,
		AuthUseCase:        authUC,
		UserUseCase:        userUC,
		MedicineUseCase:    medicineUC,
	}, nil
}

// NewTestApplicationContext creates an application context for testing with mocked dependencies
func NewTestApplicationContext(
	mockUserRepo user.UserRepositoryInterface,
	mockMedicineRepo medicine.MedicineRepositoryInterface,
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
		UserRepository:     mockUserRepo,
		MedicineRepository: mockMedicineRepo,
		AuthUseCase:        authUC,
		UserUseCase:        userUC,
		MedicineUseCase:    medicineUC,
	}
}
