package controllers

import (
	"bytes"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func bindJSON(c *gin.Context, request interface{}) (err error) {
	buf := make([]byte, 5120)
	num, _ := c.Request.Body.Read(buf)
	reqBody := string(buf[0:num])
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(reqBody)))
	err = c.ShouldBindJSON(request)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(reqBody)))
	return
}
