// Package user contains the user controller
package user

import (
	"errors"
	useCaseUser "github.com/gbrayhan/microservices-go/application/usecases/user"
	domainErrors "github.com/gbrayhan/microservices-go/domain/errors"
	"github.com/gbrayhan/microservices-go/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// Controller is a struct that contains the user service
type Controller struct {
	UserService useCaseUser.Service
}

// NewUser godoc
// @Tags user
// @Summary Create New UserName
// @Description Create new user on the system
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param data body NewUserRequest true "body data"
// @Success 200 {object} ResponseUser
// @Failure 400 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /user [post]
func (c *Controller) NewUser(ctx *gin.Context) {
	var request NewUserRequest

	if err := controllers.BindJSON(ctx, &request); err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	userModel, err := c.UserService.Create(toUsecaseMapper(&request))
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	userResponse := domainToResponseMapper(userModel)
	ctx.JSON(http.StatusOK, userResponse)
}

// GetAllUsers godoc
// @Tags user
// @Summary Get all Users
// @Security ApiKeyAuth
// @Description Get all Users on the system
// @Success 200 {object} []ResponseUser
// @Failure 400 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /user [get]
func (c *Controller) GetAllUsers(ctx *gin.Context) {
	users, err := c.UserService.GetAll()
	if err != nil {
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}

	ctx.JSON(http.StatusOK, arrayDomainToResponseMapper(users))
}

// GetUsersByID godoc
// @Tags user
// @Summary Get users by ID
// @Description Get Users by ID on the system
// @Param user_id path int true "id of user"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseUser
// @Failure 400 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /user/{user_id} [get]
func (c *Controller) GetUsersByID(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		appError := domainErrors.NewAppError(errors.New("user id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	user, err := c.UserService.GetByID(userID)
	if err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	ctx.JSON(http.StatusOK, domainToResponseMapper(user))
}

// UpdateUser godoc
// @Tags user
// @Summary Get users by ID
// @Description Get Users by ID on the system
// @Param user_id path int true "id of user"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseUser
// @Failure 400 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /user/{user_id} [get]
func (c *Controller) UpdateUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		appError := domainErrors.NewAppError(errors.New("param id is necessary in the url"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	var requestMap map[string]interface{}

	err = controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err = updateValidation(requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	user, err := c.UserService.Update(userID, requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, domainToResponseMapper(user))
}

// DeleteUser godoc
// @Tags user
// @Summary Get users by ID
// @Description Get Users by ID on the system
// @Param user_id path int true "id of user"
// @Security ApiKeyAuth
// @Success 200 {object} MessageResponse
// @Failure 400 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /user/{user_id} [get]
func (c *Controller) DeleteUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		appError := domainErrors.NewAppError(errors.New("param id is necessary in the url"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	err = c.UserService.Delete(userID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})

}
