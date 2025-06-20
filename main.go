package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	// Initialize logger first based on environment
	env := getEnvOrDefault("GO_ENV", "development")
	var loggerInstance *logger.Logger
	var err error

	if env == "development" {
		loggerInstance, err = logger.NewDevelopmentLogger()
	} else {
		loggerInstance, err = logger.NewLogger()
	}

	if err != nil {
		panic(fmt.Errorf("error initializing logger: %w", err))
	}
	defer func() {
		if err := loggerInstance.Log.Sync(); err != nil {
			loggerInstance.Log.Error("Failed to sync logger", zap.Error(err))
		}
	}()

	loggerInstance.Info("Starting microservices application")

	// Load server configuration
	serverConfig := loadServerConfig()

	// Initialize application context with dependencies and logger
	appContext, err := di.SetupDependencies(loggerInstance)
	if err != nil {
		loggerInstance.Panic("Error initializing application context", zap.Error(err))
	}

	// Setup router
	router := setupRouter(appContext, loggerInstance)

	// Setup server
	server := setupServer(router, serverConfig.Port)

	// Start server
	loggerInstance.Info("Server starting", zap.String("port", serverConfig.Port))
	if err := server.ListenAndServe(); err != nil {
		loggerInstance.Panic("Server failed to start", zap.Error(err))
	}
}

func setupRouter(appContext *di.ApplicationContext, logger *logger.Logger) *gin.Engine {
	// Configurar Gin para usar el logger de Zap basado en el entorno
	env := getEnvOrDefault("GO_ENV", "development")
	if env == "development" {
		logger.SetupGinWithZapLoggerInDevelopment()
	} else {
		logger.SetupGinWithZapLogger()
	}

	// Crear el router después de configurar el logger
	router := gin.New()

	// Agregar middlewares de recuperación y logger personalizados
	router.Use(gin.Recovery())
	router.Use(cors.Default())

	// Add middlewares
	router.Use(middlewares.ErrorHandler())
	router.Use(middlewares.GinBodyLogMiddleware)
	router.Use(middlewares.CommonHeaders)

	// Add logger middleware
	router.Use(logger.GinZapLogger())

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
