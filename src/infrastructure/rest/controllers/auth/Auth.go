package auth

import (
	"net/http"

	useCaseAuth "github.com/gbrayhan/microservices-go/src/application/usecases/auth"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
)

type IAuthController interface {
	Login(ctx *gin.Context)
	GetAccessTokenByRefreshToken(ctx *gin.Context)
}

type AuthController struct {
	authUsecase useCaseAuth.IAuthUseCase
}

func NewAuthController(authUsecase useCaseAuth.IAuthUseCase) IAuthController {
	return &AuthController{
		authUsecase: authUsecase,
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var request LoginRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	userLogin := useCaseAuth.LoginUser{
		Email:    request.Email,
		Password: request.Password,
	}
	authDataUser, err := c.authUsecase.Login(userLogin)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, authDataUser)
}

func (c *AuthController) GetAccessTokenByRefreshToken(ctx *gin.Context) {
	var request AccessTokenRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	authDataUser, err := c.authUsecase.AccessTokenByRefreshToken(request.RefreshToken)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, authDataUser)
}
