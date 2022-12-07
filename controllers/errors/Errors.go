package errors

import (
  "net/http"

  appError "github.com/gbrayhan/microservices-go/models/errors"
  "github.com/gin-gonic/gin"
)

type MessagesResponse struct {
  Message string `json:"message"`
}

// Handler is Gin middleware to handle errors.
func Handler(c *gin.Context) {
  // Execute request handlers and then handle any errors
  c.Next()

  errs := c.Errors

  if len(errs) > 0 {
    err, ok := errs[0].Err.(*appError.AppError)
    if ok {
      resp := MessagesResponse{Message: err.Error()}
      switch err.Type {
      case appError.NotFound:
        c.JSON(http.StatusNotFound, resp)
        return
      case appError.ValidationError:
        c.JSON(http.StatusBadRequest, resp)
        return
      case appError.ResourceAlreadyExists:
        c.JSON(http.StatusConflict, resp)
        return
      case appError.NotAuthenticated:
        c.JSON(http.StatusUnauthorized, resp)
        return
      case appError.NotAuthorized:
        c.JSON(http.StatusForbidden, resp)
        return
      case appError.RepositoryError:
        c.JSON(http.StatusInternalServerError, MessagesResponse{Message: "We are working to improve the flow of this request."})
        return
      default:
        c.JSON(http.StatusInternalServerError, MessagesResponse{Message: "We are working to improve the flow of this request."})
        return
      }
    }

    return
  }
}
