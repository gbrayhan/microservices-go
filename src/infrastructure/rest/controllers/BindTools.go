package controllers

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
)

func BindJSON(c *gin.Context, request any) error {
	buf := make([]byte, 5120)
	num, _ := c.Request.Body.Read(buf)
	reqBody := string(buf[0:num])
	c.Request.Body = io.NopCloser(bytes.NewBuffer([]byte(reqBody)))
	err := c.ShouldBindJSON(request)
	c.Request.Body = io.NopCloser(bytes.NewBuffer([]byte(reqBody)))
	return err
}

func BindJSONMap(c *gin.Context, request *map[string]any) error {
	buf := make([]byte, 5120)
	num, _ := c.Request.Body.Read(buf)
	reqBody := buf[0:num]
	c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	err := json.Unmarshal(reqBody, &request)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	return err
}

// MessageResponse ...
type MessageResponse struct {
	Message string `json:"message"`
}

type SortByDataRequest struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

type FieldDateRangeDataRequest struct {
	Field     string `json:"field"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}
