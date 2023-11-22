package main

import (
	"ServerNode/api"
	"strconv"

	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	// Parse arguments
	usageStr := "usage: go run main.go <port>"
	if len(os.Args) == 1 {
		fmt.Printf("[Error] port missing. %s\n", usageStr)
		os.Exit(0)
	}
	port, err := strconv.Atoi(os.Args[1])
	if err != nil || port <= 1024 {
		fmt.Printf("[Error] invalid port: %s, %s\n", os.Args[1], usageStr)
		os.Exit(0)
	}

	STORED_NBRS := 4

	// create a new router
	router := mux.NewRouter().StrictSlash(true)
	// create a handler that stores and updates our node information
	handler := api.NewHandler(fmt.Sprintf("http://10.12.103.97:%d", port), STORED_NBRS)
	// for testing purposes, you can run nodes on localhost 2000, 3000, and 4000. Then, you can remove node 3000 and it will still be successful
	// handler.NodeInfo.SuccessorArray = append(handler.NodeInfo.SuccessorArray, fmt.Sprintf("http://10.12.103.97:%d", (port+1000) % 5000))
	// handler.NodeInfo.SuccessorArray = append(handler.NodeInfo.SuccessorArray, fmt.Sprintf("http://10.12.103.97:%d", (port+2000) % 5000))
	// handler.NodeInfo.SuccessorArray = append(handler.NodeInfo.SuccessorArray, fmt.Sprintf("http://10.12.103.97:%d", (port+3000) % 5000))


	fmt.Printf("[Debug] set up node %s\n", handler.NodeInfo)

	// expose endpoints
	router.HandleFunc("/api/health", handler.HealthCheck).Methods("GET")
	router.HandleFunc("/api/cycleHealth/{StartingNodeHash}/{FinishedLoop}", handler.CycleHealthCheck).Methods("GET")

	// Catch all undefined endpoints
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Route undefined"))
	})

	// start service
	fmt.Printf("Listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
