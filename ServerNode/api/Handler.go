package api

import (
	"ServerNode/structs"
	"ServerNode/util"
)

type Handler struct {
	NodeInfo *structs.NodeInformation
	Requester util.Requester
}

func NewHandler(nodeUrl string, storedNbrs int) *Handler {
	return &Handler{
		NodeInfo: structs.NewNodeInformation(nodeUrl, storedNbrs),
		Requester: &util.HeartbeatRequester{},
	}
}
