package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gbrayhan/microservices-go/src/infrastructure/config"
	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize application context with dependencies
	appContext, err := di.SetupDependencies()
	if err != nil {
		panic(fmt.Errorf("error initializing application context: %w", err))
	}

	// Setup router
	router := setupRouter(appContext)

	// Setup server
	server := setupServer(router, cfg.Server.Port)

	// Start server
	fmt.Printf("Server running at http://localhost:%s\n", cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil {
		panic(strings.ToLower(err.Error()))
	}
}

func setupRouter(appContext *di.ApplicationContext) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	// Add middlewares
	router.Use(middlewares.ErrorHandler())
	router.Use(middlewares.GinBodyLogMiddleware)
	router.Use(middlewares.CommonHeaders)

	// Setup routes
	routes.ApplicationRouter(router, appContext)
	return router
}

func setupServer(router *gin.Engine, port string) *http.Server {
	return &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    18000 * time.Second,
		WriteTimeout:   18000 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
