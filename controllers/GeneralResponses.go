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

func ServerError(c *gin.Context, err error) {
	c.Error(err)
	c.JSON(http.StatusInternalServerError, struct {
		Messages []string `json:"messages"`
	}{
		Messages: []string{"We are working to improve the flow of this request."},
	})
}
