package entrypoint

import (
	"EntryNode/util"
	"bytes"
	"crypto/sha256"
	"encoding/json"
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
}

func New(k int) *EntryPoint {
	return &EntryPoint{
		ipHashes:        make([]big.Int, 0),
		servers:         make(map[string]string),
		toleratedFaults: k,
	}
}

type ChordSetValueReq struct {
	Key   string
	Value string
	Nonce string
}

func (entryPoint *EntryPoint) setKVP(key string, val string) {
	for i := 0; i < entryPoint.toleratedFaults; i += 1 {
		serverAddress := entryPoint.Lookup(key, strconv.Itoa(i))
		fmt.Printf("Sending %s -> %s to %s\n", key, val, serverAddress)

		j, err := json.Marshal(ChordSetValueReq{
			key,
			val,
			strconv.Itoa(i),
		})
		if err != nil {
			fmt.Println("Error creating request body")
			continue
		}

		res, err := http.Post(serverAddress+"/api",
			"application/json",
			bytes.NewBuffer(j),
		)
		if err != nil {
			fmt.Printf("An error occurred %s\n", err.Error())
		}
		fmt.Printf("Request result: %s\n", res.Status)
	}
}

type ChordGetValueRes struct {
	NodeContents map[string]string
}

func (entryPoint *EntryPoint) getKVP(key string) string {
	for i := 0; i < entryPoint.toleratedFaults; i += 1 {
		// Find the server that should have the key
		serverAddress := entryPoint.Lookup(key, strconv.Itoa(i))
		fmt.Printf("Getting %s from %s\n", key, serverAddress)

		// Send request to server
		resp, err := http.Get(
			fmt.Sprintf("%s/api/%s/%s", serverAddress, key, strconv.Itoa(i)))
		if err != nil {
			fmt.Printf("An error occurred %s\n", err.Error())
		}

		fmt.Printf("Request result: %s\n", resp.Status)

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading server node response")
			continue
		}
		return string(bodyBytes)
	}

	return ""
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

func (entryPoint *EntryPoint) addServer(ipAddress string) {
	h := sha256.New()
	h.Write([]byte(ipAddress))
	ipHash := h.Sum(nil)
	z := big.NewInt(0)
	z.SetBytes(ipHash)

	// Add server to internal tables
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

	// Debug printing
	fmt.Printf("Added node at %s with hash %x\n", ipAddress, ipHash)
	fmt.Println("Current node list:")
	for _, item := range entryPoint.ipHashes {
		fmt.Printf("\t%s->%s\n", item.Text(16), entryPoint.servers[item.Text(16)])
	}

	numServers := len(entryPoint.ipHashes)

	// If first server, set successors list to self
	if numServers == 1 {
		data, _ := json.Marshal(Successors{
			Successors: []string{ipAddress},
		})
		_, _ = http.Post(
			ipAddress+"/api/successors",
			"application/json",
			bytes.NewBuffer(data),
		)

		return
	}

	// Node join
	// Try to find the predecessor to the new node, and use it to inform the new node.
	predIndex := insertionPoint

	for predIndex != (insertionPoint+1+numServers)%numServers {
		predIndex = (predIndex - 1 + numServers) % numServers
		fmt.Printf("trying to contact at %d\n", predIndex)

		// 1. get the successors (GetSuccessors) of the node before where it belongs
		fmt.Println("Step 1")
		var succOfPred Successors
		predIp := entryPoint.servers[entryPoint.ipHashes[predIndex].Text(16)]
		// get successors of the preceding node
		resp, err := http.Get(predIp + "/api/successors")
		if err != nil {
			fmt.Println(err)
			continue
		}
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading request body")
			return
		}
		err = json.Unmarshal(bodyBytes, &succOfPred)

		// 2. set the successors of the new node to this list (SetSuccessors)
		fmt.Println("Step 2")
		data, err := json.Marshal(succOfPred)
		resp, err = http.Post(
			ipAddress+"/api/successors",
			"application/json",
			bytes.NewBuffer(data),
		)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// 3. set the successors of the old node to [newNode, oldSuccessors...]
		fmt.Println("Step 3")
		succOfPred.Successors = append([]string{ipAddress}, succOfPred.Successors...)
		data, err = json.Marshal(succOfPred)
		resp, err = http.Post(
			predIp+"/api/successors",
			"application/json",
			bytes.NewBuffer(data),
		)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// 4. call UpdateSuccessors on any node in the system (it will be safest to call it on the node directly before, but it doesn't really matter)
		fmt.Println("Step 4")
		_, err = http.NewRequest(
			http.MethodPatch,
			predIp+"/api/successors/nil/0",
			bytes.NewBuffer(data),
		)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// 5. call ReassignEntries on the node directly preceding the new node. this will move the necessary values to the new node
		fmt.Println("Step 5")
		_, err = http.NewRequest(
			http.MethodPatch,
			predIp+"/api/entries",
			bytes.NewBuffer([]byte{}),
		)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	// TODO: delete the stuff below

	// Get predecessor and successor node info, if they exist

	/*
		if numServers == 1 {
			return
		}

		succIp := entryPoint.servers[entryPoint.ipHashes[(insertionPoint+1)%numServers].Text(16)]

		var predOfSucc Predecessors
		var succOfPred Successors

		// get predecessors of the successor node
		resp, _ := http.Get(succIp + "/api/predecessors")
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading request body")
			return
		}

		err = json.Unmarshal(bodyBytes, &predOfSucc)
		if err != nil {
			fmt.Println("Error parsing request body")
			return
		}

		// get predecessors of the successor node
		resp, _ = http.Get(predIp + "/api/successors")
		bodyBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading request body")
			return
		}

		err = json.Unmarshal(bodyBytes, &succOfPred)
		if err != nil {
			fmt.Println("Error parsing request body")
			return
		}

		// Update new node
		data, err := json.Marshal(predOfSucc)
		if err != nil {
			fmt.Printf("An error occurred %s\n", err.Error())
		}
		patchRes, err := http.NewRequest(
			http.MethodPatch, ipAddress+"/api/SetPredecessors",
			bytes.NewBuffer(data),
		)
		if err != nil {
			fmt.Printf("An error occurred %s\n", err.Error())
		}
		fmt.Printf(patchRes.Response.Status)

		data, err = json.Marshal(succOfPred)
		if err != nil {
			fmt.Printf("An error occurred %s\n", err.Error())
		}
		patchRes, err = http.NewRequest(
			http.MethodPatch, ipAddress+"/api/SetSuccessors",
			bytes.NewBuffer(data),
		)
		if err != nil {
			fmt.Printf("An error occurred %s\n", err.Error())
		}
		fmt.Printf(patchRes.Response.Status)
	*/
}
