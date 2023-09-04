// Package medicine contains the repository implementation for the medicine entity
package medicine

import (
	"encoding/json"

	domainErrors "github.com/gbrayhan/microservices-go/domain/errors"
	domainMedicine "github.com/gbrayhan/microservices-go/domain/medicine"
	"gorm.io/gorm"
)

// Repository is a struct that contains the database implementation for medicine entity
type Repository struct {
	DB *gorm.DB
}

// GetAll Fetch all medicine data
func (r *Repository) GetAll(page int64, limit int64) (*PaginationResultMedicine, error) {
	var medicines []Medicine
	var total int64

	err := r.DB.Model(&Medicine{}).Count(&total).Error
	if err != nil {
		return &PaginationResultMedicine{}, err
	}
	offset := (page - 1) * limit
	err = r.DB.Limit(int(limit)).Offset(int(offset)).Find(&medicines).Error
	if err != nil {
		return &PaginationResultMedicine{}, err
	}

	numPages := (total + limit - 1) / limit
	var nextCursor, prevCursor uint
	if page < numPages {
		nextCursor = uint(page + 1)
	}
	if page > 1 {
		prevCursor = uint(page - 1)
	}

	return &PaginationResultMedicine{
		Data:       arrayToDomainMapper(&medicines),
		Total:      total,
		Limit:      limit,
		Current:    page,
		NextCursor: nextCursor,
		PrevCursor: prevCursor,
		NumPages:   numPages,
	}, nil
}

// Create ... Insert New data
func (r *Repository) Create(newMedicine *domainMedicine.Medicine) (createdMedicine *domainMedicine.Medicine, err error) {
	medicine := fromDomainMapper(newMedicine)

	tx := r.DB.Create(medicine)

	if tx.Error != nil {
		byteErr, _ := json.Marshal(tx.Error)
		var newError domainErrors.GormErr
		err = json.Unmarshal(byteErr, &newError)
		if err != nil {
			return
		}
		switch newError.Number {
		case 1062:
			err = domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
		default:
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		return
	}

	createdMedicine = medicine.toDomainMapper()
	return
}

// GetByID ... Fetch only one medicine by Id
func (r *Repository) GetByID(id int) (*domainMedicine.Medicine, error) {
	var medicine Medicine
	err := r.DB.Where("id = ?", id).First(&medicine).Error

	if err != nil {
		switch err.Error() {
		case gorm.ErrRecordNotFound.Error():
			err = domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		default:
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		return &domainMedicine.Medicine{}, err
	}

	return medicine.toDomainMapper(), nil
}

// GetOneByMap ... Fetch only one medicine by Map
func (r *Repository) GetOneByMap(medicineMap map[string]any) (*domainMedicine.Medicine, error) {
	var medicine Medicine

	err := r.DB.Where(medicineMap).Limit(1).Find(&medicine).Error
	if err != nil {
		err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		return nil, err
	}
	return medicine.toDomainMapper(), err
}

// Update ... Update medicine
func (r *Repository) Update(id int, medicineMap map[string]any) (*domainMedicine.Medicine, error) {
	var medicine Medicine

	medicine.ID = id
	err := r.DB.Model(&medicine).
		Select("name", "description", "ean_code", "laboratory").
		Updates(medicineMap).Error

	// err = config.DB.Save(medicine).Error
	if err != nil {
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		err = json.Unmarshal(byteErr, &newError)
		if err != nil {
			return &domainMedicine.Medicine{}, err
		}
		switch newError.Number {
		case 1062:
			err = domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
			return &domainMedicine.Medicine{}, err

		default:
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
			return &domainMedicine.Medicine{}, err
		}
	}

	err = r.DB.Where("id = ?", id).First(&medicine).Error

	return medicine.toDomainMapper(), err
}

// Delete ... Delete medicine
func (r *Repository) Delete(id int) (err error) {
	tx := r.DB.Delete(&domainMedicine.Medicine{}, id)
	if tx.Error != nil {
		err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		return
	}

	if tx.RowsAffected == 0 {
		err = domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}

	return
}
