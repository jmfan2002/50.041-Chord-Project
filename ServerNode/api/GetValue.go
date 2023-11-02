package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) GetValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	key := vars["ValueHash"]
	fmt.Println("Get value called")
	fmt.Println(key)

	val := h.NodeInfo.NodeContents[string(key)]
	fmt.Println(val)

	response := []byte(val)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
