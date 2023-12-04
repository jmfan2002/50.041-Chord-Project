package entrypoint

import (
	"EntryNode/util"
	"encoding/json"
	"fmt"
	"net/http"
)

type GetValueResBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (entryPoint *EntryPoint) GetValue(w http.ResponseWriter, r *http.Request) {
	// Get key from request
	queryParams := r.URL.Query()

	key := queryParams.Get("key")
	val := entryPoint.getKVP(key)

	response := &GetValueResBody{}

	err := json.Unmarshal([]byte(val), response)
	if err != nil {
		fmt.Println("Error parsing request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	util.WriteSuccessResponse(w, &response)
}
