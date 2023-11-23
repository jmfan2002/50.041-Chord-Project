package entrypoint

import (
	"EntryNode/util"
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

	sampleStruct := GetValueResBody{
		Key:   key,
		Value: val,
	}

	util.WriteSuccessResponse(w, &sampleStruct)
}
