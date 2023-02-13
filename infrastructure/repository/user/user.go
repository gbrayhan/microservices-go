// Package user contains the business logic for the user entity
package user

import (
	"encoding/json"
	domainErrors "github.com/gbrayhan/microservices-go/domain/errors"
	domainUser "github.com/gbrayhan/microservices-go/domain/user"
	"gorm.io/gorm"
)

// Repository is a struct that contains the database implementation for user entity
type Repository struct {
	DB *gorm.DB
}

// GetAll Fetch all user data
func (r *Repository) GetAll() (*[]domainUser.User, error) {
	var users []User
	err := r.DB.Find(&users).Error
	if err != nil {
		err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		return nil, err
	}

	return arrayToDomainMapper(&users), err
}

// Create ... Insert New data
func (r *Repository) Create(domainUser *domainUser.User) error {
	user := fromDomainMapper(domainUser)
	txDb := r.DB.Create(user)
	err := txDb.Error
	if err != nil {
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		err = json.Unmarshal(byteErr, &newError)
		if err != nil {
			return err
		}
		switch newError.Number {
		case 1062:
			err = domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
			return err

		default:
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}

	return err
}

// GetOneByMap ... Fetch only one user by Map values
func (r *Repository) GetOneByMap(userMap map[string]interface{}) (*domainUser.User, error) {
	var user User

	tx := r.DB.Where(userMap).Limit(1).Find(&user)
	if tx.Error != nil {
		err := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		return &domainUser.User{}, err
	}
	return user.toDomainMapper(), nil
}

// GetByID ... Fetch only one user by ID
func (r *Repository) GetByID(id int) (*domainUser.User, error) {
	var user User
	err := r.DB.Where("id = ?", id).First(&user).Error

	if err != nil {
		switch err.Error() {
		case gorm.ErrRecordNotFound.Error():
			err = domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		default:
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}

	return user.toDomainMapper(), err
}

// Update ... Update user
func (r *Repository) Update(id int, userMap map[string]interface{}) (*domainUser.User, error) {
	var user User

	user.ID = id
	err := r.DB.Model(&user).
		Select("user", "email", "firstName", "lastName").
		Updates(userMap).Error

	// err = config.DB.Save(user).Error
	if err != nil {
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		err = json.Unmarshal(byteErr, &newError)
		if err != nil {
			return &domainUser.User{}, err
		}
		switch newError.Number {
		case 1062:
			err = domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
		default:
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		return &domainUser.User{}, err

	}

	err = r.DB.Where("id = ?", id).First(&user).Error

	return user.toDomainMapper(), err
}

// Delete ... Delete user
func (r *Repository) Delete(id int) (err error) {
	tx := r.DB.Delete(&User{}, id)
	if tx.Error != nil {
		err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		return
	}

	if tx.RowsAffected == 0 {
		err = domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}

	return
}
