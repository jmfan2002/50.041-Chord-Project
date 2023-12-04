package util

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Requester is used to define how a node sends an external request
// Note: this may need to be modified to accept a body
type Requester interface {
	SendRequest(baseUrl, endpoint, httpMethod string, contents interface{}, timeoutMs int) (*http.Response, error)
}

// BasicRequester provides an implementation of Requester
type BasicRequester struct{}

func (r *BasicRequester) SendRequest(baseUrl, endpoint, httpMethod string, contents interface{}, timeoutMs int) (*http.Response, error) {
	client := http.Client{
		Timeout: time.Duration(timeoutMs) * time.Millisecond,
	}
	return client.Get(baseUrl + endpoint)
}

// HeartbeatRequester prevents a node from timing out so long as it responds to a heartbeat request
type HeartbeatRequester struct{}

func (r *HeartbeatRequester) SendRequest(baseUrl, endpoint, httpMethod string, contents interface{}, timeoutMs int) (*http.Response, error) {
	// fmt.Printf("[Debug] sending %s request to: %s%s with heartbeat timeout %d\n", httpMethod, baseUrl, endpoint, timeoutMs)

	// Create a new context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new request
	var requestBuffer *bytes.Buffer
	requestBuffer = &bytes.Buffer{}
	if contents != nil {
		jsonValue, err := json.Marshal(contents)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("[ERROR] failed to marshall JSON body: %s", err))
		}
		requestBuffer = bytes.NewBuffer(jsonValue)
	}

	req, err := http.NewRequestWithContext(ctx, httpMethod, baseUrl+endpoint, requestBuffer)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[ERROR] failed to create request: %s", err))
	}
	if contents != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Send health checks at interval of timeoutMs in a goroutine to allow us to progress to main request
	go func() {
		// fmt.Printf("[Htbt] starting heartbeat checks\n")
		ticker := time.NewTicker(time.Duration(timeoutMs) * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if ctx.Err() != nil {
					// fmt.Println("[Htbt] Heartbeat cancelled")
					return
				}

				client := http.Client{
					Timeout: time.Duration(timeoutMs) * time.Millisecond,
				}

				healthResp, healthErr := client.Get(baseUrl + "/api/health")
				if healthErr != nil || healthResp.StatusCode != http.StatusOK {
					// fmt.Printf("[Htbt] Health check failed, canceling the request. err: %s\n", healthErr)
					cancel()
					return
				} else {
					// fmt.Println("[Htbt] heartbeat health check successful!")
				}

			case <-ctx.Done():
				// If the context is done, stop the checks
				// fmt.Printf("[Htbt] stopping heartbeat checks\n")
				return
			}
		}
	}()

	// Send the main request, cancelled with the cancel() method
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[Warning] failed to send request: %s", err))
	}
	return resp, nil
}
