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
				status, message := domainErrors.AppErrorToHTTP(appErr)
				c.JSON(status, gin.H{"error": message})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			}
		}
	}
}
