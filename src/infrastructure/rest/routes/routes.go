// Package routes contains all routes of the application
package routes

import (
	// swaggerFiles for documentation
	_ "github.com/gbrayhan/microservices-go/docs"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/adapter"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// Security is a struct that contains the security of the application
// @SecurityDefinitions.jwt
type Security struct {
	Authorization string `header:"Authorization" json:"Authorization"`
}

// @title Boilerplate Golang
// @version 1.2
// @description Documentation's Boilerplate Golang
// @termsOfService http://swagger.io/terms/

// @contact.name Alejandro Gabriel Guerrero
// @contact.url http://github.com/gbrayhan
// @contact.email gbrayhan@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// ApplicationV1Router is a function that contains all routes of the application
// @host localhost:8080
// @BasePath /v1
func ApplicationV1Router(router *gin.Engine, db *gorm.DB) {
	routerV1 := router.Group("/v1")

	{
		// Documentation Swagger
		{
			routerV1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}

		AuthRoutes(routerV1, adapter.AuthAdapter(db))
		UserRoutes(routerV1, adapter.UserAdapter(db))
		MedicineRoutes(routerV1, adapter.MedicineAdapter(db))
	}
}
