package entrypoint

import (
	"EntryNode/util"
	"crypto/sha256"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
)

type EntryPoint struct {
	ipHashes []big.Int

	// Maps hashes as strings, to ips as strings
	servers map[string]string

	// number of faults to tolerate
	toleratedFaults int

	// used to send heartbeat requests
	requester util.HeartbeatRequester
}

func New(k int) *EntryPoint {
	return &EntryPoint{
		ipHashes:        make([]big.Int, 0),
		servers:         make(map[string]string),
		toleratedFaults: k,
		requester:       util.HeartbeatRequester{},
	}
}

type ChordSetValueReq struct {
	Key   string
	Value string
	Nonce string
}

func (entryPoint *EntryPoint) trySetKVP(key string, nonce string, val string) {
	serverAddress := entryPoint.Lookup(key, nonce)
	fmt.Printf("Sending %s -> %s to %s\n", key, val, serverAddress)

	res, err := entryPoint.requester.SendRequest(
		serverAddress,
		"/api",
		http.MethodPost,
		ChordSetValueReq{
			key,
			val,
			nonce,
		},
		util.REQUEST_TIMEOUT,
	)

	if err != nil {
		fmt.Printf("An error occurred %s\n", err.Error())
		h := sha256.New()
		h.Write([]byte(serverAddress))
		z := big.NewInt(0)
		z.SetBytes(h.Sum(nil))
		for i := 0; i < len(entryPoint.ipHashes); i += 1 {
			if entryPoint.ipHashes[i].Cmp(z) == 0 {
				entryPoint.ipHashes = append(entryPoint.ipHashes[:i], entryPoint.ipHashes[i+1:]...)
				break
			}
		}
		return
	}
	fmt.Printf("Request result: %s\n", res.Status)
}

func (entryPoint *EntryPoint) setKVP(key string, val string) {
	for i := 0; i <= entryPoint.toleratedFaults; i += 1 {
		go entryPoint.trySetKVP(key, strconv.Itoa(i), val)
	}
}

type ChordGetValueRes struct {
	NodeContents map[string]string
}

type Res struct {
	data    string
	isError bool
}

func (entryPoint *EntryPoint) tryGetKVP(serverAddress string, key string, nonce string, out chan Res) {

	resp, err := entryPoint.requester.SendRequest(
		serverAddress,
		fmt.Sprintf("/api/%s/%s", key, nonce),
		http.MethodGet,
		nil,
		util.REQUEST_TIMEOUT,
	)

	if err != nil {
		fmt.Printf("An error occurred %s\n", err.Error())

		// Attempt to store the key elsewhere
		h := sha256.New()
		h.Write([]byte(serverAddress))
		z := big.NewInt(0)
		z.SetBytes(h.Sum(nil))
		for i := 0; i < len(entryPoint.ipHashes); i += 1 {
			if entryPoint.ipHashes[i].Cmp(z) == 0 {
				entryPoint.ipHashes = append(entryPoint.ipHashes[:i], entryPoint.ipHashes[i+1:]...)
				break
			}
		}

		out <- Res{nonce, true}
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading server node response")
		out <- Res{"", true}
		return
	}
	out <- Res{
		string(bodyBytes),
		false,
	}
}

func (entryPoint *EntryPoint) getKVP(key string) string {
	in := make(chan Res)

	for i := 0; i <= entryPoint.toleratedFaults; i += 1 {
		// Find the server that should have the key
		serverAddress := entryPoint.Lookup(key, strconv.Itoa(i))
		fmt.Printf("Getting %s from %s\n", key, serverAddress)

		go entryPoint.tryGetKVP(serverAddress, key, strconv.Itoa(i), in)
	}

	errors := make([]string, 0)
	out := ""
	for i := 0; i < entryPoint.toleratedFaults; i += 1 {
		next := <-in
		if !next.isError {
			out = next.data
		} else {
			errors = append(errors, next.data)
		}
	}

	for _, nonce := range errors {
		entryPoint.trySetKVP(key, nonce, out)
	}

	return out
}

// Returns the IP address of the node responsible for the given key
// This is pretty much the only part the Chord paper exposes in their API,
// but we add some other functions for a nicer user application API.
func (entryPoint *EntryPoint) Lookup(key string, nonce string) string {
	h := sha256.New()
	h.Write([]byte(key + nonce))
	z := big.NewInt(0)
	z.SetBytes(h.Sum(nil))

	serverIdx := util.BinarySearch(entryPoint.ipHashes, z) % len(entryPoint.ipHashes)
	serverIpHash := entryPoint.ipHashes[serverIdx]

	serverIp := entryPoint.servers[serverIpHash.Text(16)]

	return serverIp
}

type Predecessors struct {
	Predecessors []string
}

type Successors struct {
	Successors []string
}

type NewNodeReq struct {
	// Which node the reuest "started" from, since it's passed down the ring
	Origin string
	// List of ndoes that have viewed this message
	ViewList []string
	// The ip of the new node
	NewNode string
}

func (entryPoint *EntryPoint) addServer(ipAddress string) {
	h := sha256.New()
	h.Write([]byte(ipAddress))
	ipHash := h.Sum(nil)
	z := big.NewInt(0)
	z.SetBytes(ipHash)

	// Add server to internal tables
	insertionPoint := util.BinarySearch(entryPoint.ipHashes, z)

	nodeInList := false
	for _, x := range entryPoint.ipHashes {
		if entryPoint.servers[x.Text(16)] == ipAddress {
			nodeInList = true
			break
		}
	}

	// ... actually, only if the entry doesn't already exists, though
	if len(entryPoint.ipHashes) == 0 || !nodeInList {
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
	} else {
		fmt.Printf("Readded node at %s with hash %x\n", ipAddress, ipHash)
	}

	// Debug printing
	fmt.Println("Current node list:")
	for _, item := range entryPoint.ipHashes {
		fmt.Printf("\t%s->%s\n", item.Text(16), entryPoint.servers[item.Text(16)])
	}

	numServers := len(entryPoint.ipHashes)

	// Node join
	// Try to find the predecessor to the new node, and use it to inform the new node.
	predIndex := insertionPoint

	for predIndex != (insertionPoint+1+numServers)%numServers {
		predIndex = (predIndex - 1 + numServers) % numServers
		fmt.Printf("trying to contact %s\n", entryPoint.servers[entryPoint.ipHashes[predIndex].Text(16)])
		predIp := entryPoint.servers[entryPoint.ipHashes[predIndex].Text(16)]

		_, err := entryPoint.requester.SendRequest(predIp, "/api/successors/nil/0", http.MethodPatch, nil, util.REQUEST_TIMEOUT)

		if err != nil {
			fmt.Println(err)
			continue
		}

		resp, err := entryPoint.requester.SendRequest(
			predIp,
			"/api/join",
			http.MethodPost,
			NewNodeReq{
				predIp,
				[]string{},
				ipAddress,
			},
			util.REQUEST_TIMEOUT)
		if err != nil || resp.StatusCode != http.StatusOK {
			fmt.Println("Error when joining node!")
		} else {
			return
		}
	}
}
