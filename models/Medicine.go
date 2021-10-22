package models

import (
  "encoding/json"
  "fmt"
  "github.com/gbrayhan/microservices-go/config"
  modelErrors "github.com/gbrayhan/microservices-go/models/errors"
  "time"
)

type Medicine struct {
  ID          int       `json:"id" example:"123" gorm:"primaryKey"`
  Name        string    `json:"name" example:"Paracetamol" gorm:"unique"`
  Description string    `json:"description" example:"Some Description"`
  EANCode     string    `json:"ean_code" example:"9900000124" gorm:"unique"`
  Laboratory  string    `json:"laboratory" example:"Roche"`
  CreatedAt   time.Time `json:"created_at,omitempty" example:"2021-02-24 20:19:39"`
  UpdatedAt   time.Time `json:"updated_at,omitempty" example:"2021-02-24 20:19:39"`
}

func (b *Medicine) TableName() string {
  return "medicines"
}

// GetAllMedicines Fetch all medicine data
func GetAllMedicines(medicine *[]Medicine) (err error) {
  err = config.DB.Find(medicine).Error
  if err != nil {
    return err
  }
  return nil
}

// CreateMedicine ... Insert New data
func CreateMedicine(medicine *Medicine) (err error) {
  err = config.DB.Create(medicine).Error

  if err != nil {
    byteErr, _ := json.Marshal(err)
    var newError modelErrors.GormErr
    err = json.Unmarshal(byteErr, &newError)
    if err != nil {
      return err
    }
    switch newError.Number {
    case 1062:
      err = modelErrors.NewAppErrorWithType(modelErrors.ResourceAlreadyExists)
      return

    default:
      err = modelErrors.NewAppErrorWithType(modelErrors.UnknownError)
    }
  }

  return
}

// GetMedicineByID ... Fetch only one medicine by Id
func GetMedicineByID(medicine *Medicine, id int) (err error) {
  err = config.DB.Where("id = ?", id).First(medicine).Error
  if err != nil {
    return err
  }
  return nil
}

// UpdateMedicine ... Update medicine
func UpdateMedicine(medicine *Medicine, id int) (err error) {
  fmt.Println(medicine)
  config.DB.Save(medicine)
  return nil
}

// DeleteMedicine ... Delete medicine
func DeleteMedicine(medicine *Medicine, id string) (err error) {
  config.DB.Where("id = ?", id).Delete(medicine)
  return
}
