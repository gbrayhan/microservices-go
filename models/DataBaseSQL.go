package models

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type infoDatabase struct {
	Read struct {
		Hostname  string
		Name      string
		Username  string
		Password  string
		Port      string
		Parameter string
	}
	Write struct {
		Hostname  string
		Name      string
		Username  string
		Password  string
		Port      string
		Parameter string
	}
}

// Host databases to work
var (
	dbBoilerplateGo database
	// dbOtherDB       Database
)

// Nodes read and write in database
type database struct {
	Read  *sql.DB
	Write *sql.DB
}

func init() {
	var infoDB infoDatabase
	viper.SetConfigFile("config.json")
	_ = viper.ReadInConfig()

	_ = mapstructure.Decode(viper.GetStringMap("Databases.MySQL.BoilerplateGo"), &infoDB)
	_ = dbBoilerplateGo.upConnectionMysql(&infoDB)

	// _ = mapstructure.Decode(viper.GetStringMap("Databases.MySQL.OtherDB"), &infoDB)
	// _ = dbOtherDB.upConnectionMysql(&infoDB)

	// If you need another database host, use this code HERE:
}

// Up new mysql database connection
func (db *database) upConnectionMysql(info *infoDatabase) (err error) {
	driverRead := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", info.Read.Username, info.Read.Password, info.Read.Hostname, info.Read.Port, info.Read.Name)
	db.Read, err = sql.Open("mysql", driverRead)
	db.Read.SetConnMaxLifetime(time.Second * 10)
	if err != nil {
		return
	}

	driverWrite := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", info.Write.Username, info.Write.Password, info.Write.Hostname, info.Write.Port, info.Write.Name)
	db.Write, err = sql.Open("mysql", driverWrite)
	db.Write.SetConnMaxLifetime(time.Second * 10)
	return
}
