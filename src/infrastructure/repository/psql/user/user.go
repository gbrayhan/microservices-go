package user

import (
	"encoding/json"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainUser "github.com/gbrayhan/microservices-go/src/domain/user"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type User struct {
	ID           int       `gorm:"primaryKey"`
	UserName     string    `gorm:"column:user_name;unique"`
	Email        string    `gorm:"unique"`
	FirstName    string    `gorm:"column:first_name"`
	LastName     string    `gorm:"column:last_name"`
	Status       bool      `gorm:"column:status"`
	HashPassword string    `gorm:"column:hash_password"`
	CreatedAt    time.Time `gorm:"autoCreateTime:mili"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime:mili"`
}

func (User) TableName() string {
	return "users"
}

var ColumnsUserMapping = map[string]string{
	"id":           "id",
	"userName":     "user_name",
	"email":        "email",
	"firstName":    "first_name",
	"lastName":     "last_name",
	"status":       "status",
	"hashPassword": "hash_password",
	"createdAt":    "created_at",
	"updatedAt":    "updated_at",
}

// UserRepositoryInterface defines the interface for user repository operations
type UserRepositoryInterface interface {
	GetAll() (*[]domainUser.User, error)
	Create(userDomain *domainUser.User) (*domainUser.User, error)
	GetByID(id int) (*domainUser.User, error)
	GetByEmail(email string) (*domainUser.User, error)
	Update(id int, userMap map[string]interface{}) (*domainUser.User, error)
	Delete(id int) error
	SearchPaginated(filters domain.DataFilters) (*domainUser.SearchResultUser, error)
	SearchByProperty(property string, searchText string) (*[]string, error)
}

type Repository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewUserRepository(db *gorm.DB, loggerInstance *logger.Logger) UserRepositoryInterface {
	return &Repository{DB: db, Logger: loggerInstance}
}

func (r *Repository) GetAll() (*[]domainUser.User, error) {
	var users []User
	if err := r.DB.Find(&users).Error; err != nil {
		r.Logger.Error("Error getting all users", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	r.Logger.Info("Successfully retrieved all users", zap.Int("count", len(users)))
	return arrayToDomainMapper(&users), nil
}

func (r *Repository) Create(userDomain *domainUser.User) (*domainUser.User, error) {
	r.Logger.Info("Creating new user", zap.String("email", userDomain.Email))
	userRepository := fromDomainMapper(userDomain)
	txDb := r.DB.Create(userRepository)
	err := txDb.Error
	if err != nil {
		r.Logger.Error("Error creating user", zap.Error(err), zap.String("email", userDomain.Email))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return &domainUser.User{}, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			err = domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
			return &domainUser.User{}, err
		default:
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	r.Logger.Info("Successfully created user", zap.String("email", userDomain.Email), zap.Int("id", userRepository.ID))
	return userRepository.toDomainMapper(), err
}

func (r *Repository) GetByID(id int) (*domainUser.User, error) {
	var user User
	err := r.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Logger.Warn("User not found", zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		} else {
			r.Logger.Error("Error getting user by ID", zap.Error(err), zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		return &domainUser.User{}, err
	}
	r.Logger.Info("Successfully retrieved user by ID", zap.Int("id", id))
	return user.toDomainMapper(), nil
}

func (r *Repository) GetByEmail(email string) (*domainUser.User, error) {
	var user User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Logger.Warn("User not found", zap.String("email", email))
			err = domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		} else {
			r.Logger.Error("Error getting user by email", zap.Error(err), zap.String("email", email))
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		return &domainUser.User{}, err
	}
	r.Logger.Info("Successfully retrieved user by email", zap.String("email", email))
	return user.toDomainMapper(), nil
}

func (r *Repository) Update(id int, userMap map[string]interface{}) (*domainUser.User, error) {
	var userObj User
	userObj.ID = id

	// Map JSON field names to DB column names
	updateData := make(map[string]interface{})
	for k, v := range userMap {
		if column, ok := ColumnsUserMapping[k]; ok {
			updateData[column] = v
		} else {
			updateData[k] = v
		}
	}

	err := r.DB.Model(&userObj).
		Select("user_name", "email", "first_name", "last_name", "status", "role").
		Updates(updateData).Error
	if err != nil {
		r.Logger.Error("Error updating user", zap.Error(err), zap.Int("id", id))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return &domainUser.User{}, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			return &domainUser.User{}, domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
		default:
			return &domainUser.User{}, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	if err := r.DB.Where("id = ?", id).First(&userObj).Error; err != nil {
		r.Logger.Error("Error retrieving updated user", zap.Error(err), zap.Int("id", id))
		return &domainUser.User{}, err
	}
	r.Logger.Info("Successfully updated user", zap.Int("id", id))
	return userObj.toDomainMapper(), nil
}

func (r *Repository) Delete(id int) error {
	tx := r.DB.Delete(&User{}, id)
	if tx.Error != nil {
		r.Logger.Error("Error deleting user", zap.Error(tx.Error), zap.Int("id", id))
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		r.Logger.Warn("User not found for deletion", zap.Int("id", id))
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	r.Logger.Info("Successfully deleted user", zap.Int("id", id))
	return nil
}

func (r *Repository) SearchPaginated(filters domain.DataFilters) (*domainUser.SearchResultUser, error) {
	query := r.DB.Model(&User{})

	// Apply like filters
	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" {
					column := ColumnsUserMapping[field]
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
			column := ColumnsUserMapping[field]
			if column != "" {
				query = query.Where(column+" IN ?", values)
			}
		}
	}

	// Apply date range filters
	for _, dateFilter := range filters.DateRangeFilters {
		column := ColumnsUserMapping[dateFilter.Field]
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
			column := ColumnsUserMapping[sortField]
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

	var users []User
	if err := query.Offset(offset).Limit(filters.PageSize).Find(&users).Error; err != nil {
		r.Logger.Error("Error searching users", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	totalPages := int((total + int64(filters.PageSize) - 1) / int64(filters.PageSize))

	result := &domainUser.SearchResultUser{
		Data:       arrayToDomainMapper(&users),
		Total:      total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}

	r.Logger.Info("Successfully searched users",
		zap.Int64("total", total),
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))

	return result, nil
}

func (r *Repository) SearchByProperty(property string, searchText string) (*[]string, error) {
	column := ColumnsUserMapping[property]
	if column == "" {
		r.Logger.Warn("Invalid property for search", zap.String("property", property))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.ValidationError)
	}

	var coincidences []string
	if err := r.DB.Model(&User{}).
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

// Mappers
func (u *User) toDomainMapper() *domainUser.User {
	return &domainUser.User{
		ID:           u.ID,
		UserName:     u.UserName,
		Email:        u.Email,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Status:       u.Status,
		HashPassword: u.HashPassword,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func fromDomainMapper(u *domainUser.User) *User {
	return &User{
		ID:           u.ID,
		UserName:     u.UserName,
		Email:        u.Email,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Status:       u.Status,
		HashPassword: u.HashPassword,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func arrayToDomainMapper(users *[]User) *[]domainUser.User {
	usersDomain := make([]domainUser.User, len(*users))
	for i, user := range *users {
		usersDomain[i] = *user.toDomainMapper()
	}
	return &usersDomain
}
