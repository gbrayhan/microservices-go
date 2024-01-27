// Package medicine contains the medicine controller
package medicine

import (
	"errors"
	"github.com/gbrayhan/microservices-go/src/domain"

	useCaseMedicine "github.com/gbrayhan/microservices-go/src/application/usecases/medicine"
	domainError "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainMedicine "github.com/gbrayhan/microservices-go/src/domain/medicine"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Controller is a struct that contains the medicine service
type Controller struct {
	MedicineService useCaseMedicine.Service
}

// NewMedicine godoc
// @Tags medicine
// @Summary Create New Medicine
// @Description Create new medicine on the system
// @Accept  json
// @Produce  json
// @Param data body NewMedicineRequest true "body data"
// @Success 200 {object} domainMedicine.Medicine
// @Failure 400 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /medicine [post]
func (c *Controller) NewMedicine(ctx *gin.Context) {
	var request NewMedicineRequest

	if err := controllers.BindJSON(ctx, &request); err != nil {
		appError := domainError.NewAppError(err, domainError.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	newMedicine := domainMedicine.NewMedicine{
		Name:        request.Name,
		Description: request.Description,
		Laboratory:  request.Laboratory,
		EANCode:     request.EanCode,
	}

	domainMedicine, err := c.MedicineService.Create(&newMedicine)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	responseMedicine := domainToResponseMapper(domainMedicine)

	ctx.JSON(http.StatusOK, responseMedicine)
}

// GetAllMedicines godoc
// @Tags medicine
// @Summary Get all Medicines
// @Description Get all Medicines on the system
// @Success 200 {object} []ResponseMedicine
// @Failure 400 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /medicine [get]
func (c *Controller) GetAllMedicines(ctx *gin.Context) {
	medicines, err := c.MedicineService.GetAll()

	if err != nil {
		appError := domainError.NewAppErrorWithType(domainError.UnknownError)
		_ = ctx.Error(appError)
		return
	}

	ctx.JSON(http.StatusOK, arrayDomainToResponseMapper(medicines))

}

// GetDataMedicines godoc
// @Tags medicine
// @Summary Get all Medicines
// @Description Get all Medicines on the system
// @Param   limit  query   string  true        "limit"
// @Param   page  query   string  true        "page"
// @Success 200 {object} []useCaseMedicine.PaginationResultMedicine
// @Failure 400 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /medicine [get]
func (c *Controller) GetDataMedicines(ctx *gin.Context) {
	var request DataMedicineRequest

	if err := controllers.BindJSON(ctx, &request); err != nil {
		appError := domainError.NewAppError(err, domainError.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	var dateRangeFiltersDomain []domain.DateRangeFilter = make([]domain.DateRangeFilter, len(request.FieldsDateRange))
	for i, dateRangeFilter := range request.FieldsDateRange {
		dateRangeFiltersDomain[i] = domain.DateRangeFilter{Field: dateRangeFilter.Field, Start: dateRangeFilter.StartDate, End: dateRangeFilter.EndDate}
	}

	users, err := c.MedicineService.GetData(request.Page, request.Limit, request.SorBy.Field, request.SorBy.Direction, request.Filters, request.GlobalSearch, dateRangeFiltersDomain)
	if err != nil {
		appError := domainError.NewAppErrorWithType(domainError.UnknownError)
		_ = ctx.Error(appError)
		return
	}

	numPages, nextCursor, prevCursor := controllers.PaginationValues(request.Limit, request.Page, users.Total)

	var response = PaginationResultMedicine{
		Data:       arrayDomainToResponseMapper(users.Data),
		Total:      users.Total,
		Limit:      request.Limit,
		Current:    request.Page,
		NextCursor: nextCursor,
		PrevCursor: prevCursor,
		NumPages:   numPages,
	}

	ctx.JSON(http.StatusOK, response)
}

// GetMedicinesByID godoc
// @Tags medicine
// @Summary Get medicines by ID
// @Description Get Medicines by ID on the system
// @Param medicine_id path int true "id of medicine"
// @Success 200 {object} domainMedicine.Medicine
// @Failure 400 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /medicine/{medicine_id} [get]
func (c *Controller) GetMedicinesByID(ctx *gin.Context) {
	medicineID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		appError := domainError.NewAppError(errors.New("medicine id is invalid"), domainError.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	domainMedicine, err := c.MedicineService.GetByID(medicineID)
	if err != nil {
		appError := domainError.NewAppError(err, domainError.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	ctx.JSON(http.StatusOK, domainMedicine)
}

// UpdateMedicine is the controller to update a medicine
func (c *Controller) UpdateMedicine(ctx *gin.Context) {
	medicineID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		appError := domainError.NewAppError(errors.New("param id is necessary in the url"), domainError.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	var requestMap map[string]any

	err = controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		appError := domainError.NewAppError(err, domainError.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	err = updateValidation(requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	var medicine *domainMedicine.Medicine
	medicine, err = c.MedicineService.Update(medicineID, requestMap)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, medicine)

}

// DeleteMedicine is the controller to delete a medicine
func (c *Controller) DeleteMedicine(ctx *gin.Context) {
	medicineID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		appError := domainError.NewAppError(errors.New("param id is necessary in the url"), domainError.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	err = c.MedicineService.Delete(medicineID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}
