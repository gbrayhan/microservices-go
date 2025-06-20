package psql

import (
	"fmt"
	"os"
	"strings"

	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/medicine"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// loadDatabaseConfig loads database configuration from environment variables
// Returns error if any required environment variable is missing
func loadDatabaseConfig() (DatabaseConfig, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSLMODE")

	// Check for missing required environment variables
	var missingVars []string
	if host == "" {
		missingVars = append(missingVars, "DB_HOST")
	}
	if port == "" {
		missingVars = append(missingVars, "DB_PORT")
	}
	if user == "" {
		missingVars = append(missingVars, "DB_USER")
	}
	if password == "" {
		missingVars = append(missingVars, "DB_PASSWORD")
	}
	if dbName == "" {
		missingVars = append(missingVars, "DB_NAME")
	}
	if sslMode == "" {
		missingVars = append(missingVars, "DB_SSLMODE")
	}

	if len(missingVars) > 0 {
		return DatabaseConfig{}, fmt.Errorf("missing required database environment variables: %s", strings.Join(missingVars, ", "))
	}

	return DatabaseConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
		SSLMode:  sslMode,
	}, nil
}

type PSQLRepository struct {
	DB     *gorm.DB
	Logger *logger.Logger
	Auth   AuthService
}

type AuthService interface {
	HashPassword(password string) (string, error)
}

func NewRepository(db *gorm.DB, loggerInstance *logger.Logger) *PSQLRepository {
	return &PSQLRepository{
		DB:     db,
		Logger: loggerInstance,
	}
}

func (r *PSQLRepository) SetLogger(loggerInstance *logger.Logger) {
	r.Logger = loggerInstance
}

func (r *PSQLRepository) SetAuthService(auth AuthService) {
	r.Auth = auth
}

func (r *PSQLRepository) LoadDBConfig() (DatabaseConfig, error) {
	return loadDatabaseConfig()
}

func (c DatabaseConfig) GetDSN() string {
	return "host=" + c.Host +
		" port=" + c.Port +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.DBName +
		" sslmode=" + c.SSLMode +
		" TimeZone=America/Mexico_City"
}

func (r *PSQLRepository) InitDatabase() error {
	cfg, err := loadDatabaseConfig()
	if err != nil {
		r.Logger.Error("Failed to load database configuration", zap.Error(err))
		return fmt.Errorf("failed to load database configuration: %w", err)
	}

	// Create GORM logger with zap
	gormZap := logger.NewGormLogger(r.Logger.Log).
		LogMode(gormlogger.Warn) // Silent / Error / Warn / Info

	r.DB, err = gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{
		Logger: gormZap,
	})
	if err != nil {
		r.Logger.Error("Error connecting to the database", zap.Error(err))
		return err
	}

	err = r.MigrateEntitiesGORM()
	if err != nil {
		r.Logger.Error("Error migrating the database", zap.Error(err))
		return err
	}

	err = r.SeedInitialUser()
	if err != nil {
		r.Logger.Error("Error seeding initial user", zap.Error(err))
		return err
	}

	r.Logger.Info("Database connection and migrations successful")
	return nil
}

func (r *PSQLRepository) MigrateEntitiesGORM() error {
	// Import the models to register them with GORM
	userModel := &user.User{}
	medicineModel := &medicine.Medicine{}

	// Auto migrate the models to create/update tables
	err := r.DB.AutoMigrate(userModel, medicineModel)
	if err != nil {
		r.Logger.Error("Error migrating database entities", zap.Error(err))
		return err
	}

	r.Logger.Info("Database entities migration completed successfully")
	return nil
}

func (r *PSQLRepository) SeedInitialUser() error {
	email := os.Getenv("START_USER_EMAIL")
	pw := os.Getenv("START_USER_PW")
	if email == "" || pw == "" {
		r.Logger.Info("Initial user seed skipped: START_USER_EMAIL or START_USER_PW not set")
		return nil
	}

	// Check if user already exists
	var existingUser user.User
	err := r.DB.Where("email = ?", email).First(&existingUser).Error
	if err == nil {
		r.Logger.Info("Initial user already exists, skipping seed", zap.String("email", email))
		return nil
	}

	// Create initial user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		r.Logger.Error("Error hashing password for initial user", zap.Error(err))
		return err
	}

	newUser := user.User{
		Email:        email,
		HashPassword: string(hashedPassword),
	}

	err = r.DB.Create(&newUser).Error
	if err != nil {
		r.Logger.Error("Error creating initial user", zap.Error(err))
		return err
	}

	r.Logger.Info("Initial user created successfully", zap.String("email", email))
	return nil
}

// InitPSQLDB initializes the database connection with logger
func InitPSQLDB(loggerInstance *logger.Logger) (*gorm.DB, error) {
	repo := &PSQLRepository{
		Logger: loggerInstance,
	}

	err := repo.InitDatabase()
	if err != nil {
		return nil, err
	}

	return repo.DB, nil
}
