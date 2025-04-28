package main

import (
	"fmt"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	DB, err := repository.InitDB()
	if err != nil {
		panic(fmt.Errorf("error initializing the database: %w", err))
	}

	router.Use(middlewares.ErrorHandler())
	router.Use(middlewares.GinBodyLogMiddleware)
	router.Use(middlewares.CommonHeaders)
	routes.ApplicationRouter(router, DB)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	s := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    18000 * time.Second,
		WriteTimeout:   18000 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Printf("Server running at http://localhost:%s\n", port)
	if err := s.ListenAndServe(); err != nil {
		panic(strings.ToLower(err.Error()))
	}
}
