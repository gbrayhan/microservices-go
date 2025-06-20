package user

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainUser "github.com/gbrayhan/microservices-go/src/domain/user"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
)

// Structures
type NewUserRequest struct {
	UserName  string `json:"user" binding:"required"`
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Role      string `json:"role" binding:"required"`
}

type ResponseUser struct {
	ID        int       `json:"id"`
	UserName  string    `json:"user"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type IUserController interface {
	NewUser(ctx *gin.Context)
	GetAllUsers(ctx *gin.Context)
	GetUsersByID(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
}

type UserController struct {
	userService domainUser.IUserService
}

func NewUserController(userService domainUser.IUserService) IUserController {
	return &UserController{userService}
}

func (c *UserController) NewUser(ctx *gin.Context) {
	var request NewUserRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	userModel, err := c.userService.Create(toUsecaseMapper(&request))
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	userResponse := domainToResponseMapper(userModel)
	ctx.JSON(http.StatusOK, userResponse)
}

func (c *UserController) GetAllUsers(ctx *gin.Context) {
	users, err := c.userService.GetAll()
	if err != nil {
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	ctx.JSON(http.StatusOK, arrayDomainToResponseMapper(users))
}

func (c *UserController) GetUsersByID(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		appError := domainErrors.NewAppError(errors.New("user id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	user, err := c.userService.GetByID(userID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToResponseMapper(user))
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	var requestMap map[string]any
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
	userUpdated, err := c.userService.Update(userID, requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, domainToResponseMapper(userUpdated))
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err = c.userService.Delete(userID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}

// Mappers
func domainToResponseMapper(domainUser *domainUser.User) *ResponseUser {
	return &ResponseUser{
		ID:        domainUser.ID,
		UserName:  domainUser.UserName,
		Email:     domainUser.Email,
		FirstName: domainUser.FirstName,
		LastName:  domainUser.LastName,
		Status:    domainUser.Status,
		CreatedAt: domainUser.CreatedAt,
		UpdatedAt: domainUser.UpdatedAt,
	}
}

func arrayDomainToResponseMapper(users *[]domainUser.User) *[]ResponseUser {
	res := make([]ResponseUser, len(*users))
	for i, u := range *users {
		res[i] = *domainToResponseMapper(&u)
	}
	return &res
}

func toUsecaseMapper(req *NewUserRequest) *domainUser.User {
	return &domainUser.User{
		UserName:  req.UserName,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
	}
}
