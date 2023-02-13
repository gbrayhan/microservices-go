// Package config provides the database connection
package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

var Repo Repository

type Repository struct {
	DB *gorm.DB
}

// DBConfig represents db configuration
type DBConfig struct {
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
}

func GormOpen() (gormDB *gorm.DB, err error) {
	var infoDatabase infoDatabaseSQL
	err = infoDatabase.getDiverConn("Databases.MySQL.BoilerplateGo")
	if err != nil {
		return nil, err
	}
	gormDB, err = gorm.Open(mysql.Open(infoDatabase.Write.DriverConn), &gorm.Config{
		PrepareStmt: true,
		Logger:      logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return
	}

	dialector := mysql.New(mysql.Config{
		DSN: infoDatabase.Read.DriverConn,
	})

	err = gormDB.Use(dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{dialector},
	}))
	if err != nil {
		return nil, err
	}
	var result int

	// Test the connection by executing a simple query
	if err = gormDB.Raw("SELECT 1").Scan(&result).Error; err != nil {
		return nil, err
	}

	return
}
