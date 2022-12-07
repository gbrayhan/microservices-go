package user

import (
  "encoding/json"
  "github.com/gbrayhan/microservices-go/config"
  appError "github.com/gbrayhan/microservices-go/models/errors"
  "github.com/jinzhu/gorm"
  "time"
)

type User struct {
  ID           int       `json:"id" example:"1099" gorm:"primaryKey"`
  User         string    `json:"user" example:"User" gorm:"unique"`
  Email        string    `json:"email" example:"some@mail.com" gorm:"unique"`
  FirstName    string    `json:"first_name" example:"John"`
  LastName     string    `json:"last_name" example:"Doe"`
  Status       bool      `json:"status" example:"false"`
  HashPassword string    `json:"hash_password" example:"SomeHashPass"`
  CreatedAt    time.Time `json:"created_at,omitempty" example:"2021-02-24 20:19:39" gorm:"autoCreateTime:mili"`
  UpdatedAt    time.Time `json:"updated_at,omitempty" example:"2021-02-24 20:19:39" gorm:"autoUpdateTime:mili"`
}

func (b *User) TableName() string {
  return "users"
}

// GetAllUsers Fetch all user data
func GetAllUsers() (user []User, err error) {
  err = config.DB.Find(&user).Error
  if err != nil {
    err = appError.NewAppErrorWithType(appError.UnknownError)
  }
  return
}

// CreateUser ... Insert New data
func CreateUser(user *User) (err error) {
  err = config.DB.Create(user).Error

  if err != nil {
    byteErr, _ := json.Marshal(err)
    var newError appError.GormErr
    err = json.Unmarshal(byteErr, &newError)
    if err != nil {
      return err
    }
    switch newError.Number {
    case 1062:
      err = appError.NewAppErrorWithType(appError.ResourceAlreadyExists)
      return

    default:
      err = appError.NewAppErrorWithType(appError.UnknownError)
    }
  }

  return
}

// GetUserByMap ... Fetch only one user by Map values
func GetUserByMap(userMap map[string]interface{}) (user User, err error) {
  tx := config.DB.Where(userMap).Limit(1).Find(&user)
  if tx.Error != nil {
    err = appError.NewAppErrorWithType(appError.UnknownError)
  }
  return
}

// GetUserByID ... Fetch only one user by Id
func GetUserByID(id int) (user User, err error) {
  err = config.DB.Where("id = ?", id).First(&user).Error

  if err != nil {
    switch err.Error() {
    case gorm.ErrRecordNotFound.Error():
      err = appError.NewAppErrorWithType(appError.NotFound)
    default:
      err = appError.NewAppErrorWithType(appError.UnknownError)
    }
  }

  return
}

// UpdateUser ... Update user
func UpdateUser(id int, userMap map[string]interface{}) (user User, err error) {
  user.ID = id
  err = config.DB.Model(&user).
    Select("user", "email", "firstName", "lastName").
    Updates(userMap).Error

  // err = config.DB.Save(user).Error
  if err != nil {
    byteErr, _ := json.Marshal(err)
    var newError appError.GormErr
    err = json.Unmarshal(byteErr, &newError)
    if err != nil {
      return
    }
    switch newError.Number {
    case 1062:
      err = appError.NewAppErrorWithType(appError.ResourceAlreadyExists)
      return

    default:
      err = appError.NewAppErrorWithType(appError.UnknownError)
    }
  }

  err = config.DB.Where("id = ?", id).First(&user).Error

  return
}

// DeleteUser ... Delete user
func DeleteUser(id int) (err error) {
  tx := config.DB.Delete(&User{}, id)
  if tx.Error != nil {
    err = appError.NewAppErrorWithType(appError.UnknownError)
    return
  }

  if tx.RowsAffected == 0 {
    err = appError.NewAppErrorWithType(appError.NotFound)
  }

  return
}
