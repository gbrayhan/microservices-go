package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gin-gonic/gin"
)

func ApplicationRouter(router *gin.Engine, appContext *di.ApplicationContext) {
	v1 := router.Group("/v1")

	AuthRoutes(v1, appContext.AuthController)
	UserRoutes(v1, appContext.UserController)
	MedicineRoutes(v1, appContext.MedicineController)
}
