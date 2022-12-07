package medicine

import (
  "errors"
  "github.com/gbrayhan/microservices-go/controllers"
  "github.com/gin-gonic/gin"
  "net/http"
  "strconv"

  _ "github.com/gbrayhan/microservices-go/controllers/errors"
  errorModels "github.com/gbrayhan/microservices-go/models/errors"
  model "github.com/gbrayhan/microservices-go/models/medicine"
)

// NewMedicine godoc
// @Tags medicine
// @Summary Create New Medicine
// @Description Create new medicine on the system
// @Accept  json
// @Produce  json
// @Param data body NewMedicineRequest true "body data"
// @Success 200 {object} model.Medicine
// @Failure 400 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /medicine [post]
func NewMedicine(c *gin.Context) {
  var request NewMedicineRequest

  if err := controllers.BindJSON(c, &request); err != nil {
    appError := errorModels.NewAppError(err, errorModels.ValidationError)
    _ = c.Error(appError)
    return
  }
  medicine := model.Medicine{
    Name:        request.Name,
    Description: request.Description,
    Laboratory:  request.Laboratory,
    EANCode:     request.EanCode,
  }

  err := model.CreateMedicine(&medicine)
  if err != nil {
    _ = c.Error(err)
    return
  }

  c.JSON(http.StatusOK, medicine)
}

// GetAllMedicines godoc
// @Tags medicine
// @Summary Get all Medicines
// @Description Get all Medicines on the system
// @Success 200 {object} []model.Medicine
// @Failure 400 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /medicine [get]
func GetAllMedicines(c *gin.Context) {
  medicines, err := model.GetAllMedicines()
  if err != nil {
    appError := errorModels.NewAppErrorWithType(errorModels.UnknownError)
    _ = c.Error(appError)
    return
  }
  c.JSON(http.StatusOK, medicines)
}

// GetMedicinesByID godoc
// @Tags medicine
// @Summary Get medicines by ID
// @Description Get Medicines by ID on the system
// @Param medicine_id path int true "id of medicine"
// @Success 200 {object} model.Medicine
// @Failure 400 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /medicine/{medicine_id} [get]
func GetMedicinesByID(c *gin.Context) {
  var medicine model.Medicine
  medicineID, err := strconv.Atoi(c.Param("id"))
  if err != nil {
    appError := errorModels.NewAppError(errors.New("medicine id is invalid"), errorModels.ValidationError)
    _ = c.Error(appError)
    return
  }

  err = model.GetMedicineByID(&medicine, medicineID)
  if err != nil {
    appError := errorModels.NewAppError(err, errorModels.ValidationError)
    _ = c.Error(appError)
    return
  }

  c.JSON(http.StatusOK, medicine)
}

func UpdateMedicine(c *gin.Context) {
  medicineID, err := strconv.Atoi(c.Param("id"))
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

  medicine, err := model.UpdateMedicine(medicineID, requestMap)
  if err != nil {
    _ = c.Error(err)
    return
  }

  c.JSON(http.StatusOK, medicine)

}

func DeleteMedicine(c *gin.Context) {
  medicineID, err := strconv.Atoi(c.Param("id"))
  if err != nil {
    appError := errorModels.NewAppError(errors.New("param id is necessary in the url"), errorModels.ValidationError)
    _ = c.Error(appError)
    return
  }

  err = model.DeleteMedicine(medicineID)
  if err != nil {
    _ = c.Error(err)
    return
  }
  c.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})

}
