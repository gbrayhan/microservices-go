package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
}

// loadServerConfig loads server configuration from environment variables
func loadServerConfig() ServerConfig {
	return ServerConfig{
		Port: getEnvOrDefault("SERVER_PORT", "8080"),
	}
}

func main() {
	// Load server configuration
	serverConfig := loadServerConfig()

	// Initialize application context with dependencies
	appContext, err := di.SetupDependencies()
	if err != nil {
		panic(fmt.Errorf("error initializing application context: %w", err))
	}

	// Setup router
	router := setupRouter(appContext)

	// Setup server
	server := setupServer(router, serverConfig.Port)

	// Start server
	fmt.Printf("Server running at http://localhost:%s\n", serverConfig.Port)
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

// Helper function
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
