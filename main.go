package main

import (
	"fmt"
	"github.com/aviddiviner/gin-limit"
	"github.com/banwire/microservice_golang/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func main() {

	viper.SetConfigFile("config.json")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal Error in config file: %s \n", err))
		return
	}

	router := gin.Default()
	router.Use(limit.MaxAllowed(200))
	router.Use(cors.Default())

	router.Static("/public/static", "public/static")
	router.LoadHTMLGlob("views/**/*")

	routes.ApplicationV1Router(router)

	serverPort := fmt.Sprintf(":%s", viper.GetString("ServerPort"))
	s := &http.Server{
		Addr:           serverPort,
		Handler:        router,
		ReadTimeout:    18000 * time.Second,
		WriteTimeout:   18000 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		panic(fmt.Errorf("Fatal Error Description: %s \n", err))
		return
	}

}
