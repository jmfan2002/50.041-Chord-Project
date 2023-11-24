package api

// Path: ServerNode/api/GetSuccessors.go
// Compare this snippet from EntryNode/entrypoint/GetSuccessors.go:
// package entrypoint
//
import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GetHashTableResBody struct {
	HashTable []TableEntry `json:"hashTable"`
}

type TableEntry struct {
	NodeAdress string `json:"node"`
	Key        string `json:"key"`
	Value      string `json:"value"`
}

func (h *Handler) GetHashTable(w http.ResponseWriter, r *http.Request) {
	hashTable := []TableEntry{}

	for _, value := range h.NodeInfo.NodeContents {
		hashTable = append(hashTable, TableEntry{
			NodeAdress: h.NodeInfo.NodeUrl,
			Key:        value.Key,
			Value:      value.Value,
		})
	}

	responseData := GetHashTableResBody{
		HashTable: hashTable,
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
