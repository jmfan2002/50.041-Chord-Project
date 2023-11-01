package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func CycleHealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("given nodeHash: %s", mux.Vars(r)["startingNodeHash"])
}
