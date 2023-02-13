// Package config provides the database connection
package config

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type infoDatabaseSQL struct {
	Read struct {
		Hostname   string
		Name       string
		Username   string
		Password   string
		Port       string
		Parameter  string
		DriverConn string
	}
	Write struct {
		Hostname   string
		Name       string
		Username   string
		Password   string
		Port       string
		Parameter  string
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

	infoDB.Read.DriverConn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		infoDB.Read.Username, infoDB.Read.Password, infoDB.Read.Hostname, infoDB.Read.Port, infoDB.Read.Name)
	infoDB.Write.DriverConn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		infoDB.Write.Username, infoDB.Write.Password, infoDB.Write.Hostname, infoDB.Write.Port, infoDB.Write.Name)
	return
}
