// Package routes contains all routes of the application
package routes

import (
	medicineController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/medicine"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

// MedicineRoutes is a function that contains all medicine routes
func MedicineRoutes(router *gin.RouterGroup, controller *medicineController.Controller) {

	routerMedicine := router.Group("/medicine")
	routerMedicine.Use(middlewares.AuthJWTMiddleware())

	{
		routerMedicine.GET("/", controller.GetAllMedicines)
		routerMedicine.POST("/", controller.NewMedicine)
		routerMedicine.GET("/:id", controller.GetMedicinesByID)
		routerMedicine.POST("/data", controller.GetDataMedicines)
		routerMedicine.PUT("/:id", controller.UpdateMedicine)
		routerMedicine.DELETE("/:id", controller.DeleteMedicine)

	}

}
