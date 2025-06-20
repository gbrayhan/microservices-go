package user

import (
	"encoding/json"
	"time"

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

// UserRepositoryInterface defines the interface for user repository operations
type UserRepositoryInterface interface {
	GetAll() (*[]domainUser.User, error)
	Create(userDomain *domainUser.User) (*domainUser.User, error)
	GetByID(id int) (*domainUser.User, error)
	GetByEmail(email string) (*domainUser.User, error)
	Update(id int, userMap map[string]interface{}) (*domainUser.User, error)
	Delete(id int) error
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
	err := r.DB.Model(&userObj).
		Select("user_name", "email", "first_name", "last_name", "status", "role").
		Updates(userMap).Error
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
