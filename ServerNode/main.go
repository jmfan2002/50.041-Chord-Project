package main

import (
	"ServerNode/api"
	"ServerNode/structs"

	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func setValue(Val string) {

}

func updateFingers() {

}

func main() {
	port := 4000

	var thisNode structs.NodeInformation
	thisNode.NodeValue = "{no-value}"

	// create a new router
	router := mux.NewRouter().StrictSlash(true)

	// expose endpoints
	router.HandleFunc("/health", api.HealthCheck).Methods("GET")
	router.HandleFunc("/api", api.HealthCheck).Methods("GET", "POST", "PATCH")
	router.HandleFunc("/api/cycleHealth", api.HealthCheck).Methods("GET")

	// Catch all undefined endpoints
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Route undefined"))
	})

	// start service
	fmt.Printf("Listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
