package api

import (
	"ServerNode/structs"

	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (h *Handler) SetValue(w http.ResponseWriter, r *http.Request) {
	// Get key and value from request
	var reqBody structs.SetValueReqBody

	fmt.Println("Set value called")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(string(bodyBytes))

	err = json.Unmarshal(bodyBytes, &reqBody)
	if err != nil {
		fmt.Println("Error parsing request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(reqBody)

	h.NodeInfo.NodeContents[reqBody.ValueHash] = reqBody.Data

	sampleStruct := structs.SetValueReqBody{
		ValueHash: "",
		Data:      "",
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
