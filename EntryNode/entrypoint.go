package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"
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
	fmt.Printf("Sending %s to %s\n", key, serverAddress)

	/*
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
	*/
}

func (entryPoint *EntryPoint) GetValue(key string) {
	serverAddress := entryPoint.Lookup(key)
	fmt.Printf("Getting %s from %s\n", key, serverAddress)

	/*
		h := sha256.New()
		h.Write([]byte(key))
		b := h.Sum(nil)

		resp, err := http.Get(serverAddress + "/api?valueHash=" + hex.EncodeToString(b))
		if err != nil {
			fmt.Printf("An error occurred %s\n", err.Error())
		}
		defer resp.Body.Close()

		fmt.Printf("Request result: %s\n", resp.Status)
	*/
}

// Returns the IP address of the node responsible for the given key
func (entryPoint *EntryPoint) Lookup(key string) string {
	h := sha256.New()
	h.Write([]byte(key))
	z := big.NewInt(0)
	z.SetBytes(h.Sum(nil))

	serverIdx := binarySearch(entryPoint.ipHashes, z) % len(entryPoint.ipHashes)
	serverIpHash := entryPoint.ipHashes[serverIdx]

	serverIp := entryPoint.servers[serverIpHash.Text(16)]

	return serverIp
}

func (entryPoint *EntryPoint) AddNode(ipAddress string) {
	h := sha256.New()
	h.Write([]byte(ipAddress))
	ipHash := h.Sum(nil)
	z := big.NewInt(0)
	z.SetBytes(ipHash)

	insertionPoint := binarySearch(entryPoint.ipHashes, z)

	if len(entryPoint.ipHashes) == insertionPoint {
		entryPoint.ipHashes = append(entryPoint.ipHashes, *z)
	} else {
		entryPoint.ipHashes = append(
			entryPoint.ipHashes[:insertionPoint+1],
			entryPoint.ipHashes[insertionPoint:]...,
		)
		entryPoint.ipHashes[insertionPoint] = *z
	}

	entryPoint.servers[z.Text(16)] = ipAddress

	fmt.Printf("Added node at %s with hash %x\n", ipAddress, ipHash)
	fmt.Println("Current node list:")
	for idx, item := range entryPoint.ipHashes {
		fmt.Printf("\t%s\n", item.Text(16))

		if idx < len(entryPoint.ipHashes)-1 {
			zz := entryPoint.ipHashes[idx+1]
			fmt.Printf("\t %d\n", item.Cmp(&zz))
		}
	}
	fmt.Printf("Current node map: %v\n", entryPoint.servers)
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

	return low
}
