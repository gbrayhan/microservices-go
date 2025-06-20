package di

import (
	"sync"

	authUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/auth"
	medicineUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/medicine"
	userUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/user"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
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
	Logger             *logger.Logger
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

var (
	loggerInstance *logger.Logger
	loggerOnce     sync.Once
)

func GetLogger() *logger.Logger {
	loggerOnce.Do(func() {
		loggerInstance, _ = logger.NewLogger()
	})
	return loggerInstance
}

// SetupDependencies creates a new application context with all dependencies
func SetupDependencies(loggerInstance *logger.Logger) (*ApplicationContext, error) {
	// Initialize database with logger
	db, err := psql.InitPSQLDB(loggerInstance)
	if err != nil {
		return nil, err
	}

	// Initialize JWT service (manages its own configuration)
	jwtService := security.NewJWTService()

	// Initialize repositories with logger
	userRepo := user.NewUserRepository(db, loggerInstance)
	medicineRepo := medicine.NewMedicineRepository(db, loggerInstance)

	// Initialize use cases with logger
	authUC := authUseCase.NewAuthUseCase(userRepo, jwtService, loggerInstance)
	userUC := userUseCase.NewUserUseCase(userRepo, loggerInstance)
	medicineUC := medicineUseCase.NewMedicineUseCase(medicineRepo, loggerInstance)

	// Initialize controllers with logger
	authController := authController.NewAuthController(authUC, loggerInstance)
	userController := userController.NewUserController(userUC, loggerInstance)
	medicineController := medicineController.NewMedicineController(medicineUC, loggerInstance)

	return &ApplicationContext{
		DB:                 db,
		Logger:             loggerInstance,
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
	loggerInstance *logger.Logger,
) *ApplicationContext {
	// Initialize use cases with mocked repositories and logger
	authUC := authUseCase.NewAuthUseCase(mockUserRepo, mockJWTService, loggerInstance)
	userUC := userUseCase.NewUserUseCase(mockUserRepo, loggerInstance)
	medicineUC := medicineUseCase.NewMedicineUseCase(mockMedicineRepo, loggerInstance)

	// Initialize controllers with logger
	authController := authController.NewAuthController(authUC, loggerInstance)
	userController := userController.NewUserController(userUC, loggerInstance)
	medicineController := medicineController.NewMedicineController(medicineUC, loggerInstance)

	return &ApplicationContext{
		Logger:             loggerInstance,
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
