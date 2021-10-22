package config

import (
	"database/sql"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type infoDatabaseSQL struct {
	Read struct {
		Hostname  string
		Name      string
		Username  string
		Password  string
		Port      string
		Parameter string
		DriverConn string
	}
	Write struct {
		Hostname  string
		Name      string
		Username  string
		Password  string
		Port      string
		Parameter string
		DriverConn string
	}
}
func (infoDB *infoDatabaseSQL) getDiverConn(nameMap string) (err error) {
	viper.SetConfigFile("config.json")
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = mapstructure.Decode(viper.GetStringMap(nameMap), infoDB)
	if err != nil {
		return
	}

	infoDB.Read.DriverConn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", infoDB.Read.Username, infoDB.Read.Password, infoDB.Read.Hostname, infoDB.Read.Port, infoDB.Read.Name)
	infoDB.Write.DriverConn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", infoDB.Write.Username, infoDB.Write.Password, infoDB.Write.Hostname, infoDB.Write.Port, infoDB.Write.Name)
	return
}



// Nodes read and write in databaseSQL
type databaseSQL struct {
	Read  *sql.DB
	Write *sql.DB
}

// Up new mysql databaseSQL connection
func (db *databaseSQL) upConnectionMysql(info *infoDatabaseSQL) (err error) {
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
