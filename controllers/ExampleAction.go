package controllers

import (
	"github.com/banwire/microservice_golang/models"
	"github.com/gin-gonic/gin"
)

/*
 * Section to validate input data
 */

func validateExample(order *models.ExampleElement, response *models.ResponseBase) {
	// Rules to validate
	response.AppendSuccessMessage("successful validation")
}

/*
 * Actions Controllers
 */

func ExampleAction(c *gin.Context) {
	// Custom Request Client
	type requestClient struct {
		ID       int    `json:"id"`
		UserName string `json:"user_name"`
	}
	type responseClient struct {
		Messages  []string `json:"messages"`
		Merchant  string   `json:"merchant"`
		Reference string   `json:"reference"`
	}

	var (
		reqClient      requestClient
		responseGlobal models.ResponseBase
		element        models.ExampleElement
		respClient     responseClient
	)

	defer responseGlobal.ShowResponseJSON(c, &respClient)

	if err := c.ShouldBindJSON(&reqClient); err != nil {
		responseGlobal.AppendErrorRequest(err.Error())
		return
	}

	element.ID = reqClient.ID

	if element.CompleteDataID(&responseGlobal); !responseGlobal.Success {
		return
	}

}
