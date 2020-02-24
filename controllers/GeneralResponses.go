package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)


func BadRequest(c *gin.Context, messages []string) {
	c.JSON(http.StatusBadRequest, struct {
		Messages []string `json:"messages"`
	}{Messages: messages})
}

func ServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, []byte(`{
		"messages": ["We are working to improve the flow of this request."]
	}`))
}
