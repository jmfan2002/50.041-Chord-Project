package entrypoint

import (
	"EntryNode/util"
	"fmt"
	"net/http"
)

type JoinReqBody struct {
	NewNodeAddress string `json:"NewNodeAddress"`
}

type JoinResBody struct {
	Data string
}

func (entryPoint *EntryPoint) AddNode(w http.ResponseWriter, r *http.Request) {
	// Get key and value from request
	var reqBody JoinReqBody

	util.ReadRequestBody(w, r, &reqBody)

	fmt.Println("Receiving join for", reqBody.NewNodeAddress)
	entryPoint.addServer(reqBody.NewNodeAddress)

	sampleStruct := JoinResBody{
		Data: "test123",
	}

	util.WriteSuccessResponse(w, &sampleStruct)

	fmt.Println("Ending join")
}
