package entrypoint

import (
	"EntryNode/util"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strconv"
)

type EntryPoint struct {
	IpHashes []big.Int

	// Maps hashes as strings, to ips as strings
	Servers map[string]string

	// number of faults to tolerate
	ToleratedFaults int

	// used to send heartbeat requests
	requester util.HeartbeatRequester
}

func New(k int) *EntryPoint {
	return &EntryPoint{
		IpHashes:        make([]big.Int, 0),
		Servers:         make(map[string]string),
		ToleratedFaults: k,
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

	_, err := entryPoint.requester.SendRequest(
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
		for i := 0; i < len(entryPoint.IpHashes); i += 1 {
			if entryPoint.IpHashes[i].Cmp(z) == 0 {
				entryPoint.IpHashes = append(entryPoint.IpHashes[:i], entryPoint.IpHashes[i+1:]...)
				break
			}
		}
		return
	}
	// fmt.Printf("Request result: %s\n", res.Status)
}

func (entryPoint *EntryPoint) setKVP(key string, val string) {
	for i := 0; i <= entryPoint.ToleratedFaults; i += 1 {
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
		for i := 0; i < len(entryPoint.IpHashes); i += 1 {
			if entryPoint.IpHashes[i].Cmp(z) == 0 {
				entryPoint.IpHashes = append(entryPoint.IpHashes[:i], entryPoint.IpHashes[i+1:]...)
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

	for i := 0; i <= entryPoint.ToleratedFaults; i += 1 {
		// Find the server that should have the key
		serverAddress := entryPoint.Lookup(key, strconv.Itoa(i))
		fmt.Printf("Getting %s from %s\n", key, serverAddress)

		go entryPoint.tryGetKVP(serverAddress, key, strconv.Itoa(i), in)
	}

	ostream := make(chan string)

	go func() {
		errors := make([]string, 0)
		out := ""
		flag := false

		for i := 0; i < entryPoint.ToleratedFaults; i += 1 {
			next := <-in
			if !next.isError && next.data != "" {
				ostream <- next.data
				flag = true
			} else {
				errors = append(errors, next.data)
			}
		}

		if flag {
			ostream <- ""
		}

		for _, nonce := range errors {
			entryPoint.trySetKVP(key, nonce, out)
		}
	}()

	return <-ostream
}

// Returns the IP address of the node responsible for the given key
// This is pretty much the only part the Chord paper exposes in their API,
// but we add some other functions for a nicer user application API.
func (entryPoint *EntryPoint) Lookup(key string, nonce string) string {
	h := sha256.New()
	h.Write([]byte(key + nonce))
	z := big.NewInt(0)
	z.SetBytes(h.Sum(nil))

	serverIdx := util.BinarySearch(entryPoint.IpHashes, z) % len(entryPoint.IpHashes)
	serverIpHash := entryPoint.IpHashes[serverIdx]

	serverIp := entryPoint.Servers[serverIpHash.Text(16)]

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
	insertionPoint := util.BinarySearch(entryPoint.IpHashes, z)

	nodeInList := false
	for _, x := range entryPoint.IpHashes {
		if entryPoint.Servers[x.Text(16)] == ipAddress {
			nodeInList = true
			break
		}
	}

	// ... actually, only if the entry doesn't already exists, though
	if len(entryPoint.IpHashes) == 0 || !nodeInList {
		if len(entryPoint.IpHashes) == insertionPoint {
			entryPoint.IpHashes = append(entryPoint.IpHashes, *z)
		} else {
			entryPoint.IpHashes = append(
				entryPoint.IpHashes[:insertionPoint+1],
				entryPoint.IpHashes[insertionPoint:]...,
			)
			entryPoint.IpHashes[insertionPoint] = *z
		}

		entryPoint.Servers[z.Text(16)] = ipAddress
		fmt.Printf("Added node at %s with hash %x\n", ipAddress, ipHash)
	} else {
		fmt.Printf("Readded node at %s with hash %x\n", ipAddress, ipHash)
	}

	// Debug printing
	fmt.Println("Current node list:")
	for _, item := range entryPoint.IpHashes {
		fmt.Printf("\t%s->%s\n", item.Text(16), entryPoint.Servers[item.Text(16)])
	}

	numServers := len(entryPoint.IpHashes)

	// Node join
	// Try to find the predecessor to the new node, and use it to inform the new node.
	predIndex := insertionPoint

	for predIndex != (insertionPoint+1+numServers)%numServers {
		fmt.Println("[DEBUG] looping")
		predIndex = (predIndex - 1 + numServers) % numServers
		// fmt.Printf("trying to contact %s\n", entryPoint.Servers[entryPoint.IpHashes[predIndex].Text(16)])
		predIp := entryPoint.Servers[entryPoint.IpHashes[predIndex].Text(16)]

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

func (EntryPoint *EntryPoint) WriteState() {
	file, err := os.Create("entrypoint\\state.txt")
	if err != nil {
		fmt.Println("Error creating state.txt")
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(EntryPoint); err != nil {
		fmt.Printf("Error encoding EntryPoint to state.txt file. %v\n", err)
		return
	}

	fmt.Println("Entrypoint successfully backed up to state.txt")
}

func ReadState() *EntryPoint {
	var data EntryPoint

	file, err := os.Open("entrypoint\\state.txt")
	if err != nil {
		fmt.Println("Error opening state.txt file to read from, blank EntryPoint returned.")
		return &EntryPoint{}
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		fmt.Println("Error decoding state.txt, blank EntryPoint returned.")
		return &EntryPoint{}
	}

	// Unfortunately, go's bigInt is garbage, so we need to regenerate the array on our own
	data.IpHashes = make([]big.Int, 0)
	for key := range data.Servers {
		z := big.NewInt(0)
		z.SetString(key, 16)

		insertion := util.BinarySearch(data.IpHashes, z)

		data.IpHashes = append(
			data.IpHashes[:insertion],
			append(
				[]big.Int{*z},
				data.IpHashes[insertion:]...,
			)...)
	}

	return &data
}
