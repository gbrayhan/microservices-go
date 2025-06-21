package routes

import (
	"net/http"

	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gin-gonic/gin"
)

func ApplicationRouter(router *gin.Engine, appContext *di.ApplicationContext) {
	v1 := router.Group("/v1")

	v1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Service is running",
		})
	})

	AuthRoutes(v1, appContext.AuthController)
	UserRoutes(v1, appContext.UserController)
	MedicineRoutes(v1, appContext.MedicineController)
}
