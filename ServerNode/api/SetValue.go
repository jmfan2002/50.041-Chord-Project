package api

import (
	"ServerNode/structs"
	"encoding/json"
	"fmt"
	"net/http"
)

func SetValue(w http.ResponseWriter, r *http.Request, nodeData *structs.NodeInformation) {
	fmt.Println("Set value called")

	response, err := json.Marshal(nodeData)
	if err != nil {
		fmt.Println("error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	nodeData.NodeContents = response

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
