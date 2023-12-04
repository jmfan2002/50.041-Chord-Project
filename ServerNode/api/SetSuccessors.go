package api

import (
	"ServerNode/structs"
	"ServerNode/util"
	"fmt"
	"net/http"
)

func (h *Handler) SetSuccessors(w http.ResponseWriter, r *http.Request) {
	var newSuccessors = &structs.SuccessorsResponse{}
	err := util.ReadBody(r.Body, &newSuccessors)
	if err != nil {
		fmt.Printf("[ERROR] failed to read request: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// fmt.Printf("[Msg] setting successors to %s\n", newSuccessors.Successors)
	h.NodeInfo.SuccessorArray = newSuccessors.Successors
}
