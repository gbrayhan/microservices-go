package auth

import (
	useCaseAuth "github.com/gbrayhan/microservices-go/application/usecases/auth"
	domainErrors "github.com/gbrayhan/microservices-go/domain/errors"
	"github.com/gbrayhan/microservices-go/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	AuthService useCaseAuth.Service
}

// Login godoc
// @Tags auth
// @Summary Login UserName
// @Description Auth user by email and password
// @Param data body LoginRequest true "body data"
// @Success 200 {object} useCaseAuth.DataUserAuthenticated
// @Failure 400 {object} controllers.MessageResponse
// @Failure 500 {object} controllers.MessageResponse
// @Router /auth/login [post]
func (c *Controller) Login(ctx *gin.Context) {
	var request LoginRequest

	if err := controllers.BindJSON(ctx, &request); err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	user := useCaseAuth.LoginUser{
		Email:    request.Email,
		Password: request.Password,
	}

	authDataUser, err := c.AuthService.Login(user)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, authDataUser)
}
