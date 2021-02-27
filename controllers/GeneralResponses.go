package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSONSwagger struct {
}

type GeneralResponse struct {
	Messages []string `json:"messages" example:"error description,other error description"`
}

func badRequest(c *gin.Context, messages []string) {
	c.JSON(http.StatusBadRequest,
		GeneralResponse{Messages: messages})
}

func serverError(c *gin.Context, err error) {
	_ = c.Error(err)
	c.JSON(http.StatusInternalServerError,
		GeneralResponse{Messages: []string{"We are working to improve the flow of this request."}})
}
