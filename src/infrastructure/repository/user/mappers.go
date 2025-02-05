package user

import (
	domainUser "github.com/gbrayhan/microservices-go/src/domain/user"
)

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
