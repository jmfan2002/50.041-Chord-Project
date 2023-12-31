package main

import (
	"EntryNode/entrypoint"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// Note: values MUST start with capital letter to be exported
type SampleStruct struct {
	Val     string `json:"val"`
	Data    string `json:"data"`
	Message string `json:"message"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Index route called")
	http.ServeFile(w, r, "client/index.html")
}

func periodicWrite(entry entrypoint.EntryPoint) {
	for {
		time.Sleep(5 * time.Second)
		entry.WriteState()
	}
}

func main() {
	// Command line flags
	portPtr := flag.Int("port", 3000, "The port to serve the entrypoint on")
	kPtr := flag.Int("k", 1, "The number of Chord nodes to save values to. This allows for k-1 fault tolerance")
	flag.Parse()

	port := *portPtr
	numFaults := *kPtr - 1

	// create entrypoint
	var entryServer *entrypoint.EntryPoint
	if _, err := os.Stat("entrypoint\\state.txt"); err == nil {
		fmt.Println("state.txt exists, loading from file...")
		entryServer = entrypoint.ReadState()
	} else {
		fmt.Println("state.txt does not exist, creating new entrypoint...")
		entryServer = entrypoint.New(numFaults)
	}

	// create a new router
	router := mux.NewRouter().StrictSlash(true)

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./client/static"))))

	// expose endpoints
	router.HandleFunc("/health", entryServer.HealthCheck).Methods("GET")
	router.HandleFunc("/cycleHealth", entryServer.CycleHealth).Methods("GET")

	router.HandleFunc("/data", entryServer.GetValue).Methods("GET")
	router.HandleFunc("/data", entryServer.SetValue).Methods("POST")
	router.HandleFunc("/data/hashTable", entryServer.GetHashTable).Methods("GET")

	router.HandleFunc("/join", entryServer.AddNode).Methods("POST")

	router.HandleFunc("/nodes", entryServer.GetNodes).Methods("GET")

	// Serve webpage
	router.HandleFunc("/", indexHandler)

	// Catch all undefined endpoints
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Route undefined"))
	})

	go periodicWrite(*entryServer)

	// start service
	fmt.Printf("Listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
