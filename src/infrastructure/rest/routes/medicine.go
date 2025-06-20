package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/medicine"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func MedicineRoutes(router *gin.RouterGroup, controller medicine.IMedicineController) {
	med := router.Group("/medicine")
	med.Use(middlewares.AuthJWTMiddleware())
	{
		med.GET("/", controller.GetAllMedicines)
		med.POST("/", controller.NewMedicine)
		med.GET("/:id", controller.GetMedicinesByID)
		med.PUT("/:id", controller.UpdateMedicine)
		med.DELETE("/:id", controller.DeleteMedicine)
	}
}
