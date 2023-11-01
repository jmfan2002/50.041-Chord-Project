package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
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

func (entryPoint *EntryPoint) GetValue(key string) string {
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
	for _, item := range entryPoint.ipHashes {
		fmt.Printf("\t%s->%s\n", item.Text(16), entryPoint.servers[item.Text(16)])
	}
}

// Finds the insertion index of item in arr
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

// Request handler stuff

func readRequestBody(w http.ResponseWriter, r *http.Request, reqBody any) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(bodyBytes, reqBody)
	if err != nil {
		fmt.Println("Error parsing request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func writeSuccessResponse(w http.ResponseWriter, body any) {
	response, err := json.Marshal(body)
	if err != nil {
		fmt.Println("error marshalling data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

type AddDataBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (entryPoint *EntryPoint) AddData(w http.ResponseWriter, r *http.Request) {
	// Get key and value from request
	var reqBody AddDataBody

	readRequestBody(w, r, &reqBody)

	entryPoint.SetValue(reqBody.Key, reqBody.Value)

	sampleStruct := SampleStruct{
		Data: "test123",
	}

	writeSuccessResponse(w, &sampleStruct)
}

func (entryPoint *EntryPoint) GetData(w http.ResponseWriter, r *http.Request) {
	// Get key from request
	queryParams := r.URL.Query()

	key := queryParams.Get("key")
	val := entryPoint.GetValue(key)

	sampleStruct := SampleStruct{
		Val:  val,
		Data: "test123",
	}

	writeSuccessResponse(w, &sampleStruct)
}

type JoinReqBody struct {
	NewNodeAddress string `json:"NewNodeAddress"`
}

func (entryPoint *EntryPoint) JoinReq(w http.ResponseWriter, r *http.Request) {
	// Get key and value from request
	var reqBody JoinReqBody

	readRequestBody(w, r, &reqBody)

	fmt.Println("Receiving join for %s\n", reqBody.NewNodeAddress)
	entryPoint.AddNode(reqBody.NewNodeAddress)

	sampleStruct := SampleStruct{
		Data: "test123",
	}

	writeSuccessResponse(w, &sampleStruct)
}
