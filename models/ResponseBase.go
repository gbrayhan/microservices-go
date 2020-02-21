package models

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseBase struct {
	BadRequest     bool     `json:"-"`
	FatalError     bool     `json:"-"`
	Success        bool     `json:"-"`
	MessagesClient []string `json:"messages"`
	MessagesFatal  []string `json:"-"`
	Language       string   `json:"-"`
}

func (resp *ResponseBase) ShowResponseJSON(c *gin.Context, responseClient interface{}) {
	type Data interface{}

	if resp.FatalError {
		// TODO: Add action to send notification internal server
		c.JSON(http.StatusInternalServerError, []byte(`
			"messages": ["We are working to improve the flow of this request."]
		`))
		return
	}

	if resp.BadRequest {
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, responseClient)
}



// Return status code 200 Status OK
func (resp *ResponseBase) AppendSuccessMessage(message string) {
	resp.MessagesClient = append(resp.MessagesClient, message)
	resp.Success = true
}

// Return an error 400 Bad Request
func (resp *ResponseBase) AppendErrorRequest(message string) {
	resp.MessagesClient = append(resp.MessagesClient, message)
	resp.BadRequest = true
	resp.Success = false
}

// Return an error 500 Internal Server Error
func (resp *ResponseBase) AppendFatalError(message string) {
	resp.MessagesFatal = append(resp.MessagesFatal, message)
	resp.FatalError = true
	resp.Success = false
}
