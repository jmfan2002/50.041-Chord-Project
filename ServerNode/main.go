package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Note: values MUST start with capital letter to be exported
type SampleStruct struct {
	Val string `json:"val"`
}

type NodeInformation struct {
	NodeValue string
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

func getValue(w http.ResponseWriter, r *http.Request, nodeData *NodeInformation) {
	fmt.Println("Get value called")

	response, err := json.Marshal(nodeData)
	if err != nil {
		fmt.Println("error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response = []byte(nodeData.NodeValue)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func setValue(Val string) {

}

func updateFingers() {

}

func main() {
	port := 4000

	var thisNode NodeInformation
	thisNode.NodeValue = "{no-value}"

	// create a new router
	router := mux.NewRouter().StrictSlash(true)

	// expose endpoints
	router.HandleFunc("/health", healthCheck).Methods("GET")
	router.HandleFunc("/api", healthCheck).Methods("GET", "POST", "PATCH")
	router.HandleFunc("/api/cycleHealth", healthCheck).Methods("GET")

	// Catch all undefined endpoints
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Route undefined"))
	})

	// start service
	fmt.Printf("Listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
