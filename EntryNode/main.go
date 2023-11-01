package main

import (
	"EntryNode/entrypoint"
	"flag"

	"encoding/json"
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

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Health check called")

	sampleStruct := SampleStruct{
		Val: "success!",
	}

	response, err := json.Marshal(sampleStruct)
	if err != nil {
		fmt.Println("error marshalling data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Index route called")
	http.ServeFile(w, r, "client/index.html")
}

func getData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get data called")

	sampleStruct := SampleStruct{
		Data: "test123",
	}

	response, err := json.Marshal(sampleStruct)
	if err != nil {
		fmt.Println("error marshalling data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func addData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Add data called")

	// Get data

	sampleStruct := SampleStruct{
		Message: "Data added",
	}

	response, err := json.Marshal(sampleStruct)
	if err != nil {
		fmt.Println("error marshalling data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
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
	router.HandleFunc("/health", healthCheck).Methods("GET")

	router.HandleFunc("/data", entryServer.GetValue).Methods("GET")
	router.HandleFunc("/data", entryServer.SetValue).Methods("POST")

	router.HandleFunc("/join", entryServer.AddNode).Methods("POST")

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
