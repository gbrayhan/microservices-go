package routes

import (
	"github.com/gbrayhan/microservices-go/controllers"
	"github.com/gbrayhan/microservices-go/middlewares"
	"github.com/gin-gonic/gin"
)

func ApplicationV1Router(router *gin.Engine) {
	router.Use(middlewares.GinBodyLogMiddleware)

	{
		router.POST("/example", controllers.ExampleAction)
	}
}
