package medicine

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"

	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainMedicine "github.com/gbrayhan/microservices-go/src/domain/medicine"
	"gorm.io/gorm"
)

// MedicineRepositoryInterface defines the interface for medicine repository operations
type MedicineRepositoryInterface interface {
	GetAll() (*[]domainMedicine.Medicine, error)
	GetData(page int64, limit int64, sortBy string, sortDirection string, filters map[string][]string, searchText string, dateRangeFilters []domain.DateRangeFilter) (*domainMedicine.DataMedicine, error)
	GetByID(id int) (*domainMedicine.Medicine, error)
	Create(medicine *domainMedicine.Medicine) (*domainMedicine.Medicine, error)
	GetByMap(medicineMap map[string]any) (*domainMedicine.Medicine, error)
	Delete(id int) error
	Update(id int, medicineMap map[string]any) (*domainMedicine.Medicine, error)
}

// Structures
type Medicine struct {
	ID          int    `gorm:"primaryKey"`
	Name        string `gorm:"unique"`
	Description string
	EANCode     string `gorm:"unique"`
	Laboratory  string
	CreatedAt   time.Time `gorm:"autoCreateTime:milli"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:milli"`
}

type PaginationResultMedicine struct {
	Data       *[]domainMedicine.Medicine
	Total      int64
	Limit      int64
	Current    int64
	NextCursor uint
	PrevCursor uint
	NumPages   int64
}

func (*Medicine) TableName() string {
	return "medicines"
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

type Repository struct {
	DB *gorm.DB
}

func NewMedicineRepository(DB *gorm.DB) MedicineRepositoryInterface {
	return &Repository{
		DB: DB}
}

func (r *Repository) GetData(page int64, limit int64, sortBy string, sortDirection string, filters map[string][]string, searchText string, dateRangeFilters []domain.DateRangeFilter) (*domainMedicine.DataMedicine, error) {
	var medicines []Medicine
	var total int64
	offset := (page - 1) * limit

	var searchColumns = []string{"name", "description", "ean_code", "laboratory"}

	countResult := make(chan error)
	go func() {
		err := r.DB.Model(&Medicine{}).
			Scopes(ApplyFilters(ColumnsMedicineMapping, filters, dateRangeFilters, searchText, searchColumns)).
			Count(&total).Error
		countResult <- err
	}()

	queryResult := make(chan error)
	go func() {
		query, err := ComplementSearch((*gorm.DB)(nil), sortBy, sortDirection, limit, offset, filters, dateRangeFilters, searchText, searchColumns, ColumnsMedicineMapping)
		if err != nil {
			queryResult <- err
			return
		}
		if query == nil {
			query = r.DB
		} else {
			query = r.DB.Scopes(ApplyFilters(ColumnsMedicineMapping, filters, dateRangeFilters, searchText, searchColumns))
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
		if !IsZeroValue(value) {
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

// Mappers
func (m *Medicine) toDomainMapper() *domainMedicine.Medicine {
	return &domainMedicine.Medicine{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		EanCode:     m.EANCode,
		Laboratory:  m.Laboratory,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func arrayToDomainMapper(medicines *[]Medicine) *[]domainMedicine.Medicine {
	medicinesDomain := make([]domainMedicine.Medicine, len(*medicines))
	for i, medicine := range *medicines {
		medicinesDomain[i] = *medicine.toDomainMapper()
	}
	return &medicinesDomain
}

// IsZeroValue checks if a value is the zero value of its type
func IsZeroValue(value any) bool {
	return reflect.DeepEqual(value, reflect.Zero(reflect.TypeOf(value)).Interface())
}

// ComplementSearch function moved from Utils.go
func ComplementSearch(r *gorm.DB, sortBy string, sortDirection string, limit int64, offset int64, filters map[string][]string, dateRangeFilters []domain.DateRangeFilter, searchText string, searchColumns []string, columnMapping map[string]string) (query *gorm.DB, err error) {
	if r == nil {
		return nil, nil
	}

	query = r
	if sortBy != "" {
		orderClause := fmt.Sprintf("%s %s", columnMapping[sortBy], sortDirection)
		query = query.Order(orderClause).Limit(int(limit)).Offset(int(offset))
	} else {
		query = query.Limit(int(limit)).Offset(int(offset))
	}

	if len(filters) > 0 {
		filters = UpdateFilterKeys(filters, columnMapping)
		for key, values := range filters {
			query = query.Where(fmt.Sprintf("%s IN (?)", key), values)
		}
	}

	if len(dateRangeFilters) > 0 {
		for i := range dateRangeFilters {
			if newFieldName, ok := columnMapping[dateRangeFilters[i].Field]; ok {
				dateRangeFilters[i].Field = newFieldName
			}
		}
		for _, filter := range dateRangeFilters {
			query = query.Where(fmt.Sprintf("%s BETWEEN ? AND ?", filter.Field), filter.Start, filter.End)
		}
	}

	if searchText != "" {
		var orConditions []string
		for _, column := range searchColumns {
			orConditions = append(orConditions, fmt.Sprintf("%s LIKE '%%%s%%'", column, searchText))
		}
		searchQuery := fmt.Sprintf("AND (%s)", strings.Join(orConditions, " OR "))
		query = query.Where(fmt.Sprintf("1=1 %s", searchQuery))
	}
	return
}

// UpdateFilterKeys function moved from Utils.go
func UpdateFilterKeys(filters map[string][]string, columnMapping map[string]string) map[string][]string {
	updatedFilters := make(map[string][]string)
	for key, value := range filters {
		if updatedKey, ok := columnMapping[key]; ok {
			updatedFilters[updatedKey] = value
		} else {
			updatedFilters[key] = value
		}
	}
	return updatedFilters
}

// ApplyFilters function moved from Utils.go
func ApplyFilters(columnMapping map[string]string, filters map[string][]string, dateRangeFilters []domain.DateRangeFilter, searchText string, searchColumns []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := db
		if len(filters) > 0 {
			filters = UpdateFilterKeys(filters, columnMapping)
			for key, values := range filters {
				query = query.Where(fmt.Sprintf("%s IN (?)", key), values)
			}
		}
		if len(dateRangeFilters) > 0 {
			for _, filter := range dateRangeFilters {
				if newFieldName, ok := columnMapping[filter.Field]; ok {
					filter.Field = newFieldName
				}
				query = query.Where(fmt.Sprintf("%s BETWEEN ? AND ?", filter.Field), filter.Start, filter.End)
			}
		}
		if searchText != "" && len(searchColumns) > 0 {
			var orConditions []string
			var args []interface{}
			for _, column := range searchColumns {
				orConditions = append(orConditions, fmt.Sprintf("%s LIKE ?", column))
				args = append(args, "%"+searchText+"%")
			}
			searchQuery := fmt.Sprintf("(%s)", strings.Join(orConditions, " OR "))
			query = query.Where(searchQuery, args...)
		}
		return query
	}
}
