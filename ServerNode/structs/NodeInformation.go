package structs

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type NodeInformation struct {
	NodeUrl          string
	NodeHash         string
	NodeContents     map[string]string
	PredecessorArray []string
	SuccessorArray   []string
	StoredNbrs       int
}

func (n NodeInformation) String() string {
	return fmt.Sprintf("{%s - PredecessorArray: %s, SuccessorArray: %s, StoredNbrs: %d, NodeHash: %s}", n.NodeUrl, n.PredecessorArray, n.SuccessorArray, n.StoredNbrs, n.NodeHash)
}

func NewNodeInformation(nodeUrl string, storedNbrs int) NodeInformation {
	hash := sha256.Sum256([]byte(nodeUrl))
	hashStr := hex.EncodeToString(hash[:])
	return NodeInformation{
		NodeUrl:          nodeUrl,
		NodeHash:         hashStr,
		NodeContents:     make(map[string]string),
		PredecessorArray: make([]string, 0),
		SuccessorArray:   make([]string, 0),
		StoredNbrs:       storedNbrs,
	}
}
