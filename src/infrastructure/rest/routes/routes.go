package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ApplicationRouter(router *gin.Engine, appContext *di.ApplicationContext) {
	v1 := router.Group("/v1")

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	AuthRoutes(v1, appContext.AuthController)
	UserRoutes(v1, appContext.UserController)
	MedicineRoutes(v1, appContext.MedicineController)
}
