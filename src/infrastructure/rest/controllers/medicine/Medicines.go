package medicine

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainError "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainMedicine "github.com/gbrayhan/microservices-go/src/domain/medicine"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
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
	GetDataMedicines(ctx *gin.Context)
	GetMedicinesByID(ctx *gin.Context)
	UpdateMedicine(ctx *gin.Context)
	DeleteMedicine(ctx *gin.Context)
}

type Controller struct {
	MedicineService domainMedicine.IMedicineService
}

func NewMedicineController(medicineService domainMedicine.IMedicineService) IMedicineController {
	return &Controller{MedicineService: medicineService}
}

func (c *Controller) NewMedicine(ctx *gin.Context) {
	var request NewMedicineRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		appError := domainError.NewAppError(err, domainError.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	newMed := domainMedicine.Medicine{
		Name:        request.Name,
		Description: request.Description,
		Laboratory:  request.Laboratory,
		EanCode:     request.EanCode,
	}
	dMed, err := c.MedicineService.Create(&newMed)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	resp := domainToResponseMapper(dMed)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) GetAllMedicines(ctx *gin.Context) {
	medicines, err := c.MedicineService.GetAll()
	if err != nil {
		appError := domainError.NewAppErrorWithType(domainError.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	ctx.JSON(http.StatusOK, arrayDomainToResponseMapper(medicines))
}

func (c *Controller) GetDataMedicines(ctx *gin.Context) {
	var request DataMedicineRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		appError := domainError.NewAppError(err, domainError.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	var dateRangeFiltersDomain []domain.DateRangeFilter
	for _, f := range request.FieldsDateRange {
		dateRangeFiltersDomain = append(dateRangeFiltersDomain, domain.DateRangeFilter{
			Field: f.Field, Start: f.StartDate, End: f.EndDate,
		})
	}
	result, err := c.MedicineService.GetData(
		request.Page, request.Limit,
		request.SorBy.Field, request.SorBy.Direction,
		request.Filters, request.GlobalSearch,
		dateRangeFiltersDomain,
	)
	if err != nil {
		appError := domainError.NewAppErrorWithType(domainError.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	numPages, nextCursor, prevCursor := controllers.PaginationValues(request.Limit, request.Page, result.Total)
	resp := PaginationResultMedicine{
		Data:       arrayDomainToResponseMapper(result.Data),
		Total:      result.Total,
		Limit:      request.Limit,
		Current:    request.Page,
		NextCursor: nextCursor,
		PrevCursor: prevCursor,
		NumPages:   numPages,
	}
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) GetMedicinesByID(ctx *gin.Context) {
	medicineID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		appError := domainError.NewAppError(errors.New("medicine id is invalid"), domainError.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	dMed, err := c.MedicineService.GetByID(medicineID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dMed)
}

func (c *Controller) UpdateMedicine(ctx *gin.Context) {
	medicineID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		appError := domainError.NewAppError(errors.New("param id is necessary"), domainError.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	var requestMap map[string]any
	if err := controllers.BindJSONMap(ctx, &requestMap); err != nil {
		appError := domainError.NewAppError(err, domainError.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err = updateValidation(requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	updated, err := c.MedicineService.Update(medicineID, requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, updated)
}

func (c *Controller) DeleteMedicine(ctx *gin.Context) {
	medicineID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		appError := domainError.NewAppError(errors.New("param id is necessary"), domainError.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	if err = c.MedicineService.Delete(medicineID); err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}

// Mappers
func domainToResponseMapper(m *domainMedicine.Medicine) *ResponseMedicine {
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

func arrayDomainToResponseMapper(m *[]domainMedicine.Medicine) *[]ResponseMedicine {
	res := make([]ResponseMedicine, len(*m))
	for i, med := range *m {
		res[i] = *domainToResponseMapper(&med)
	}
	return &res
}
