package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ReadRequestBody(w http.ResponseWriter, r *http.Request, reqBody any) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(bodyBytes, reqBody)
	if err != nil {
		fmt.Println("Error parsing request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func WriteSuccessResponse(w http.ResponseWriter, body any) {
	response, err := json.Marshal(body)
	if err != nil {
		fmt.Println("error marshalling data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
