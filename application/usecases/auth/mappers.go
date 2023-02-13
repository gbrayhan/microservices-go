// Package auth provides the use case for authentication
package auth

import (
	"github.com/gbrayhan/microservices-go/application/security/jwt"
	userDomain "github.com/gbrayhan/microservices-go/domain/user"
)

func secAuthUserMapper(domainUser *userDomain.User, authInfo *jwt.Auth) *SecurityAuthenticatedUser {
	return &SecurityAuthenticatedUser{
		Data: DataUserAuthenticated{
			UserName:  domainUser.UserName,
			Email:     domainUser.Email,
			FirstName: domainUser.FirstName,
			LastName:  domainUser.LastName,
			ID:        domainUser.ID,
			Status:    domainUser.Status,
		},
		Security: DataSecurityAuthenticated{
			JWTAccessToken:        authInfo.AccessToken,
			JWTRefreshToken:       authInfo.RefreshToken,
			ExpirationAccessTime:  authInfo.ExpirationAccessTime,
			ExpirationRefreshTime: authInfo.ExpirationRefreshTime,
		},
	}

}
