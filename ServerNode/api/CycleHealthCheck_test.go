package api_test

import (
	"ServerNode/api"

	"net/http"
	"net/http/httptest"
	"testing"
	"errors"
	"io"
	"bytes"

	"github.com/gorilla/mux"
)

type MockRequester struct {
	Response *http.Response
	Error    error
}

func (m *MockRequester) SendRequest(baseUrl, endpoint, httpMethod string, timeoutMs int) (*http.Response, error) {
	return m.Response, m.Error
}

func TestCycleHealthCheck(t *testing.T) {
	handler := api.NewHandler("", 5)
	handler.NodeInfo.NodeHash = "b"

	r := mux.NewRouter()
	r.HandleFunc("/cycleHealth/{StartingNodeHash}/{FinishedLoop}", handler.CycleHealthCheck).Methods("GET")

	t.Run("After looping, we reach the original node. we should return ok", func(t *testing.T) {
		RunHttpTest(r, t, "/cycleHealth/b/true", http.StatusOK, ``)
	})

	t.Run("After looping, we reach past original node, we should return ok", func(t *testing.T) {
		handler.NodeInfo.NodeHash = "c"
		RunHttpTest(r, t, "/cycleHealth/b/true", http.StatusOK, ``)
	})

	t.Run("Loop but don't reach original, call to next node succeeds", func(t *testing.T) {
		handler.NodeInfo.NodeHash = "a"
		handler.NodeInfo.SuccessorArray = []string{"b","c"}
		handler.Requester = &MockRequester{
			Response: &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBufferString("hello"))},
			Error:    nil,
		}
		RunHttpTest(r, t, "/cycleHealth/b/true", http.StatusOK, ``)
	})

	t.Run("No loop, call to next node times out", func(t *testing.T) {
		handler.NodeInfo.NodeHash = "a"
		handler.NodeInfo.SuccessorArray = []string{"b","c"}
		handler.Requester = &MockRequester{
			Response: nil,
			Error: errors.New("timeout"),
		}
		RunHttpTest(r, t, "/cycleHealth/b/false", http.StatusInternalServerError, ``)
	})
}

// Fails the test if failure, does nothing on pass
func RunHttpTest(router *mux.Router, t *testing.T, url string, expectedCode int, expectedContents string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != expectedCode {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := expectedContents
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
