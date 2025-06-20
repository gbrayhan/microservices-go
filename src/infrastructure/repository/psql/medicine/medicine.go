package medicine

import (
	"encoding/json"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainMedicine "github.com/gbrayhan/microservices-go/src/domain/medicine"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MedicineRepositoryInterface defines the interface for medicine repository operations
type MedicineRepositoryInterface interface {
	GetAll() (*[]domainMedicine.Medicine, error)
	GetByID(id int) (*domainMedicine.Medicine, error)
	Create(medicine *domainMedicine.Medicine) (*domainMedicine.Medicine, error)
	Delete(id int) error
	Update(id int, medicineMap map[string]any) (*domainMedicine.Medicine, error)
	SearchPaginated(filters domain.DataFilters) (*domainMedicine.SearchResultMedicine, error)
	SearchByProperty(property string, searchText string) (*[]string, error)
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
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewMedicineRepository(DB *gorm.DB, loggerInstance *logger.Logger) MedicineRepositoryInterface {
	return &Repository{
		DB:     DB,
		Logger: loggerInstance,
	}
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
		r.Logger.Error("Error creating medicine", zap.Error(tx.Error), zap.String("name", newMedicine.Name))
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
	r.Logger.Info("Successfully created medicine", zap.String("name", newMedicine.Name), zap.Int("id", medicine.ID))
	return medicine.toDomainMapper(), nil
}

func (r *Repository) GetByID(id int) (*domainMedicine.Medicine, error) {
	var medicine Medicine
	err := r.DB.Where("id = ?", id).First(&medicine).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Logger.Warn("Medicine not found", zap.Int("id", id))
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		r.Logger.Error("Error getting medicine by ID", zap.Error(err), zap.Int("id", id))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	r.Logger.Info("Successfully retrieved medicine by ID", zap.Int("id", id))
	return medicine.toDomainMapper(), nil
}

func (r *Repository) Update(id int, medicineMap map[string]any) (*domainMedicine.Medicine, error) {
	var med Medicine
	med.ID = id
	err := r.DB.Model(&med).
		Select("name", "description", "ean_code", "laboratory").
		Updates(medicineMap).Error
	if err != nil {
		r.Logger.Error("Error updating medicine", zap.Error(err), zap.Int("id", id))
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
	if err != nil {
		r.Logger.Error("Error retrieving updated medicine", zap.Error(err), zap.Int("id", id))
		return nil, err
	}
	r.Logger.Info("Successfully updated medicine", zap.Int("id", id))
	return med.toDomainMapper(), nil
}

func (r *Repository) Delete(id int) error {
	tx := r.DB.Delete(&Medicine{}, id)
	if tx.Error != nil {
		r.Logger.Error("Error deleting medicine", zap.Error(tx.Error), zap.Int("id", id))
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		r.Logger.Warn("Medicine not found for deletion", zap.Int("id", id))
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	r.Logger.Info("Successfully deleted medicine", zap.Int("id", id))
	return nil
}

func (r *Repository) GetAll() (*[]domainMedicine.Medicine, error) {
	var medicines []Medicine
	if err := r.DB.Find(&medicines).Error; err != nil {
		r.Logger.Error("Error getting all medicines", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	r.Logger.Info("Successfully retrieved all medicines", zap.Int("count", len(medicines)))
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

func (r *Repository) SearchPaginated(filters domain.DataFilters) (*domainMedicine.SearchResultMedicine, error) {
	query := r.DB.Model(&Medicine{})

	// Apply like filters
	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" {
					column := ColumnsMedicineMapping[field]
					if column != "" {
						query = query.Where(column+" ILIKE ?", "%"+value+"%")
					}
				}
			}
		}
	}

	// Apply exact matches
	for field, values := range filters.Matches {
		if len(values) > 0 {
			column := ColumnsMedicineMapping[field]
			if column != "" {
				query = query.Where(column+" IN ?", values)
			}
		}
	}

	// Apply date range filters
	for _, dateFilter := range filters.DateRangeFilters {
		column := ColumnsMedicineMapping[dateFilter.Field]
		if column != "" {
			if dateFilter.Start != nil {
				query = query.Where(column+" >= ?", dateFilter.Start)
			}
			if dateFilter.End != nil {
				query = query.Where(column+" <= ?", dateFilter.End)
			}
		}
	}

	// Apply sorting
	if len(filters.SortBy) > 0 && filters.SortDirection.IsValid() {
		for _, sortField := range filters.SortBy {
			column := ColumnsMedicineMapping[sortField]
			if column != "" {
				query = query.Order(column + " " + string(filters.SortDirection))
			}
		}
	}

	// Count total records
	var total int64
	clonedQuery := query
	clonedQuery.Count(&total)

	// Apply pagination
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 10
	}
	offset := (filters.Page - 1) * filters.PageSize

	var medicines []Medicine
	if err := query.Offset(offset).Limit(filters.PageSize).Find(&medicines).Error; err != nil {
		r.Logger.Error("Error searching medicines", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	totalPages := int((total + int64(filters.PageSize) - 1) / int64(filters.PageSize))

	result := &domainMedicine.SearchResultMedicine{
		Data:       arrayToDomainMapper(&medicines),
		Total:      total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}

	r.Logger.Info("Successfully searched medicines",
		zap.Int64("total", total),
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))

	return result, nil
}

func (r *Repository) SearchByProperty(property string, searchText string) (*[]string, error) {
	column := ColumnsMedicineMapping[property]
	if column == "" {
		r.Logger.Warn("Invalid property for search", zap.String("property", property))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.ValidationError)
	}

	var coincidences []string
	if err := r.DB.Model(&Medicine{}).
		Distinct(column).
		Where(column+" ILIKE ?", "%"+searchText+"%").
		Limit(20).
		Pluck(column, &coincidences).Error; err != nil {
		r.Logger.Error("Error searching by property", zap.Error(err), zap.String("property", property))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	r.Logger.Info("Successfully searched by property",
		zap.String("property", property),
		zap.Int("results", len(coincidences)))

	return &coincidences, nil
}
