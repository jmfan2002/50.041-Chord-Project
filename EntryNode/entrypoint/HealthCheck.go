package entrypoint

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HealthResBody struct {
	Message string `json:"message"`
}

func (entryPoint *EntryPoint) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Health check called")

	queryParams := r.URL.Query()
	serverAddress := queryParams.Get("node")

	// Send health request to server
	resp, err := http.Get(serverAddress + "/api/health")
	if err != nil {
		fmt.Printf("An error occurred %s\n", err.Error())
	}

	fmt.Printf("Request result: %s\n", resp.Status)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading server node response")
		return
	}

	responseData := HealthResBody{
		Message: string(bodyBytes),
	}

	response, err := json.Marshal(responseData)
	if err != nil {
		fmt.Println("error marshalling data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
