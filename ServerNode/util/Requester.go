package util

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Requester is used to define how a node sends an external request
// Note: this may need to be modified to accept a body
type Requester interface {
	SendRequest(baseUrl, endpoint, method string, timeoutMs int) (*http.Response, error)
}

// BasicRequester provides an implementation of Requester
type BasicRequester struct{}

func (r *BasicRequester) SendRequest(baseUrl, endpoint, method string, timeoutMs int) (*http.Response, error) {
	client := http.Client{
		Timeout: time.Duration(timeoutMs) * time.Millisecond,
	}
	return client.Get(baseUrl + endpoint)
}

// HeartbeatRequester prevents a node from timing out so long as it responds to a heartbeat request
type HeartbeatRequester struct{}

func (r *HeartbeatRequester) SendRequest(baseUrl, endpoint, httpMethod string, timeoutMs int) (*http.Response, error) {
	fmt.Printf("[Debug] sending %s request to: %s%s with heartbeat timeout %d\n", httpMethod, baseUrl, endpoint, timeoutMs)

	// Create a new context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Make a new request
	req, err := http.NewRequestWithContext(ctx, httpMethod, baseUrl + endpoint, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[ERROR] failed to create request: %s", err))
	}

	
	// Send health checks at interval of timeoutMs in a goroutine to allow us to progress to main request
	ticker := time.NewTicker(time.Duration(timeoutMs) * time.Millisecond)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ticker.C:
				client := http.Client {
					Timeout: time.Duration(timeoutMs) * time.Millisecond,
				}

				resp, err := client.Get(baseUrl + "/api/health")
				if err != nil || resp.StatusCode != http.StatusOK {
					fmt.Println("[Warning] Health check failed, canceling the request")
					cancel()
					return
				} else {
					fmt.Println("[Debug] heartbeat health check successful!")
				}
			case <-ctx.Done():
				// If the context is done, stop the health checks
				return
			}
		}
	}()

	// Send the main request, cancelled with the cancel() method
	fmt.Printf("[Debug] Request start\n")
	resp, err := http.DefaultClient.Do(req)
	fmt.Printf("[Debug] Request finish\n")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[Warning] failed to send request: %s", err))
	}
	fmt.Println("[Debug] Successfully received response!")
	return resp, nil
}