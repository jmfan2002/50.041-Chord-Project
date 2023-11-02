package api_test

import (
	"ServerNode/api"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestCycleHealthCheck(t *testing.T) {
	handler := api.NewHandler("b", 5)
	r := mux.NewRouter()
	r.HandleFunc("/cycleHealth/{StartingNodeHash}/{FinishedLoop}", handler.CycleHealthCheck).Methods("GET")

	// After looping, we reach the original node
	RunHttpTest(r, t, "/cycleHealth/b/true", http.StatusOK, ``)
	// After looping, we reach past original node
	RunHttpTest(r, t, "/cycleHealth/c/true", http.StatusOK, ``)
	// After looping, we haven't reached original
	RunHttpTest(r, t, "/cycleHealth/a/true", http.StatusOK, ``)

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
