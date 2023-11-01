package entrypoint

import (
	"EntryNode/util"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type EntryPoint struct {
	ipHashes []big.Int

	// Maps hashes as strings, to ips as strings
	servers map[string]string
}

func New() *EntryPoint {
	return &EntryPoint{
		make([]big.Int, 0),
		make(map[string]string),
	}
}

func (entryPoint *EntryPoint) setKVP(key string, val string) {
	serverAddress := entryPoint.Lookup(key)
	fmt.Printf("Sending %s -> %s to %s\n", key, val, serverAddress)

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

func (entryPoint *EntryPoint) getKVP(key string) string {
	serverAddress := entryPoint.Lookup(key)
	fmt.Printf("Getting %s from %s\n", key, serverAddress)

	return "a"
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
// This is pretty much the only part the Chord paper exposes in their API,
// but we add some other functions for a nicer user application API.
func (entryPoint *EntryPoint) Lookup(key string) string {
	h := sha256.New()
	h.Write([]byte(key))
	z := big.NewInt(0)
	z.SetBytes(h.Sum(nil))

	serverIdx := util.BinarySearch(entryPoint.ipHashes, z) % len(entryPoint.ipHashes)
	serverIpHash := entryPoint.ipHashes[serverIdx]

	serverIp := entryPoint.servers[serverIpHash.Text(16)]

	return serverIp
}

func (entryPoint *EntryPoint) addServer(ipAddress string) {
	h := sha256.New()
	h.Write([]byte(ipAddress))
	ipHash := h.Sum(nil)
	z := big.NewInt(0)
	z.SetBytes(ipHash)

	insertionPoint := util.BinarySearch(entryPoint.ipHashes, z)

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
	for _, item := range entryPoint.ipHashes {
		fmt.Printf("\t%s->%s\n", item.Text(16), entryPoint.servers[item.Text(16)])
	}
}
