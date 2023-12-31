package structs

import (
	"ServerNode/util"
	"fmt"
)

type NodeInformation struct {
	NodeUrl          string
	NodeHash         string
	NodeContents     map[string]EntryResponse
	PredecessorArray []string
	SuccessorArray   []string
	StoredNbrs       int

	// Heartbeat information to detect timeouts
	HeartbeatTarget string
	HeartbeatSuccessful bool
}

func (n NodeInformation) String() string {
	return fmt.Sprintf("{%s - PredecessorArray: %s, SuccessorArray: %s, StoredNbrs: %d, NodeHash: %s}", n.NodeUrl, n.PredecessorArray, n.SuccessorArray, n.StoredNbrs, n.NodeHash)
}

func NewNodeInformation(nodeUrl string, storedNbrs int) *NodeInformation {
	hashStr := util.Sha256String(nodeUrl)
	return &NodeInformation{
		NodeUrl:          nodeUrl,
		NodeHash:         hashStr,
		NodeContents:     make(map[string]EntryResponse),
		PredecessorArray: make([]string, 0),
		SuccessorArray:   make([]string, 0),
		StoredNbrs:       storedNbrs,
	}
}
