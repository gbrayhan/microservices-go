package routes

import (
  "github.com/gbrayhan/microservices-go/controllers/medicine"
  "github.com/gbrayhan/microservices-go/controllers/user"
  "github.com/gin-gonic/gin"
  swaggerFiles "github.com/swaggo/files"
  ginSwagger "github.com/swaggo/gin-swagger"

  _ "github.com/gbrayhan/microservices-go/docs"
)

// @title Boilerplate Golang
// @version 1.0
// @description Documentation's Boilerplate Golang
// @termsOfService http://swagger.io/terms/

// @contact.name Alejandro Gabriel Guerrero
// @contact.url http://github.com/gbrayhan
// @contact.email gbrayhan@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /v1
func ApplicationV1Router(router *gin.Engine) {
  v1 := router.Group("/v1")
  {
    // Documentation Swagger
    {
      v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    }
    // Medicines
    v1Medicines := v1.Group("/medicine")
    {
      v1Medicines.POST("/", medicine.NewMedicine)
      v1Medicines.GET("/:id", medicine.GetMedicinesByID)
      v1Medicines.GET("/", medicine.GetAllMedicines)
      v1Medicines.PUT("/:id", medicine.UpdateMedicine)
      v1Medicines.DELETE("/:id", medicine.DeleteMedicine)
    }

    // Users
    v1User := v1.Group("/user")
    {
      v1User.POST("/", user.NewUser)
      v1User.GET("/:id", user.GetUsersByID)
      v1User.GET("/", user.GetAllUsers)
      v1User.PUT("/:id", user.UpdateUser)
      v1User.DELETE("/:id", user.DeleteUser)
      v1User.POST("/login", user.Login)
    }
  }
}
