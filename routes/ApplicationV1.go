package routes

import (
	"github.com/gbrayhan/microservices-go/controllers"
	"github.com/gin-gonic/gin"
)

func ApplicationV1Router(router *gin.Engine) {
	{
		router.POST("/example", controllers.ExampleAction)
	}
}
