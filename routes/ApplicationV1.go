package routes

import (
	"github.com/banwire/microservice_golang/controllers"
	"github.com/gin-gonic/gin"
)

func ApplicationV1Router(router *gin.Engine) {

	{
		router.POST("/example", controllers.ExampleAction)
	}

}
