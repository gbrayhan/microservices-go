package main

import (
	"fmt"
	limit "github.com/aviddiviner/gin-limit"
	"github.com/gbrayhan/microservices-go/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func main() {
	router := gin.Default()

	initialGinConfig(router)
	routes.ApplicationV1Router(router)
	startServer(router)

}

func initialGinConfig(router *gin.Engine) {
	router.Use(limit.MaxAllowed(200))
	router.Use(cors.Default())
	router.Static("/public/static", "public/static")
	router.LoadHTMLGlob("views/**/*")
}

func startServer(router *gin.Engine) {
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error in config file: %s \n", err))
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
		panic(fmt.Errorf("Fatal Error Description: %s \n", err))
	}

}
