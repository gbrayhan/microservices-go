package routes

import (
	userController "github.com/gbrayhan/microservices-go/infrastructure/rest/controllers/user"
	"github.com/gbrayhan/microservices-go/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup, controller *userController.Controller) {
	routerAuth := router.Group("/user")
	routerAuth.Use(middlewares.AuthJWTMiddleware())
	{
		routerAuth.POST("/", controller.NewUser)
		routerAuth.GET("/:id", controller.GetUsersByID)
		routerAuth.GET("/", controller.GetAllUsers)
		routerAuth.PUT("/:id", controller.UpdateUser)
		routerAuth.DELETE("/:id", controller.DeleteUser)
	}
}
