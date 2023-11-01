package util

import (
	"net/http"
	"time"
)

func SendRequest(requestUrl string, timeoutMs int) (*http.Response, error){
	client := http.Client {
		Timeout: time.Duration(timeoutMs) * time.Millisecond,
	}
	return client.Get(requestUrl)
}
