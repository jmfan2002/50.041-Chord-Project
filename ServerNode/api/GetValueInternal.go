package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) GetValueInternal(w http.ResponseWriter, r *http.Request) {
	// Parse variables from url -------------------------------------------
	PreviousNodeHash := mux.Vars(r)["PreviousNodeHash"]
	Key := mux.Vars(r)["Key"]
	Nonce := mux.Vars(r)["Nonce"]
	// fmt.Printf("[Debug] GetValueInternal called on key %s nonce %s PreviousNodeHash %s\n", Key, Nonce, PreviousNodeHash)

	h.GetValueHelper(w, Key, Nonce, PreviousNodeHash)
}
