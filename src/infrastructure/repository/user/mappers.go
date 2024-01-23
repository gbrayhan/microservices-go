// Package user contains the business logic for the user entity
package user

import (
	domainUser "github.com/gbrayhan/microservices-go/src/domain/user"
)

func (user *User) toDomainMapper() *domainUser.User {
	return &domainUser.User{
		ID:           user.ID,
		UserName:     user.UserName,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Status:       user.Status,
		HashPassword: user.HashPassword,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

func fromDomainMapper(user *domainUser.User) *User {
	return &User{
		ID:           user.ID,
		UserName:     user.UserName,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Status:       user.Status,
		HashPassword: user.HashPassword,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

func arrayToDomainMapper(users *[]User) *[]domainUser.User {
	usersDomain := make([]domainUser.User, len(*users))
	for i, user := range *users {
		usersDomain[i] = *user.toDomainMapper()
	}

	return &usersDomain
}
