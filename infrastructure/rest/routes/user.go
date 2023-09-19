// Package routes contains all routes of the application
package routes

import (
	userController "github.com/gbrayhan/microservices-go/infrastructure/rest/controllers/user"
	"github.com/gbrayhan/microservices-go/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

// UserRoutes is a function that contains all routes of the user
func UserRoutes(router *gin.RouterGroup, controller *userController.Controller) {
	routerUser := router.Group("/user")
	routerUser.Use(middlewares.AuthJWTMiddleware())
	{
		routerUser.POST("/", controller.NewUser)
		routerUser.GET("/:id", controller.GetUsersByID)
		routerUser.GET("/", controller.GetAllUsers)
		routerUser.PUT("/:id", controller.UpdateUser)
		routerUser.DELETE("/:id", controller.DeleteUser)
	}
}
