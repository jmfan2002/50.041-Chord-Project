package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SetSuccReqBody struct {
	Succ string `json:"succ"`
}

func (h *Handler) SetSucc(w http.ResponseWriter, r *http.Request) {
	// Get key and value from request
	var reqBody SetSuccReqBody

	readRequestBody(w, r, &reqBody)

	h.NodeInfo.SuccessorArray[0] = reqBody.Succ

	fmt.Println("Set successors called, new successor: %s", reqBody.Succ)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func readRequestBody(w http.ResponseWriter, r *http.Request, reqBody any) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(string(bodyBytes))

	err = json.Unmarshal(bodyBytes, reqBody)
	if err != nil {
		fmt.Println("Error parsing request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
