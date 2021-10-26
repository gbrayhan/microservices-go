package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSONSwagger struct {
}

type generalResponse struct {
	Messages []string `json:"messages" example:"error description,other error description"`
}

func BadRequest(c *gin.Context, messages []string) {
	c.JSON(http.StatusBadRequest,
		generalResponse{Messages: messages})
}

func ServerError(c *gin.Context, err error) {
	_ = c.Error(err)
	c.JSON(http.StatusInternalServerError,
		generalResponse{Messages: []string{"We are working to improve the flow of this request."}})
}
