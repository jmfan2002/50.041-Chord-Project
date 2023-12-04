package entrypoint

import (
	"EntryNode/util"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HealthResBody struct {
	Val string `json:"val"`
}

func (entryPoint *EntryPoint) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Health check called")

	queryParams := r.URL.Query()
	serverAddress := queryParams.Get("node")

	if serverAddress == "" {
		fmt.Println("No server address provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Send health request to server
	resp, err := http.Get(serverAddress + "/api/health")
	if err != nil {
		fmt.Printf("An error occurred %s\n", err.Error())
		util.WriteSuccessResponse(w, &HealthResBody{Val: "false"})
		return
	}

	fmt.Printf("Request result: %s\n", resp.Status)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading server node response")
		return
	}

	response := HealthResBody{}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		fmt.Println("Error parsing request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	util.WriteSuccessResponse(w, &response)
}
