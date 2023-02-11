package middlewares

import (
	"bytes"
	"fmt"
	"github.com/gbrayhan/microservices-go/application/services"
	"io"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gbrayhan/microservices-go/utils"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func GinBodyLogMiddleware(c *gin.Context) {
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	c.Next()

	buf := make([]byte, 4096)
	num, err := c.Request.Body.Read(buf)
	if err != nil && err.Error() != "EOF" {

		_ = fmt.Errorf("error reading buffer: %s", err.Error())
	}
	reqBody := string(buf[0:num])
	c.Request.Body = io.NopCloser(bytes.NewBuffer([]byte(reqBody)))

	loc, _ := time.LoadLocation("America/Mexico_City")
	allDataIO := map[string]interface{}{
		"ruta":          c.FullPath(),
		"request_uri":   c.Request.RequestURI,
		"raw_request":   reqBody,
		"status_code":   c.Writer.Status(),
		"body_response": blw.body.String(),
		"errors":        c.Errors.Errors(),
		"created_at":    time.Now().In(loc).Format("2006-01-02T15:04:05"),
	}
	_ = fmt.Sprintf("%v", allDataIO)

	allLogs := []string{
		"/payment-with-recurrence",
		"/buy-console",
		"/other-route",
	}

	if existAll, _ := utils.InArray(c.FullPath(), allLogs); existAll {
		if c.Writer.Status() == 500 {
			go func() { err = services.SendSimpleMail() }()
		}
		// go SaveLogs(allDataIO)
	}
}
