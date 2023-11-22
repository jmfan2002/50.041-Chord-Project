package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) GetValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["ValueHash"]

	queryParams := r.URL.Query()
	src := queryParams.Get("src")
	fmt.Println("Get value called")

	if src == h.NodeInfo.NodeUrl {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if src == "" {
		src = h.NodeInfo.NodeUrl
	}

	val, inMap := h.NodeInfo.NodeContents[string(key)]
	if !inMap {
		if h.NodeInfo.SuccessorArray[0] == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		fmt.Println("uh oh, key not found. Linear probing instead. Target:")
		res, _ := http.Get(
			h.NodeInfo.SuccessorArray[0] + "/api/" + key + "?src=" + src,
		)

		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading request body")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		val = string(bodyBytes)
	}

	fmt.Println("val:")
	fmt.Println(val)
	response := []byte(val)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
