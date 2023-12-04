package entrypoint

import (
	"EntryNode/util"
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

	// fmt.Println("[Msg] Get hash table called")
	// fmt.Println("Getting hash table from", len(handler.EntryPoint.servers), "nodes")

	// Go through all nodes
	for _, nodeURL := range entryPoint.Servers {
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

	response := GetHashTableResBody{
		HashTable: hashTable,
	}

	util.WriteSuccessResponse(w, &response)
}
