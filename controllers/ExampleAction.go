package controllers

import (
	"github.com/gbrayhan/microservices-go/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GeneralRequest struct {
	ID       int    `json:"id"`
	UserName string `json:"user_name"`
}

// Validate input data
func validateExample(request *GeneralRequest) (messages []string) {
	// Rules to validate
	if request.ID == 0 {
		messages = append(messages, "Field (id) is required.")
	}

	return
}

/*
 * Actions Controllers
 */
func ExampleAction(c *gin.Context) {
	var (
		request GeneralRequest
		element models.ExampleElement
	)
	_ = bindJSON(c, &request)


	if messagesError := validateExample(&request); messagesError != nil {
		BadRequest(c, messagesError)
		return
	}

	element.ID = request.ID
	if err := element.CompleteDataID(); err != nil {
		ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, element)
}
