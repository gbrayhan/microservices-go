package repository

import (
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// InitDB initializes the database connection and migrations using the new repository structure
func InitDB() (*gorm.DB, error) {
	logger := log.New(os.Stdout, "[DB] ", log.LstdFlags)

	repo := NewRepository(nil, logger)
	repo.SetAuthService(NewAuthService())

	err := repo.InitDatabase()
	if err != nil {
		logger.Printf("Error initializing database: %v", err)
		return nil, err
	}

	sqlDB, err := repo.DB.DB()
	if err != nil {
		logger.Printf("Error retrieving sql.DB from gorm.DB: %v", err)
		return nil, err
	}

	maxIdleConns := getEnvAsInt("DB_MAX_IDLE_CONNS", 10)
	maxOpenConns := getEnvAsInt("DB_MAX_OPEN_CONNS", 50)
	connMaxLifetime := getEnvAsInt("DB_CONN_MAX_LIFETIME", 300)

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)

	if err = sqlDB.Ping(); err != nil {
		logger.Printf("Error pinging database: %v", err)
		return nil, err
	}

	return repo.DB, nil
}

func getEnvAsInt(key string, defaultVal int) int {
	if valStr, ok := os.LookupEnv(key); ok {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}
