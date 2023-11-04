package util

import (
	"net/http"
	"time"
)

// Requester is used to define how a node sends an external request
// Note: this may need to be modified to accept a body
type Requester interface {
	SendRequest(requestID string, timeoutMs int) (*http.Response, error)
}

// BasicRequester provides an implementation of Requester
type BasicRequester struct{}

func (r *BasicRequester) SendRequest(requestUrl string, timeoutMs int) (*http.Response, error) {
	client := http.Client{
		Timeout: time.Duration(timeoutMs) * time.Millisecond,
	}
	return client.Get(requestUrl)
}
