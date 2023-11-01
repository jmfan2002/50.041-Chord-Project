package main

import (
	"ServerNode/api"
	"ServerNode/structs"

	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	port := 4000

	var thisNode structs.NodeInformation
	thisNode.NodeValue = "{no-value}"

	// create a new router
	router := mux.NewRouter().StrictSlash(true)

	// expose endpoints
	router.HandleFunc("/api/health", api.HealthCheck).Methods("GET")
	router.HandleFunc("/api/cycleHealth/{StartingNodeHash}", api.CycleHealthCheck).Methods("GET")

	// Catch all undefined endpoints
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Route undefined"))
	})

	// start service
	fmt.Printf("Listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
