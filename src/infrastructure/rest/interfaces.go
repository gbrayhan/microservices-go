package rest

import (
	"github.com/gin-gonic/gin"
)

// AuthControllerInterface defines the interface for auth controller operations
type AuthControllerInterface interface {
	Login(c *gin.Context)
	RefreshToken(c *gin.Context)
}

// UserControllerInterface defines the interface for user controller operations
type UserControllerInterface interface {
	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

// MedicineControllerInterface defines the interface for medicine controller operations
type MedicineControllerInterface interface {
	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetData(c *gin.Context)
}

// MiddlewareInterface defines the interface for middleware operations
type MiddlewareInterface interface {
	Handle(c *gin.Context)
}

// RouterInterface defines the interface for router operations
type RouterInterface interface {
	SetupRoutes(router *gin.Engine)
}
