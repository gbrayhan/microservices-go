package user

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainUser "github.com/gbrayhan/microservices-go/src/domain/user"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	SearchPaginated(ctx *gin.Context)
	SearchByProperty(ctx *gin.Context)
}

type UserController struct {
	userService domainUser.IUserService
	Logger      *logger.Logger
}

func NewUserController(userService domainUser.IUserService, loggerInstance *logger.Logger) IUserController {
	return &UserController{userService: userService, Logger: loggerInstance}
}

func (c *UserController) NewUser(ctx *gin.Context) {
	c.Logger.Info("Creating new user")
	var request NewUserRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new user", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	userModel, err := c.userService.Create(toUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating user", zap.Error(err), zap.String("email", request.Email))
		_ = ctx.Error(err)
		return
	}
	userResponse := domainToResponseMapper(userModel)
	c.Logger.Info("User created successfully", zap.String("email", request.Email), zap.Int("id", userModel.ID))
	ctx.JSON(http.StatusOK, userResponse)
}

func (c *UserController) GetAllUsers(ctx *gin.Context) {
	c.Logger.Info("Getting all users")
	users, err := c.userService.GetAll()
	if err != nil {
		c.Logger.Error("Error getting all users", zap.Error(err))
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Successfully retrieved all users", zap.Int("count", len(*users)))
	ctx.JSON(http.StatusOK, arrayDomainToResponseMapper(users))
}

func (c *UserController) GetUsersByID(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid user ID parameter", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("user id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Getting user by ID", zap.Int("id", userID))
	user, err := c.userService.GetByID(userID)
	if err != nil {
		c.Logger.Error("Error getting user by ID", zap.Error(err), zap.Int("id", userID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Successfully retrieved user by ID", zap.Int("id", userID))
	ctx.JSON(http.StatusOK, domainToResponseMapper(user))
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid user ID parameter for update", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Updating user", zap.Int("id", userID))
	var requestMap map[string]any
	err = controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		c.Logger.Error("Error binding JSON for user update", zap.Error(err), zap.Int("id", userID))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err = updateValidation(requestMap)
	if err != nil {
		c.Logger.Error("Validation error for user update", zap.Error(err), zap.Int("id", userID))
		_ = ctx.Error(err)
		return
	}
	userUpdated, err := c.userService.Update(userID, requestMap)
	if err != nil {
		c.Logger.Error("Error updating user", zap.Error(err), zap.Int("id", userID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("User updated successfully", zap.Int("id", userID))
	ctx.JSON(http.StatusOK, domainToResponseMapper(userUpdated))
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid user ID parameter for deletion", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Deleting user", zap.Int("id", userID))
	err = c.userService.Delete(userID)
	if err != nil {
		c.Logger.Error("Error deleting user", zap.Error(err), zap.Int("id", userID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("User deleted successfully", zap.Int("id", userID))
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}

func (c *UserController) SearchPaginated(ctx *gin.Context) {
	c.Logger.Info("Searching users with pagination")

	// Parse query parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if pageSize < 1 {
		pageSize = 10
	}

	// Build filters
	filters := domain.DataFilters{
		Page:     page,
		PageSize: pageSize,
	}

	// Parse like filters
	likeFilters := make(map[string][]string)
	for field := range user.ColumnsUserMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	// Parse exact matches
	matches := make(map[string][]string)
	for field := range user.ColumnsUserMapping {
		if values := ctx.QueryArray(field + "_match"); len(values) > 0 {
			matches[field] = values
		}
	}
	filters.Matches = matches

	// Parse date range filters
	var dateRanges []domain.DateRangeFilter
	for field := range user.ColumnsUserMapping {
		startStr := ctx.Query(field + "_start")
		endStr := ctx.Query(field + "_end")

		if startStr != "" || endStr != "" {
			dateRange := domain.DateRangeFilter{Field: field}

			if startStr != "" {
				if startTime, err := time.Parse(time.RFC3339, startStr); err == nil {
					dateRange.Start = &startTime
				}
			}

			if endStr != "" {
				if endTime, err := time.Parse(time.RFC3339, endStr); err == nil {
					dateRange.End = &endTime
				}
			}

			dateRanges = append(dateRanges, dateRange)
		}
	}
	filters.DateRangeFilters = dateRanges

	// Parse sorting
	sortBy := ctx.QueryArray("sortBy")
	if len(sortBy) > 0 {
		filters.SortBy = sortBy
	}

	sortDirection := domain.SortDirection(ctx.DefaultQuery("sortDirection", "asc"))
	if sortDirection.IsValid() {
		filters.SortDirection = sortDirection
	}

	result, err := c.userService.SearchPaginated(filters)
	if err != nil {
		c.Logger.Error("Error searching users", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	response := gin.H{
		"data":       arrayDomainToResponseMapper(result.Data),
		"total":      result.Total,
		"page":       result.Page,
		"pageSize":   result.PageSize,
		"totalPages": result.TotalPages,
		"filters":    filters,
	}

	c.Logger.Info("Successfully searched users",
		zap.Int64("total", result.Total),
		zap.Int("page", result.Page))
	ctx.JSON(http.StatusOK, response)
}

func (c *UserController) SearchByProperty(ctx *gin.Context) {
	property := ctx.Query("property")
	searchText := ctx.Query("searchText")

	if property == "" || searchText == "" {
		c.Logger.Error("Missing property or searchText parameter")
		appError := domainErrors.NewAppError(errors.New("missing property or searchText parameter"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	// Validate property
	allowed := map[string]bool{
		"userName":  true,
		"email":     true,
		"firstName": true,
		"lastName":  true,
		"status":    true,
	}
	if !allowed[property] {
		c.Logger.Error("Invalid property for search", zap.String("property", property))
		appError := domainErrors.NewAppError(errors.New("invalid property"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	coincidences, err := c.userService.SearchByProperty(property, searchText)
	if err != nil {
		c.Logger.Error("Error searching by property", zap.Error(err), zap.String("property", property))
		_ = ctx.Error(err)
		return
	}

	c.Logger.Info("Successfully searched by property",
		zap.String("property", property),
		zap.Int("results", len(*coincidences)))
	ctx.JSON(http.StatusOK, coincidences)
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
