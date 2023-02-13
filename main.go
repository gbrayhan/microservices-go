package main

import (
	"fmt"
	"github.com/gbrayhan/microservices-go/infrastructure/repository/config"
	errorsController "github.com/gbrayhan/microservices-go/infrastructure/rest/controllers/errors"
	"github.com/gbrayhan/microservices-go/infrastructure/rest/middlewares"
	"net/http"
	"strings"
	"time"

	limit "github.com/aviddiviner/gin-limit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/gbrayhan/microservices-go/infrastructure/rest/routes"
)

func main() {
	router := gin.Default()
	router.Use(limit.MaxAllowed(200))
	router.Use(cors.Default())
	var err error
	DB, err := config.GormOpen()
	if err != nil {
		_ = fmt.Errorf("fatal error in database file: %s \n", err)
		panic(err)
	}
	router.Use(middlewares.GinBodyLogMiddleware)
	router.Use(errorsController.Handler)
	routes.ApplicationV1Router(router, DB)
	startServer(router)

}

func startServer(router http.Handler) {
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		_ = fmt.Errorf("fatal error in config file: %s \n", err.Error())
		panic(err)

	}
	serverPort := fmt.Sprintf(":%s", viper.GetString("ServerPort"))
	s := &http.Server{
		Addr:           serverPort,
		Handler:        router,
		ReadTimeout:    18000 * time.Second,
		WriteTimeout:   18000 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		_ = fmt.Errorf("fatal error description: %s \n", strings.ToLower(err.Error()))
		panic(err)

	}
}
