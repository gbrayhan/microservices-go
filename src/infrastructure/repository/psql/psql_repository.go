package psql

import (
	"errors"
	"log"
	"os"

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
func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASSWORD", "password"),
		DBName:   getEnvOrDefault("DB_NAME", "microservices_go"),
		SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
	}
}

type PSQLRepository struct {
	DB     *gorm.DB
	Logger *log.Logger
	Auth   AuthService
}

type AuthService interface {
	HashPassword(password string) (string, error)
}

func NewRepository(db *gorm.DB, logger *log.Logger) *PSQLRepository {
	return &PSQLRepository{
		DB:     db,
		Logger: logger,
	}
}

func (r *PSQLRepository) SetLogger(logger *log.Logger) {
	r.Logger = logger
}

func (r *PSQLRepository) SetAuthService(auth AuthService) {
	r.Auth = auth
}

func (r *PSQLRepository) LoadDBConfig() (DatabaseConfig, error) {
	config := loadDatabaseConfig()

	// Check if any required environment variables are missing
	if config.Host == "" || config.Port == "" || config.User == "" || config.Password == "" || config.DBName == "" || config.SSLMode == "" {
		return DatabaseConfig{}, errors.New("missing required database environment variables")
	}

	return config, nil
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

func (r *PSQLRepository) MigrateEntitiesGORM() error {
	// Note: This will be called from the main application after importing the models
	r.Logger.Println("Database entities migration should be handled by the main application")
	return nil
}

func (r *PSQLRepository) SeedInitialUser() error {
	email := os.Getenv("START_USER_EMAIL")
	pw := os.Getenv("START_USER_PW")
	if email == "" || pw == "" {
		r.Logger.Println("Initial user seed skipped: START_USER_EMAIL or START_USER_PW not set")
		return nil
	}

	// Note: This will be handled by the user repository
	r.Logger.Println("User seeding should be handled by the user repository")
	return nil
}

// InitPSQLDB initializes the database connection
func InitPSQLDB() (*gorm.DB, error) {
	repo := &PSQLRepository{
		Logger: log.New(os.Stdout, "DB: ", log.LstdFlags),
	}

	err := repo.InitDatabase()
	if err != nil {
		return nil, err
	}

	return repo.DB, nil
}

// Helper function
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
