package tools

import (
	"bytes"
	"context"
	"net/http"
	"time"
)

func SendRequestJson(url string, jsonData []byte, timeOutMs float64) (response *http.Response, err error) {
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(timeOutMs))

	request = request.WithContext(ctx)
	client := &http.Client{}

	return client.Do(request)
}

