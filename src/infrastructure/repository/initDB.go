package repository

import (
	"fmt"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/medicine"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strconv"
	"time"
)

func InitDB() (*gorm.DB, error) {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASS", "")
	dbName := getEnv("DB_NAME", "postgres")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Mexico_City", dbHost, dbUser, dbPass, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error retrieving sql.DB from gorm.DB: %v", err)
		return nil, err
	}

	maxIdleConns := getEnvAsInt("DB_MAX_IDLE_CONNS", 10)
	maxOpenConns := getEnvAsInt("DB_MAX_OPEN_CONNS", 50)
	connMaxLifetime := getEnvAsInt("DB_CONN_MAX_LIFETIME", 300)

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)

	if err = sqlDB.Ping(); err != nil {
		log.Printf("Error pinging database: %v", err)
		return nil, err
	}

	if err = db.AutoMigrate(&user.User{}, &medicine.Medicine{}); err != nil {
		log.Printf("Error auto-migrating database schema: %v", err)
		return nil, err
	}

	return db, nil
}

func getEnv(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if valStr, ok := os.LookupEnv(key); ok {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}
