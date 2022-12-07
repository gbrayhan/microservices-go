package main

import (
  "fmt"
  "github.com/gbrayhan/microservices-go/config"
  errorsController "github.com/gbrayhan/microservices-go/controllers/errors"
  "net/http"
  "strings"
  "time"

  limit "github.com/aviddiviner/gin-limit"
  "github.com/gin-contrib/cors"
  "github.com/gin-gonic/gin"
  "github.com/spf13/viper"

  "github.com/gbrayhan/microservices-go/middlewares"
  "github.com/gbrayhan/microservices-go/routes"
)

func main() {
  router := gin.Default()
  initialGinConfig(router)
  router.Use(middlewares.GinBodyLogMiddleware)
  router.Use(errorsController.Handler)
  routes.ApplicationV1Router(router)
  startServer(router)

}

func initialGinConfig(router *gin.Engine) {
  router.Use(limit.MaxAllowed(200))
  router.Use(cors.Default())
  var err error
  config.DB, err = config.GormOpen()

  if err != nil {
    _ = fmt.Errorf("fatal error in database file: %s \n", err)
  }

}

func startServer(router http.Handler) {
  viper.SetConfigFile("config.json")
  if err := viper.ReadInConfig(); err != nil {
    _ = fmt.Errorf("fatal error in config file: %s \n", err.Error())
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
  }
}
