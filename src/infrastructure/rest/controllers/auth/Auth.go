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

	domainUser, authTokens, err := c.authUsecase.Login(request.Email, request.Password)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	response := LoginResponse{
		Data: UserData{
			UserName:  domainUser.UserName,
			Email:     domainUser.Email,
			FirstName: domainUser.FirstName,
			LastName:  domainUser.LastName,
			Status:    domainUser.Status,
			ID:        domainUser.ID,
		},
		Security: SecurityData{
			JWTAccessToken:            authTokens.AccessToken,
			JWTRefreshToken:           authTokens.RefreshToken,
			ExpirationAccessDateTime:  authTokens.ExpirationAccessDateTime,
			ExpirationRefreshDateTime: authTokens.ExpirationRefreshDateTime,
		},
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *AuthController) GetAccessTokenByRefreshToken(ctx *gin.Context) {
	var request AccessTokenRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	domainUser, authTokens, err := c.authUsecase.AccessTokenByRefreshToken(request.RefreshToken)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	response := LoginResponse{
		Data: UserData{
			UserName:  domainUser.UserName,
			Email:     domainUser.Email,
			FirstName: domainUser.FirstName,
			LastName:  domainUser.LastName,
			Status:    domainUser.Status,
			ID:        domainUser.ID,
		},
		Security: SecurityData{
			JWTAccessToken:            authTokens.AccessToken,
			JWTRefreshToken:           authTokens.RefreshToken,
			ExpirationAccessDateTime:  authTokens.ExpirationAccessDateTime,
			ExpirationRefreshDateTime: authTokens.ExpirationRefreshDateTime,
		},
	}

	ctx.JSON(http.StatusOK, response)
}
