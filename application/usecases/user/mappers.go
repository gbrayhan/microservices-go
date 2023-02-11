package user

import (
	domainUser "github.com/gbrayhan/microservices-go/domain/user"
)

func (n *NewUser) toDomainMapper() *domainUser.User {
	return &domainUser.User{
		UserName:  n.UserName,
		Email:     n.Email,
		FirstName: n.FirstName,
		LastName:  n.LastName,
		Role:      n.Role,
	}
}
