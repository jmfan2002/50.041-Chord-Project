package api

import (
	"ServerNode/structs"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetValue(w http.ResponseWriter, r *http.Request, nodeData *structs.NodeInformation) {
	fmt.Println("Get value called")

	response, err := json.Marshal(nodeData)
	if err != nil {
		fmt.Println("error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response = []byte("[Error] need to get specific value from map")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
