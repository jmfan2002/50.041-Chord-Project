package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
)

type SetValueBody struct {
	ValueHash string `json:"valueHash"`
	Data      string `json:"data"`
}

type EntryPoint struct {
	ipHashes []big.Int

	// Maps hashes as strings, to ips as strings
	servers map[string]string
}

func (entryPoint *EntryPoint) SetValue(key string, val string) {
	serverAddress := entryPoint.Lookup(key)

	h := sha256.New()
	h.Write([]byte(key))
	b := h.Sum(nil)

	j, _ := json.Marshal(SetValueBody{
		hex.EncodeToString(b), val,
	})

	res, err := http.Post(serverAddress,
		"application/json",
		bytes.NewBuffer(j),
	)
	if err != nil {
		fmt.Printf("An error occurred %s\n", err.Error())
	}
	fmt.Printf("Request result: %s\n", res.Status)
}

func (entryPoint *EntryPoint) GetValue(key string) {
	serverAddress := entryPoint.Lookup(key)

	h := sha256.New()
	h.Write([]byte(key))
	b := h.Sum(nil)

	resp, err := http.Get(serverAddress + "/api?valueHash=" + hex.EncodeToString(b))
	if err != nil {
		fmt.Printf("An error occurred %s\n", err.Error())
	}
	defer resp.Body.Close()

	fmt.Printf("Request result: %s\n", resp.Status)
}

// Returns the IP address of the node responsible for the given key
func (entryPoint *EntryPoint) Lookup(key string) string {
	h := sha256.New()
	h.Write([]byte(key))
	z := big.NewInt(0)
	z.SetBytes(h.Sum(nil))

	serverIpHash := entryPoint.ipHashes[binarySearch(entryPoint.ipHashes, z)]

	serverIp := entryPoint.servers[serverIpHash.String()]

	return serverIp
}

func AddNode(ipAddress string) {

}

func binarySearch(arr []big.Int, item *big.Int) int {
	low := 0
	high := len(arr) - 1
	for low <= high {
		mid := low + (high-low)/2
		cmp := arr[mid].Cmp(item)
		if cmp == -1 {
			low = mid + 1
		} else if cmp == 1 {
			high = mid - 1
		} else {
			return mid
		}
	}

	if low == len(arr) {
		return 0
	}
	return low
}
