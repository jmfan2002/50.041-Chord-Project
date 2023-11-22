package api

import (
	"ServerNode/structs"
	"ServerNode/util"
	"net/http"
)

func (h *Handler) GetSuccessors(w http.ResponseWriter, r *http.Request) {
	util.WriteResponse(w, structs.SuccessorsResponse{Successors: h.NodeInfo.SuccessorArray}, http.StatusOK)
}
