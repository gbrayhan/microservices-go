package user

import (
	userDomain "github.com/gbrayhan/microservices-go/src/domain/user"
)

func domainToResponseMapper(domainUser *userDomain.User) *ResponseUser {
	return &ResponseUser{
		ID:        domainUser.ID,
		UserName:  domainUser.UserName,
		Email:     domainUser.Email,
		FirstName: domainUser.FirstName,
		LastName:  domainUser.LastName,
		Status:    domainUser.Status,
		CreatedAt: domainUser.CreatedAt,
		UpdatedAt: domainUser.UpdatedAt,
	}
}

func arrayDomainToResponseMapper(users *[]userDomain.User) *[]ResponseUser {
	res := make([]ResponseUser, len(*users))
	for i, u := range *users {
		res[i] = *domainToResponseMapper(&u)
	}
	return &res
}

func toUsecaseMapper(req *NewUserRequest) *userDomain.User {
	return &userDomain.User{
		UserName:  req.UserName,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
	}
}
