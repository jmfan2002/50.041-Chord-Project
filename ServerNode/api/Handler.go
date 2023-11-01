package api

import (
	"ServerNode/structs"
)

type Handler struct {
	NodeInfo *structs.NodeInformation
}

func NewHandler(nodeUrl string, storedNbrs int) *Handler {
	return &Handler{NodeInfo: structs.NewNodeInformation(nodeUrl, storedNbrs)}
}