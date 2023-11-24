package entrypoint

import (
	"EntryNode/util"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CycleHealthResBody struct {
	Message string `json:"message"`
}

type CycleHealthResponse struct {
	CycleSize int `json:"cycleSize"`
}

func (entryPoint *EntryPoint) CycleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Cycling health check")
	// Checks if there are any nodes in the network
	var resp *http.Response
	var err error
	for _, nodeAdress := range entryPoint.servers {
		resp, err = http.Get(nodeAdress + "/api/cycleHealth/nil/")
		if err == nil {
			break
		}
	}

	if err != nil {
		fmt.Println("error getting data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := CycleHealthResponse{}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		fmt.Println("Error parsing request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	util.WriteSuccessResponse(w, &response)
}
