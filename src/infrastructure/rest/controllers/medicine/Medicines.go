package medicine

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	domainError "github.com/gbrayhan/microservices-go/src/domain/errors"
	medicineDomain "github.com/gbrayhan/microservices-go/src/domain/medicine"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Structures
type NewMedicineRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Laboratory  string `json:"laboratory" binding:"required"`
	EanCode     string `json:"eanCode" binding:"required"`
}

type DataMedicineRequest struct {
	Limit           int64                                   `json:"limit" example:"10"`
	Page            int64                                   `json:"page" example:"1"`
	GlobalSearch    string                                  `json:"globalSearch" example:"John"`
	Filters         map[string][]string                     `json:"filters"`
	SorBy           controllers.SortByDataRequest           `json:"sortBy"`
	FieldsDateRange []controllers.FieldDateRangeDataRequest `json:"fieldsDateRange"`
}

type ResponseMedicine struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	EanCode     string    `json:"eanCode"`
	Laboratory  string    `json:"laboratory"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}

type PaginationResultMedicine struct {
	Data       *[]ResponseMedicine `json:"data"`
	Total      int64               `json:"total"`
	Limit      int64               `json:"limit"`
	Current    int64               `json:"current"`
	NextCursor int64               `json:"nextCursor"`
	PrevCursor int64               `json:"prevCursor"`
	NumPages   int64               `json:"numPages"`
}

type IMedicineController interface {
	NewMedicine(ctx *gin.Context)
	GetAllMedicines(ctx *gin.Context)
	GetMedicinesByID(ctx *gin.Context)
	UpdateMedicine(ctx *gin.Context)
	DeleteMedicine(ctx *gin.Context)
}

type Controller struct {
	medicineService medicineDomain.IMedicineService
	Logger          *logger.Logger
}

func NewMedicineController(medicineService medicineDomain.IMedicineService, loggerInstance *logger.Logger) IMedicineController {
	return &Controller{medicineService: medicineService, Logger: loggerInstance}
}

func (c *Controller) NewMedicine(ctx *gin.Context) {
	c.Logger.Info("Creating new medicine")
	var request NewMedicineRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new medicine", zap.Error(err))
		appError := domainError.NewAppError(err, domainError.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	newMed := medicineDomain.Medicine{
		Name:        request.Name,
		Description: request.Description,
		Laboratory:  request.Laboratory,
		EanCode:     request.EanCode,
	}
	dMed, err := c.medicineService.Create(&newMed)
	if err != nil {
		c.Logger.Error("Error creating medicine", zap.Error(err), zap.String("name", request.Name))
		_ = ctx.Error(err)
		return
	}
	resp := domainToResponseMapper(dMed)
	c.Logger.Info("Medicine created successfully", zap.String("name", request.Name), zap.Int("id", dMed.ID))
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) GetAllMedicines(ctx *gin.Context) {
	c.Logger.Info("Getting all medicines")
	medicines, err := c.medicineService.GetAll()
	if err != nil {
		c.Logger.Error("Error getting all medicines", zap.Error(err))
		appError := domainError.NewAppErrorWithType(domainError.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Successfully retrieved all medicines", zap.Int("count", len(*medicines)))
	ctx.JSON(http.StatusOK, arrayDomainToResponseMapper(medicines))
}

func (c *Controller) GetMedicinesByID(ctx *gin.Context) {
	medicineID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid medicine ID parameter", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainError.NewAppError(errors.New("medicine id is invalid"), domainError.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Getting medicine by ID", zap.Int("id", medicineID))
	dMed, err := c.medicineService.GetByID(medicineID)
	if err != nil {
		c.Logger.Error("Error getting medicine by ID", zap.Error(err), zap.Int("id", medicineID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Successfully retrieved medicine by ID", zap.Int("id", medicineID))
	ctx.JSON(http.StatusOK, dMed)
}

func (c *Controller) UpdateMedicine(ctx *gin.Context) {
	medicineID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid medicine ID parameter for update", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainError.NewAppError(errors.New("param id is necessary"), domainError.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Updating medicine", zap.Int("id", medicineID))
	var requestMap map[string]any
	if err := controllers.BindJSONMap(ctx, &requestMap); err != nil {
		c.Logger.Error("Error binding JSON for medicine update", zap.Error(err), zap.Int("id", medicineID))
		appError := domainError.NewAppError(err, domainError.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err = updateValidation(requestMap)
	if err != nil {
		c.Logger.Error("Validation error for medicine update", zap.Error(err), zap.Int("id", medicineID))
		_ = ctx.Error(err)
		return
	}
	updated, err := c.medicineService.Update(medicineID, requestMap)
	if err != nil {
		c.Logger.Error("Error updating medicine", zap.Error(err), zap.Int("id", medicineID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Medicine updated successfully", zap.Int("id", medicineID))
	ctx.JSON(http.StatusOK, updated)
}

func (c *Controller) DeleteMedicine(ctx *gin.Context) {
	medicineID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid medicine ID parameter for deletion", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainError.NewAppError(errors.New("param id is necessary"), domainError.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Deleting medicine", zap.Int("id", medicineID))
	if err = c.medicineService.Delete(medicineID); err != nil {
		c.Logger.Error("Error deleting medicine", zap.Error(err), zap.Int("id", medicineID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Medicine deleted successfully", zap.Int("id", medicineID))
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}

// Mappers
func domainToResponseMapper(m *medicineDomain.Medicine) *ResponseMedicine {
	return &ResponseMedicine{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		EanCode:     m.EanCode,
		Laboratory:  m.Laboratory,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func arrayDomainToResponseMapper(m *[]medicineDomain.Medicine) *[]ResponseMedicine {
	res := make([]ResponseMedicine, len(*m))
	for i, medicine := range *m {
		res[i] = *domainToResponseMapper(&medicine)
	}
	return &res
}
