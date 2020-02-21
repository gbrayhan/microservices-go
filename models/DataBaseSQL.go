package models

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"time"
)

var database Database

type Database struct {
	DBRead  *sql.DB
	DBWrite *sql.DB
}

func init() {
	var err interface{}

	viper.SetConfigFile("config.json")
	if err = viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
		return
	}

	dbUser := viper.GetString("Database.MySQL.Read.Username")
	dbPassword := viper.GetString("Database.MySQL.Read.Password")
	dbHost := viper.GetString("Database.MySQL.Read.Hostname")
	dbPort := viper.GetString("Database.MySQL.Read.Port")
	dbDataBase := viper.GetString("Database.MySQL.Read.Name")
	if database.DBRead, err =
		sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+dbHost+":"+dbPort+")/"+dbDataBase); err != nil {
		panic(fmt.Errorf("Description: %s \n", err))
		return
	}
	database.DBRead.SetConnMaxLifetime(time.Second * 10)

	dbUser = viper.GetString("Database.MySQL.Write.Username")
	dbPassword = viper.GetString("Database.MySQL.Write.Password")
	dbHost = viper.GetString("Database.MySQL.Write.Hostname")
	dbPort = viper.GetString("Database.MySQL.Write.Port")
	dbDataBase = viper.GetString("Database.MySQL.Write.Name")
	if database.DBWrite, err = sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+dbHost+":"+dbPort+")/"+dbDataBase); err != nil {
		panic(fmt.Errorf("Description: %s \n", err))
		return
	}
	database.DBWrite.SetConnMaxLifetime(time.Second * 10)

}
