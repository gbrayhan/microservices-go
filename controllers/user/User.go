package user

import (
  "errors"
  "github.com/gbrayhan/microservices-go/controllers"
  usecases "github.com/gbrayhan/microservices-go/usecases/user"
  "github.com/gin-gonic/gin"
  "net/http"
  "strconv"

  _ "github.com/gbrayhan/microservices-go/controllers/errors"
  errorModels "github.com/gbrayhan/microservices-go/models/errors"
  model "github.com/gbrayhan/microservices-go/models/user"
)

// NewUser godoc
// @Tags user
// @Summary Create New User
// @Description Create new user on the system
// @Accept  json
// @Produce  json
// @Param data body NewUserRequest true "body data"
// @Success 200 {object} model.User
// @Failure 400 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /user [post]
func NewUser(c *gin.Context) {
  var request NewUserRequest

  if err := controllers.BindJSON(c, &request); err != nil {
    appError := errorModels.NewAppError(err, errorModels.ValidationError)
    _ = c.Error(appError)
    return
  }

  user := usecases.NewUser{
    User:      request.User,
    Email:     request.Email,
    FirstName: request.FirstName,
    LastName:  request.LastName,
    Password:  request.Password,
  }
  userModel, err := usecases.CreateUser(&user)
  if err != nil {
    _ = c.Error(err)
    return
  }
  userResponse := UserModelToResponseMapper(userModel)
  c.JSON(http.StatusOK, userResponse)
}

// GetAllUsers godoc
// @Tags user
// @Summary Get all Users
// @Description Get all Users on the system
// @Success 200 {object} []model.User
// @Failure 400 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /user [get]
func GetAllUsers(c *gin.Context) {
  users, err := model.GetAllUsers()
  if err != nil {
    appError := errorModels.NewAppErrorWithType(errorModels.UnknownError)
    _ = c.Error(appError)
    return
  }
  var usersResponse []UserResponse
  for _, item := range users {
    usersResponse = append(usersResponse, *UserModelToResponseMapper(item))
  }
  c.JSON(http.StatusOK, usersResponse)
}

// GetUsersByID godoc
// @Tags user
// @Summary Get users by ID
// @Description Get Users by ID on the system
// @Param user_id path int true "id of user"
// @Success 200 {object} model.User
// @Failure 400 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /user/{user_id} [get]
func GetUsersByID(c *gin.Context) {
  var user model.User
  userID, err := strconv.Atoi(c.Param("id"))
  if err != nil {
    appError := errorModels.NewAppError(errors.New("user id is invalid"), errorModels.ValidationError)
    _ = c.Error(appError)
    return
  }

  user, err = model.GetUserByID(userID)
  if err != nil {
    appError := errorModels.NewAppError(err, errorModels.ValidationError)
    _ = c.Error(appError)
    return
  }

  c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
  userID, err := strconv.Atoi(c.Param("id"))
  if err != nil {
    appError := errorModels.NewAppError(errors.New("param id is necessary in the url"), errorModels.ValidationError)
    _ = c.Error(appError)
    return
  }
  var requestMap map[string]interface{}

  err = controllers.BindJSONMap(c, &requestMap)
  if err != nil {
    appError := errorModels.NewAppError(err, errorModels.ValidationError)
    _ = c.Error(appError)
    return
  }

  err = updateValidation(requestMap)
  if err != nil {
    _ = c.Error(err)
    return
  }

  user, err := model.UpdateUser(userID, requestMap)
  if err != nil {
    _ = c.Error(err)
    return
  }

  c.JSON(http.StatusOK, user)

}

func DeleteUser(c *gin.Context) {
  userID, err := strconv.Atoi(c.Param("id"))
  if err != nil {
    appError := errorModels.NewAppError(errors.New("param id is necessary in the url"), errorModels.ValidationError)
    _ = c.Error(appError)
    return
  }

  err = model.DeleteUser(userID)
  if err != nil {
    _ = c.Error(err)
    return
  }
  c.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})

}

func Login(c *gin.Context) {
  var request LoginRequest

  if err := controllers.BindJSON(c, &request); err != nil {
    appError := errorModels.NewAppError(err, errorModels.ValidationError)
    _ = c.Error(appError)
    return
  }
  user := usecases.LoginUser{
    Email:    request.Email,
    Password: request.Password,
  }

  usecases.Login(user)
}
