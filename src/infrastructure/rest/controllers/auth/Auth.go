package auth

import (
	"net/http"

	useCaseAuth "github.com/gbrayhan/microservices-go/src/application/usecases/auth"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IAuthController interface {
	Login(ctx *gin.Context)
	GetAccessTokenByRefreshToken(ctx *gin.Context)
}

type AuthController struct {
	authUseCase useCaseAuth.IAuthUseCase
	Logger      *logger.Logger
}

func NewAuthController(authUsecase useCaseAuth.IAuthUseCase, loggerInstance *logger.Logger) IAuthController {
	return &AuthController{
		authUseCase: authUsecase,
		Logger:      loggerInstance,
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	c.Logger.Info("User login request")
	var request LoginRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for login", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	domainUser, authTokens, err := c.authUseCase.Login(request.Email, request.Password)
	if err != nil {
		c.Logger.Error("Login failed", zap.Error(err), zap.String("email", request.Email))
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

	c.Logger.Info("Login successful", zap.String("email", request.Email), zap.Int("userID", domainUser.ID))
	ctx.JSON(http.StatusOK, response)
}

func (c *AuthController) GetAccessTokenByRefreshToken(ctx *gin.Context) {
	c.Logger.Info("Token refresh request")
	var request AccessTokenRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for token refresh", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	domainUser, authTokens, err := c.authUseCase.AccessTokenByRefreshToken(request.RefreshToken)
	if err != nil {
		c.Logger.Error("Token refresh failed", zap.Error(err))
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

	c.Logger.Info("Token refresh successful", zap.Int("userID", domainUser.ID))
	ctx.JSON(http.StatusOK, response)
}
