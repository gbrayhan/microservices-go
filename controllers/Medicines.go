package controllers

import (
  "net/http"
  "strconv"

  "github.com/gin-gonic/gin"

  "github.com/gbrayhan/microservices-go/models"
  "github.com/gbrayhan/microservices-go/validator"
)

type MedicineController struct {
  Name        string `json:"name" example:"Paracetamol"`
  Description string `json:"description" example:"Something"`
  Laboratory  string `json:"laboratory" example:"Roche"`
  EanCode     string `json:"ean_code" example:"122000000021"`
}

//  NewMedicine godoc
// @Tags medicine
// @Summary Create New Medicine
// @Description Create new medicine on the system
// @Accept  json
// @Produce  json
// @Param data body MedicineController true "body data"
// @Success 200 {object} models.Medicine
// @Failure 400 {object} GeneralResponse
// @Failure 500 {object} GeneralResponse
// @Router /medicine/new [post]
func NewMedicine(c *gin.Context) {
  var request MedicineController

  _ = bindJSON(c, &request)

  if messagesError := validator.General(request, nil); messagesError != nil {
    badRequest(c, messagesError)
    return
  }

  medicineModel := models.Medicine{Name: request.Name, Description: request.Description, Laboratory: request.Laboratory, EANCode: request.EanCode}
  if err := medicineModel.SaveMedicine(); err != nil {
    serverError(c, err)
    return
  }
  c.JSON(http.StatusOK, medicineModel)
}

//  GetAllMedicine godoc
// @Tags medicine
// @Summary Get all Medicines
// @Description Get all Medicines on the system
// @Success 200 {object} []models.Medicine
// @Failure 400 {object} GeneralResponse
// @Failure 500 {object} GeneralResponse
// @Router /medicine/get-all [get]
func GetAllMedicines(c *gin.Context) {
  medicines, err := models.GetAllMedicines()
  if err != nil {
    serverError(c, err)
    return
  }

  c.JSON(http.StatusOK, medicines)
}

//  GetMedicinesByID godoc
// @Tags medicine
// @Summary Get medicines by ID
// @Description Get Medicines by ID on the system
// @Param medicine_id path int true "Id of medicine"
// @Success 200 {object} models.Medicine
// @Failure 400 {object} GeneralResponse
// @Failure 500 {object} GeneralResponse
// @Router /medicine/get-by-id/{medicine_id} [get]
func GetMedicineByID(c *gin.Context) {
  var medicine models.Medicine
  medicineID, err := strconv.Atoi(c.Param("medicine-id"))
  medicine.ID = medicineID
  if err != nil {
    badRequest(c, []string{"Medicine ID is invalid"})
    return
  }
  err = medicine.GetMedicineByID()
  if err != nil {
    serverError(c, err)
    return
  }

  if medicine.Name == "" {
    badRequest(c, []string{"Medicine ID is not found"})
    return
  }

  c.JSON(http.StatusOK, medicine)
}
