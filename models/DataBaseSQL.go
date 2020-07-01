package models

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"time"
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
	dbCompanyIT Database
	dbCompanyOp Database
)

// Nodes read and write in database
type Database struct {
	Read  *sql.DB
	Write *sql.DB
}

func init() {
	var infoDB infoDatabase
	viper.SetConfigFile("config.json")
	viper.ReadInConfig()

	mapstructure.Decode(viper.GetStringMap("Databases.MySQL.CompanyIT"), &infoDB)
	dbCompanyIT.upConnectionMysql(&infoDB)

	mapstructure.Decode(viper.GetStringMap("Databases.MySQL.CompanyOp"), &infoDB)
	dbCompanyOp.upConnectionMysql(&infoDB)

	// If you need another database host, use this code HERE:
	// mapstructure.Decode(viper.GetStringMap("Databases.MySQL.NAME"), &infoCompanyOp)
	// dbCompanyOp, _ = infoCompanyOp.upConnectionMysql()

}

// Up new mysql database connection
func (db *Database) upConnectionMysql(info *infoDatabase) (err error) {
	driverRead := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", info.Read.Username, info.Read.Password, info.Read.Hostname, info.Read.Port, info.Read.Name)
	db.Read, err = sql.Open("mysql", driverRead)
	db.Read.SetConnMaxLifetime(time.Second * 10)

	driverWrite := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", info.Write.Username, info.Write.Password, info.Write.Hostname, info.Write.Port, info.Write.Name)
	db.Write, err = sql.Open("mysql", driverWrite)
	db.Write.SetConnMaxLifetime(time.Second * 10)
	return
}
