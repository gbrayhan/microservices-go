package errors

import (
  "net/http"

  models "github.com/gbrayhan/microservices-go/models/errors"
  "github.com/gin-gonic/gin"
)

type messagesResponse struct {
  Messages []string `json:"messages"`
}

// Handler is Gin middleware to handle errors.
func Handler(c *gin.Context) {
  // Execute request handlers and then handle any errors
  c.Next()

  errs := c.Errors

  if len(errs) > 0 {
    err, ok := errs[0].Err.(*models.AppError)
    if ok {
      resp := messagesResponse{Messages: []string{err.Error()}}
      switch err.Type {
      case models.NotFound:
        c.JSON(http.StatusNotFound, resp)
        return
      case models.ValidationError:
        c.JSON(http.StatusBadRequest, resp)
        return
      case models.ResourceAlreadyExists:
        c.JSON(http.StatusConflict, resp)
        return
      case models.NotAuthenticated:
        c.JSON(http.StatusUnauthorized, resp)
        return
      case models.NotAuthorized:
        c.JSON(http.StatusForbidden, resp)
        return
      case models.RepositoryError:
        c.JSON(http.StatusInternalServerError, messagesResponse{Messages: []string{"We are working to improve the flow of this request."}})
        return
      default:
        c.JSON(http.StatusInternalServerError, messagesResponse{Messages: []string{"We are working to improve the flow of this request."}})
        return
      }
    }

    // Error is not AppError, return a generic internal server error
    c.JSON(http.StatusInternalServerError, "Internal Server Error")
    return
  }
}
