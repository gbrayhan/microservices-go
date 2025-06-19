package medicine

import (
	"encoding/json"
	"fmt"

	"github.com/gbrayhan/microservices-go/src/domain"

	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainMedicine "github.com/gbrayhan/microservices-go/src/domain/medicine"
	"github.com/gbrayhan/microservices-go/src/infrastructure"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/utils"
	"gorm.io/gorm"
)

func (*Medicine) TableName() string {
	return "medicines"
}

type Repository struct {
	DB *gorm.DB
}

func NewMedicineRepository(DB *gorm.DB) infrastructure.MedicineRepositoryInterface {
	return &Repository{
		DB: DB}
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

func (r *Repository) GetData(page int64, limit int64, sortBy string, sortDirection string, filters map[string][]string, searchText string, dateRangeFilters []domain.DateRangeFilter) (*domainMedicine.DataMedicine, error) {
	var medicines []Medicine
	var total int64
	offset := (page - 1) * limit

	var searchColumns = []string{"name", "description", "ean_code", "laboratory"}

	countResult := make(chan error)
	go func() {
		err := r.DB.Model(&Medicine{}).
			Scopes(utils.ApplyFilters(ColumnsMedicineMapping, filters, dateRangeFilters, searchText, searchColumns)).
			Count(&total).Error
		countResult <- err
	}()

	queryResult := make(chan error)
	go func() {
		query, err := utils.ComplementSearch((*gorm.DB)(nil), sortBy, sortDirection, limit, offset, filters, dateRangeFilters, searchText, searchColumns, ColumnsMedicineMapping)
		if err != nil {
			queryResult <- err
			return
		}
		if query == nil {
			query = r.DB
		} else {
			query = r.DB.Scopes(utils.ApplyFilters(ColumnsMedicineMapping, filters, dateRangeFilters, searchText, searchColumns))
		}

		err = query.Find(&medicines).Error
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
		Data:  arrayToDomainMapper(&medicines),
		Total: total,
	}, nil
}

func (r *Repository) Create(newMedicine *domainMedicine.Medicine) (*domainMedicine.Medicine, error) {
	medicine := &Medicine{
		Name:        newMedicine.Name,
		Description: newMedicine.Description,
		EANCode:     newMedicine.EanCode,
		Laboratory:  newMedicine.Laboratory,
	}

	tx := r.DB.Create(medicine)
	if tx.Error != nil {
		byteErr, _ := json.Marshal(tx.Error)
		var newError domainErrors.GormErr
		err := json.Unmarshal(byteErr, &newError)
		if err != nil {
			return nil, err
		}
		switch newError.Number {
		case 1062:
			return nil, domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
		default:
			return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	return medicine.toDomainMapper(), nil
}

func (r *Repository) GetByID(id int) (*domainMedicine.Medicine, error) {
	var medicine Medicine
	err := r.DB.Where("id = ?", id).First(&medicine).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return medicine.toDomainMapper(), nil
}

func (r *Repository) GetByMap(medicineMap map[string]any) (*domainMedicine.Medicine, error) {
	var medicine Medicine
	tx := r.DB.Limit(1)
	for key, value := range medicineMap {
		if !utils.IsZeroValue(value) {
			tx = tx.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}
	err := tx.Find(&medicine).Error
	if err != nil {
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return medicine.toDomainMapper(), nil
}

func (r *Repository) Update(id int, medicineMap map[string]any) (*domainMedicine.Medicine, error) {
	var med Medicine
	med.ID = id
	err := r.DB.Model(&med).
		Select("name", "description", "ean_code", "laboratory").
		Updates(medicineMap).Error
	if err != nil {
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return nil, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			return nil, domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
		default:
			return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	err = r.DB.Where("id = ?", id).First(&med).Error
	return med.toDomainMapper(), err
}

func (r *Repository) Delete(id int) error {
	tx := r.DB.Delete(&Medicine{}, id)
	if tx.Error != nil {
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	return nil
}

func (r *Repository) GetAll() (*[]domainMedicine.Medicine, error) {
	var medicines []Medicine
	err := r.DB.Find(&medicines).Error
	if err != nil {
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return arrayToDomainMapper(&medicines), nil
}
