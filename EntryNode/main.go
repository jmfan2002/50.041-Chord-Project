package main

import (
	"EntryNode/entrypoint"
	"flag"
	"fmt"
	"log"
	"net/http"

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

func main() {
	// Command line flags
	portPtr := flag.Int("port", 3000, "The port to serve the entrypoint on")
	flag.Parse()

	port := *portPtr

	// create entrypoint
	entryServer := entrypoint.New()

	// create a new router
	router := mux.NewRouter().StrictSlash(true)

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./client/static"))))

	// expose endpoints
	router.HandleFunc("/health", entryServer.HealthCheck).Methods("GET")

	router.HandleFunc("/data", entryServer.GetValue).Methods("GET")
	router.HandleFunc("/data", entryServer.SetValue).Methods("POST")

	router.HandleFunc("/join", entryServer.AddNode).Methods("POST")

	router.HandleFunc("/nodes", entryServer.GetNodes).Methods("GET")

	// Serve webpage
	router.HandleFunc("/", indexHandler)

	// Catch all undefined endpoints
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Route undefined"))
	})

	// start service
	fmt.Printf("Listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
