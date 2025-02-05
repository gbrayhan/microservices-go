package routes

import (
	authController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/auth"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.RouterGroup, controller authController.IAuthController) {
	routerAuth := router.Group("/auth")
	{
		routerAuth.POST("/login", controller.Login)
		routerAuth.POST("/access-token", controller.GetAccessTokenByRefreshToken)
	}
}
