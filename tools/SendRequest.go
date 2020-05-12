package tools

import (
	"context"
	"net/http"
	"strings"
	"time"
)

func SendRequest(method string, url string, headers map[string]string, body string) (response *http.Response, err error) {
	request, _ := http.NewRequest(method, url, strings.NewReader(body))

	for key, value := range headers {
		request.Header.Set(key, value)
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*1000000)
	request = request.WithContext(ctx)
	client := &http.Client{}

	return client.Do(request)
}

