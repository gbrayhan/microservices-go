package repository

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Repository struct {
	DB     *gorm.DB
	Logger *log.Logger
	Auth   AuthService
}

type AuthService interface {
	HashPassword(password string) (string, error)
}

func NewRepository(db *gorm.DB, logger *log.Logger) *Repository {
	return &Repository{
		DB:     db,
		Logger: logger,
	}
}

func (r *Repository) SetLogger(logger *log.Logger) {
	r.Logger = logger
}

func (r *Repository) SetAuthService(auth AuthService) {
	r.Auth = auth
}

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string
}

func (r *Repository) LoadDBConfig() (Config, error) {
	dbHost := loadEnvVar("DB_HOST")
	dbPort := loadEnvVar("DB_PORT")
	dbUser := loadEnvVar("DB_USER")
	dbPassword := loadEnvVar("DB_PASS")
	dbName := loadEnvVar("DB_NAME")
	sslMode := loadEnvVar("DB_SSLMODE")

	// Check if any required environment variables are missing
	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" || sslMode == "" {
		return Config{}, errors.New("missing required database environment variables")
	}

	return Config{
		DBHost:     dbHost,
		DBPort:     dbPort,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,
		SSLMode:    sslMode,
	}, nil
}

func loadEnvVar(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return ""
}

func (c Config) GetDSN() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=" + c.SSLMode +
		" TimeZone=America/Mexico_City"
}

func (r *Repository) InitDatabase() error {
	cfg, err := r.LoadDBConfig()
	if err != nil {
		return err
	}

	r.DB, err = gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		r.Logger.Printf("Error connecting to the database: %v", err)
		return err
	}

	err = r.MigrateEntitiesGORM()
	if err != nil {
		r.Logger.Printf("Error migrating the database: %v", err)
		return err
	}

	r.Logger.Println("Database connection and migrations successful")
	return nil
}

func (r *Repository) MigrateEntitiesGORM() error {
	err := r.DB.AutoMigrate(&User{}, &Medicine{})
	if err != nil {
		r.Logger.Printf("Error migrating database entities: %v", err)
		return err
	}
	r.Logger.Println("Database entities migrated successfully")

	if err := r.SeedInitialUser(); err != nil {
		r.Logger.Printf("Error seeding initial user: %v", err)
		return err
	}
	r.Logger.Println("Seeding completed")
	return nil
}

func (r *Repository) SeedInitialUser() error {
	email := os.Getenv("START_USER_EMAIL")
	pw := os.Getenv("START_USER_PW")
	if email == "" || pw == "" {
		r.Logger.Println("Initial user seed skipped: START_USER_EMAIL or START_USER_PW not set")
		return nil
	}

	var count int64
	if err := r.DB.Model(&User{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		r.Logger.Printf("Error checking existing user: %v", err)
		return err
	}

	if count == 0 {
		hashed, err := r.Auth.HashPassword(pw)
		if err != nil {
			r.Logger.Printf("Error hashing initial user password: %v", err)
			return err
		}

		user := User{
			Email:        email,
			HashPassword: hashed,
			FirstName:    "John",
			LastName:     "Doe",
			UserName:     "admin",
			Status:       true,
		}
		if err := r.DB.Create(&user).Error; err != nil {
			r.Logger.Printf("Error creating initial user: %v", err)
			return err
		}
		r.Logger.Printf("Created initial user: %s", email)
	} else {
		r.Logger.Printf("Initial user already exists: %s (count: %d)", email, count)
	}

	return nil
}

// Error types and AppError struct
const (
	NotFound                     = "NotFound"
	notFoundMessage              = "record not found"
	ValidationError              = "ValidationError"
	validationErrorMessage       = "validation error"
	ResourceAlreadyExists        = "ResourceAlreadyExists"
	alreadyExistsErrorMessage    = "resource already exists"
	RepositoryError              = "RepositoryError"
	repositoryErrorMessage       = "error in repository operation"
	NotAuthenticated             = "NotAuthenticated"
	notAuthenticatedErrorMessage = "not authenticated"
	TokenGeneratorError          = "TokenGeneratorError"
	tokenGeneratorErrorMessage   = "error in token generation"
	NotAuthorized                = "NotAuthorized"
	notAuthorizedErrorMessage    = "not authorized"
	UnknownError                 = "UnknownError"
	unknownErrorMessage          = "something went wrong"
)

type AppError struct {
	Err  error
	Type string
}

func NewAppError(err error, errType string) *AppError {
	return &AppError{
		Err:  err,
		Type: errType,
	}
}

func NewAppErrorWithType(errType string) *AppError {
	var err error

	switch errType {
	case NotFound:
		err = errors.New(notFoundMessage)
	case ValidationError:
		err = errors.New(validationErrorMessage)
	case ResourceAlreadyExists:
		err = errors.New(alreadyExistsErrorMessage)
	case RepositoryError:
		err = errors.New(repositoryErrorMessage)
	case NotAuthenticated:
		err = errors.New(notAuthenticatedErrorMessage)
	case NotAuthorized:
		err = errors.New(notAuthorizedErrorMessage)
	case TokenGeneratorError:
		err = errors.New(tokenGeneratorErrorMessage)
	default:
		err = errors.New(unknownErrorMessage)
	}

	return &AppError{
		Err:  err,
		Type: errType,
	}
}

func (appErr *AppError) Error() string {
	return appErr.Err.Error()
}

func GenerateNewUUID() (string, error) {
	uuid := make([]byte, 16)
	if _, err := rand.Read(uuid); err != nil {
		return "", err
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
