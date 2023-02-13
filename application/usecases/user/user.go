// Package user provides the use case for user
package user

import (
	userDomain "github.com/gbrayhan/microservices-go/domain/user"
	userRepository "github.com/gbrayhan/microservices-go/infrastructure/repository/user"
	"golang.org/x/crypto/bcrypt"
)

// Service is a struct that contains the repository implementation for user use case
type Service struct {
	UserRepository userRepository.Repository
}

// GetAll is a function that returns all users
func (s *Service) GetAll() (*[]userDomain.User, error) {
	return s.UserRepository.GetAll()
}

// GetByID is a function that returns a user by id
func (s *Service) GetByID(id int) (*userDomain.User, error) {
	return s.UserRepository.GetByID(id)
}

// Create is a function that creates a new user
func (s *Service) Create(newUser *NewUser) (*userDomain.User, error) {
	domain := newUser.toDomainMapper()

	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return &userDomain.User{}, err
	}
	domain.HashPassword = string(hash)
	domain.Status = true
	err = s.UserRepository.Create(domain)

	return domain, err
}

// GetOneByMap is a function that returns a user by map
func (s *Service) GetOneByMap(userMap map[string]interface{}) (*userDomain.User, error) {
	return s.UserRepository.GetOneByMap(userMap)
}

// Delete is a function that deletes a user by id
func (s *Service) Delete(id int) error {
	return s.UserRepository.Delete(id)
}

// Update is a function that updates a user by id
func (s *Service) Update(id int, userMap map[string]interface{}) (*userDomain.User, error) {
	return s.UserRepository.Update(id, userMap)
}
