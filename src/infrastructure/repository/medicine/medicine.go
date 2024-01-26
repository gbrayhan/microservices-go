// Package medicine contains the repository implementation for the medicine entity
package medicine

import (
	"encoding/json"
	"fmt"
	"github.com/gbrayhan/microservices-go/src/domain"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository"

	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainMedicine "github.com/gbrayhan/microservices-go/src/domain/medicine"
	"gorm.io/gorm"
)

// Repository is a struct that contains the database implementation for medicine entity
type Repository struct {
	DB *gorm.DB
}

var ColumnsMedicineMapping = map[string]string{
	"id":          "id",
	"name":        "name",
	"description": "description",
	"eanCode":     "ean_code",
	"laboratory":  "laboratory",
	"createdAt":   "created_at",
	"updatedAt":   "updated_at",
}

var ColumnsMedicineStructure = map[string]string{
	"id":          "ID",
	"name":        "Name",
	"description": "Description",
	"eanCode":     "EanCode",
	"laboratory":  "Laboratory",
	"createdAt":   "CreatedAt",
	"updatedAt":   "UpdatedAt",
}

// GetAll Fetch all medicine data
func (r *Repository) GetAll(page int64, limit int64, sortBy string, sortDirection string, filters map[string][]string, searchText string, dateRangeFilters []domain.DateRangeFilter) (*domainMedicine.DataMedicine, error) {
	var users []Medicine
	var total int64
	offset := (page - 1) * limit

	var searchColumns = []string{"name", "description", "ean_code", "laboratory"}

	countResult := make(chan error)
	go func() {
		err := r.DB.Model(&Medicine{}).Scopes(repository.ApplyFilters(ColumnsMedicineMapping, filters, dateRangeFilters, searchText, searchColumns)).Count(&total).Error
		countResult <- err
	}()

	queryResult := make(chan error)
	go func() {
		query, err := repository.ComplementSearch((*repository.Repository)(r), sortBy, sortDirection, limit, offset, filters, dateRangeFilters, searchText, searchColumns, ColumnsMedicineMapping)
		if err != nil {
			queryResult <- err
			return
		}
		err = query.Find(&users).Error
		queryResult <- err
	}()

	var countErr, queryErr error
	for i := 0; i < 2; i++ {
		select {
		case err := <-countResult:
			countErr = err
		case err := <-queryResult:
			queryErr = err
		}
	}

	if countErr != nil {
		return &domainMedicine.DataMedicine{}, countErr
	}
	if queryErr != nil {
		return &domainMedicine.DataMedicine{}, queryErr
	}

	return &domainMedicine.DataMedicine{
		Data:  arrayToDomainMapper(&users),
		Total: total,
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

	tx := r.DB.Limit(1)
	for key, value := range medicineMap {
		if !repository.IsZeroValue(value) {
			tx = tx.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}

	err := tx.Find(&medicine).Error
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
