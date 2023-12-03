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

func (entryPoint *EntryPoint) trySetKVP(key string, nonce string, val string) {
	serverAddress := entryPoint.Lookup(key, nonce)
	fmt.Printf("Sending %s -> %s to %s\n", key, val, serverAddress)

	j, err := json.Marshal(ChordSetValueReq{
		key,
		val,
		nonce,
	})
	if err != nil {
		fmt.Println("Error creating request body")
		return
	}

	res, err := http.Post(serverAddress+"/api",
		"application/json",
		bytes.NewBuffer(j),
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
	resp, err := http.Get(
		fmt.Sprintf("%s/api/%s/%s", serverAddress, key, nonce))
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
	WiewList []string
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

	// ... actually, only if the entry doesn't already exists, though
	if len(entryPoint.ipHashes) == 0 || entryPoint.ipHashes[(insertionPoint-1)%len(entryPoint.ipHashes)].Cmp(z) != 0 {
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
		fmt.Printf("trying to contact %s\n", entryPoint.servers[entryPoint.ipHashes[predIndex].Text(16)])
		predIp := entryPoint.servers[entryPoint.ipHashes[predIndex].Text(16)]

		data, _ := json.Marshal(NewNodeReq{
			predIp,
			[]string{},
			ipAddress,
		})
		resp, err := http.Post(
			predIp+"/api/join",
			"application/json",
			bytes.NewBuffer(data),
		)

		if err != nil || resp.StatusCode != http.StatusOK {
			fmt.Println("Error when joining node!")
		} else {
			return
		}

		/*
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
		*/
	}
}
