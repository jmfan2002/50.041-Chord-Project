package api

import (
	"ServerNode/structs"

	"encoding/json"
	"fmt"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Health check called")

	sampleStruct := structs.SampleStruct{
		Val: "success!",
	}

	response, err := json.Marshal(sampleStruct)
	if err != nil {
		fmt.Println("error marshalling data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
