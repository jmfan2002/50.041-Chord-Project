package entrypoint

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

func (entryPoint *EntryPoint) GetHashTable(w http.ResponseWriter, r *http.Request) {
	hashTable := []TableEntry{}

	// Go through all nodes
	for _, nodeURL := range entryPoint.servers {
		// Ask for hashT table
		resp, err := http.Get(nodeURL + "/api/hashTable")
		if err != nil {
			fmt.Println("error getting data")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Decode response
		var data GetHashTableResBody
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			fmt.Println("error decoding data")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Add to hashTable
		hashTable = append(hashTable, data.HashTable...)

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
