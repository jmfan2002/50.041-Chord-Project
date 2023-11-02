package main

import (
	"ServerNode/api"
	"bytes"
	"encoding/json"
	"flag"
	"time"

	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type JoinReq struct {
	NewNodeAddress string
}

func main() {
	// Command line flags
	portPtr := flag.Int("port", 4000, "The port to serve the entrypoint on")
	entryPtr := flag.String("entry", "127.0.0.1:3000", "The ip of the entry point")
	flag.Parse()
	entry := *entryPtr
	port := *portPtr

	STORED_NBRS := 4

	// create a new router
	router := mux.NewRouter().StrictSlash(true)
	// create a handler that stores and updates our node information
	handler := api.NewHandler(fmt.Sprintf("http://localhost:%d", port), STORED_NBRS)
	fmt.Printf("[Debug] set up node %s\n", handler.NodeInfo)

	// Set/Get values
	router.HandleFunc("/api/{ValueHash}", handler.GetValue).Methods("GET")
	router.HandleFunc("/api", handler.SetValue).Methods("POST")

	// Update successors
	router.HandleFunc("/api/succ", handler.SetSucc).Methods("POST")

	// expose endpoints
	router.HandleFunc("/api/health", handler.HealthCheck).Methods("GET")
	router.HandleFunc("/api/cycleHealth/{StartingNodeHash}/{FinishedLoop}", handler.CycleHealthCheck).Methods("GET")

	// Catch all undefined endpoints
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Route undefined"))
	})

	j, _ := json.Marshal(JoinReq{
		fmt.Sprintf("http://localhost:%d", port),
	})

	// Notify entrypoint
	go func() {
		time.Sleep(time.Second * 1)
		fmt.Println("Waking up, notifying entry point")
		http.Post("http://"+entry+"/join", "application/json",
			bytes.NewBuffer(j))
	}()

	// start service
	fmt.Printf("Listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
