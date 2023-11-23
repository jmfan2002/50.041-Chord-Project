package entrypoint

import (
	"EntryNode/util"
	"fmt"
	"net/http"
)

type SetValueReqBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SetValueResBody struct {
	Message string `json:"message"`
}

func (entryPoint *EntryPoint) SetValue(w http.ResponseWriter, r *http.Request) {
	// Get key and value from request
	var reqBody SetValueReqBody
	fmt.Println("Got data set")

	util.ReadRequestBody(w, r, &reqBody)

	entryPoint.setKVP(reqBody.Key, reqBody.Value)

	response := SetValueResBody{
		Message: reqBody.Key + " set to " + reqBody.Value,
	}

	util.WriteSuccessResponse(w, &response)
}
