package entrypoint

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GetNodesResBody struct {
	Nodes []string `json:"nodes"`
}

func (handler *Handler) GetNodes(w http.ResponseWriter, r *http.Request) {
	nodeList := []string{}

	// Go through all nodes
	for _, nodeURL := range handler.EntryPoint.servers {
		nodeList = append(nodeList, nodeURL)
	}

	responseData := GetNodesResBody{
		Nodes: nodeList,
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
