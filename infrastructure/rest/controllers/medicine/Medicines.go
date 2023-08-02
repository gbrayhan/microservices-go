// Package medicine contains the medicine controller
package medicine

import (
	"errors"

	useCaseMedicine "github.com/gbrayhan/microservices-go/application/usecases/medicine"
	domainError "github.com/gbrayhan/microservices-go/domain/errors"
	domainMedicine "github.com/gbrayhan/microservices-go/domain/medicine"
	"github.com/gbrayhan/microservices-go/infrastructure/rest/controllers"

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
	newMedicine := useCaseMedicine.NewMedicine{
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

	ctx.JSON(http.StatusOK, domainMedicine)
}

// GetAllMedicines godoc
// @Tags medicine
// @Summary Get all Medicines
// @Description Get all Medicines on the system
// @Param   limit  query   string  true        "limit"
// @Param   page  query   string  true        "page"
// @Success 200 {object} []useCaseMedicine.PaginationResultMedicine
// @Failure 400 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /medicine [get]
func (c *Controller) GetAllMedicines(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "20")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		appError := domainError.NewAppError(errors.New("param page is necessary to be an integer"), domainError.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		appError := domainError.NewAppError(errors.New("param limit is necessary to be an integer"), domainError.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	medicines, err := c.MedicineService.GetAll(page, limit)
	if err != nil {
		appError := domainError.NewAppErrorWithType(domainError.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	ctx.JSON(http.StatusOK, medicines)
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
