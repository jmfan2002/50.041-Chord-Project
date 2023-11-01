package entrypoint

import (
	"EntryNode/util"
	"net/http"
)

type SetValueReqBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SetValueResBody struct {
	Data string
}

func (entryPoint *EntryPoint) SetValue(w http.ResponseWriter, r *http.Request) {
	// Get key and value from request
	var reqBody SetValueReqBody

	util.ReadRequestBody(w, r, &reqBody)

	entryPoint.setKVP(reqBody.Key, reqBody.Value)

	sampleStruct := SetValueResBody{
		Data: "test123",
	}

	util.WriteSuccessResponse(w, &sampleStruct)
}
