package middlewares

import (
	"errors"
	"net/http"

	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var appErr *domainErrors.AppError
			if errors.As(err, &appErr) {
				switch appErr.Type {
				case domainErrors.NotFound:
					c.JSON(http.StatusNotFound, gin.H{"error": appErr.Error()})
				case domainErrors.ValidationError:
					c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Error()})
				case domainErrors.RepositoryError:
					c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.Error()})
				case domainErrors.NotAuthenticated:
					c.JSON(http.StatusUnauthorized, gin.H{"error": appErr.Error()})
				case domainErrors.NotAuthorized:
					c.JSON(http.StatusForbidden, gin.H{"error": appErr.Error()})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			}
		}
	}
}
