package main

import (
	"ServerNode/api"
	"ServerNode/structs"
	"ServerNode/util"
	"bytes"
	"encoding/json"
	"strconv"
	"time"

	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	// Parse arguments
	usageStr := "usage: go run main.go <port> <storedNbrs> <baseURL> <entrypoint>"
	if len(os.Args) < 5 {
		fmt.Printf("[Error] port or storedNbrs or baseURL or entrypoint missing. %s\n", usageStr)
		os.Exit(0)
	}
	port, err := strconv.Atoi(os.Args[1])
	if err != nil || port <= 1024 {
		fmt.Printf("[Error] invalid port: %s, %s\n", os.Args[1], usageStr)
		os.Exit(0)
	}

	STORED_NBRS, err := strconv.Atoi(os.Args[2])
	if err != nil || STORED_NBRS < 1 {
		fmt.Printf("[Error] invalid storedNbrs: %s, %s\n", os.Args[2], usageStr)
		os.Exit(0)
	}

	BASE_URL := os.Args[3]
	if BASE_URL == "" {
		fmt.Printf("[Error] invalid baseURL (set to localhost if testing locally): %s, %s\n", os.Args[3], usageStr)
		os.Exit(0)
	}

	ENTRYPOINT := os.Args[4]
	if ENTRYPOINT == "" {
		fmt.Printf("[Error] invalid baseURL entrypoint url (set to localhost if testing locally): %s, %s\n", os.Args[4], usageStr)
		os.Exit(0)
	}

	// create a new router
	router := mux.NewRouter().StrictSlash(true)
	// create a handler that stores and updates our node information
	handler := api.NewHandler(fmt.Sprintf("http://%s:%d", BASE_URL, port), STORED_NBRS)
	// for testing purposes, you can run nodes on localhost 2000, 3000, and 4000. Then, you can remove node 3000 and it will still be successful
	handler.NodeInfo.NodeHash = util.Sha256String(fmt.Sprintf("http://%s:%d", BASE_URL, port)) // DEBUG: REMOVE WHEN DONE

	fmt.Printf("[Debug] set up node %s\n", handler.NodeInfo)

	// expose endpoints
	router.HandleFunc("/api/health", handler.HealthCheck).Methods("GET")
	router.HandleFunc("/api/cycleHealth/{StartNodeHash}", handler.CycleHealthCheck).Methods("GET")
	router.HandleFunc("/api/successors", handler.SetSuccessors).Methods("POST")
	router.HandleFunc("/api/successors", handler.GetSuccessors).Methods("GET")
	router.HandleFunc("/api/successors/{StartingNodeHash}/{CurrentOverlap}", handler.UpdateSuccessors).Methods("PATCH")
	router.HandleFunc("/api/entries", handler.ReassignEntries).Methods("PATCH")

	router.HandleFunc("/api/hashTable", handler.GetHashTable).Methods("GET")
	router.HandleFunc("/api/{Key}/{Nonce}", handler.GetValue).Methods("GET")
	router.HandleFunc("/api", handler.SetValue).Methods("POST")

	router.HandleFunc("/api/join", handler.NewNode).Methods("POST")

	// Internal endpoints
	router.HandleFunc("/api/{Key}/{Nonce}/{PreviousNodeHash}", handler.GetValueInternal).Methods("GET")

	// Catch all undefined endpoints
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Route undefined"))
	})

	// start service
	// log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
	// ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	j, _ := json.Marshal(structs.JoinReq{
		fmt.Sprintf("http://%s:%d", BASE_URL, port),
	})
	// entry := "entry_node:3000"

	fmt.Printf("Listening on port %d\n", port)
	go func() {
		fmt.Println("[Debug]")
		time.Sleep(1 * time.Second)
		fmt.Println("Waking up, notifying entry point")
		http.Post(fmt.Sprintf("http://%s", ENTRYPOINT)+"/join", "application/json",
			bytes.NewBuffer(j))
		fmt.Printf("[Debug] set up node %s\n", handler.NodeInfo)
	}()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))

}
