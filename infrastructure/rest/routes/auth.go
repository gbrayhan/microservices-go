// Package routes contains all routes of the application
package routes

import (
	authController "github.com/gbrayhan/microservices-go/infrastructure/rest/controllers/auth"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.RouterGroup, controller *authController.Controller) {

	routerAuth := router.Group("/auth")
	{
		routerAuth.POST("/login", controller.Login)
	}

}
