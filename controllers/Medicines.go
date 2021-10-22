package controllers

import (
  "net/http"
  "strconv"

  "github.com/gin-gonic/gin"

  "github.com/gbrayhan/microservices-go/models"
)

type medicineController struct {
  Name        string `json:"name" example:"Paracetamol"`
  Description string `json:"description" example:"Something"`
  Laboratory  string `json:"laboratory" example:"Roche"`
  EanCode     string `json:"ean_code" example:"122000000021"`
}

// NewMedicine godoc
// @Tags medicine
// @Summary Create New Medicine
// @Description Create new medicine on the system
// @Accept  json
// @Produce  json
// @Param data body medicineController true "body data"
// @Success 200 {object} models.Medicine
// @Failure 400 {object} generalResponse
// @Failure 500 {object} generalResponse
// @Router /medicine/new [post]
func NewMedicine(c *gin.Context) {
  request := struct {
    Name        string `json:"name" example:"Paracetamol" gorm:"unique"`
    Description string `json:"description" example:"Something"`
    Laboratory  string `json:"laboratory" example:"Roche"`
    EanCode     string `json:"ean_code" example:"122000000021" gorm:"unique"`
  }{}

  if err := bindJSON(c, &request); err != nil {
    badRequest(c, []string{err.Error()})
    return
  }
  medicine := models.Medicine{
    Name:        request.Name,
    Description: request.Description,
    Laboratory:  request.Laboratory,
    EANCode:     request.EanCode,
  }

  err := models.CreateMedicine(&medicine)
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
// @Success 200 {object} []models.Medicine
// @Failure 400 {object} generalResponse
// @Failure 500 {object} generalResponse
// @Router /medicine/get-all [get]
func GetAllMedicines(c *gin.Context) {
  var medicines []models.Medicine
  err := models.GetAllMedicines(&medicines)
  if err != nil {
    serverError(c, err)
    return
  }
  c.JSON(http.StatusOK, medicines)
}

// GetMedicinesByID godoc
// @Tags medicine
// @Summary Get medicines by ID
// @Description Get Medicines by ID on the system
// @Param medicine_id path int true "id of medicine"
// @Success 200 {object} models.Medicine
// @Failure 400 {object} generalResponse
// @Failure 500 {object} generalResponse
// @Router /medicine/get-by-id/{medicine_id} [get]
func GetMedicinesByID(c *gin.Context) {
  var medicine models.Medicine
  medicineID, err := strconv.Atoi(c.Param("medicine-id"))
  if err != nil {
    badRequest(c, []string{"Medicine ID is invalid"})
    return
  }
  err = models.GetMedicineByID(&medicine, medicineID)
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
