package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) GetValue(w http.ResponseWriter, r *http.Request) {
	Key := mux.Vars(r)["Key"]
	Nonce := mux.Vars(r)["Nonce"]
	fmt.Printf("[Debug] GetValue called on key %s nonce %s\n", Key, Nonce)

	h.GetValueHelper(w, Key, Nonce, "nil")
}
