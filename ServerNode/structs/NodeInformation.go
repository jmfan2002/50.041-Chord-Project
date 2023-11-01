package structs

import (
	"crypto/sha256"
	"encoding/hex"
)

type NodeInformation struct {
	NodeUrl string
	NodeHash string
	NodeContents map[string]string
}

func NewNodeInformation(nodeUrl string) NodeInformation {
	hash := sha256.Sum256([]byte(nodeUrl))
	hashStr := hex.EncodeToString(hash[:])
	return NodeInformation{
		NodeUrl: nodeUrl,
		NodeHash: hashStr,
		NodeContents: make(map[string]string),
	}
}
