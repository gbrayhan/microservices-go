package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/adapter"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func ApplicationRouter(router *gin.Engine, db *gorm.DB) {
	v1 := router.Group("/v1")

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	AuthRoutes(v1, adapter.AuthAdapter(db))
	UserRoutes(v1, adapter.UserAdapter(db))
	MedicineRoutes(v1, adapter.MedicineAdapter(db))
}
