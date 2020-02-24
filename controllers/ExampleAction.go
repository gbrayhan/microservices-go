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

/*
 * Section to validate input data
 */

func (request *GeneralRequest) validateExample() (err error, messages []string) {
	// Rules to validate

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

	if err := c.ShouldBindJSON(&request); err != nil {
		BadRequest(c, []string{err.Error()})
		return
	}

	if err, messages := request.validateExample(); err != nil {
		BadRequest(c, messages)
		return
	}

	element.ID = request.ID
	if err := element.CompleteDataID(); err != nil {
		// TODO: Action log to internal server
		ServerError(c)
		return
	}
	c.JSON(http.StatusOK, element)
}
