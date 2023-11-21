package api

import (
	"ServerNode/constants"
	"ServerNode/util"
	"fmt"
	"net/http"
	"time"

	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) CycleHealthCheck(w http.ResponseWriter, r *http.Request) {
	// Parse variables from url
	StartingNodeHash := mux.Vars(r)["StartingNodeHash"]
	FinishedLoop, err := strconv.ParseBool(mux.Vars(r)["FinishedLoop"])
	if err != nil {
		fmt.Printf("[Error] FinishedLoop %s must be a boolean\n", mux.Vars(r)["FinishedLoop"])
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Printf("[Debug] given StartingNodeHash: %s, FinishedLoop: %t\n", StartingNodeHash, FinishedLoop)
	fmt.Printf("[Debug] current node hash: %s\n", h.NodeInfo.NodeHash)

	// Case: we reach an equal or greater node from the start after having looped
	if FinishedLoop && h.NodeInfo.NodeHash >= StartingNodeHash {
		fmt.Printf("[Debug] Finished the loop\n")
		w.WriteHeader(http.StatusOK)
		return

	} else {
		fmt.Printf("[Debug] continuing on the loop\n")
		for i := 0; i < min(h.NodeInfo.StoredNbrs, len(h.NodeInfo.SuccessorArray)); i++ {
			fmt.Printf("[Debug] sending msg to %s\n",h.NodeInfo.SuccessorArray[i])
			// Case: we've looped
			if util.Sha256String(h.NodeInfo.SuccessorArray[i]) <= StartingNodeHash {
				FinishedLoop = true
			}

			// Case: check the next descendant with timeout
			requestEndpoint := fmt.Sprintf("/%s/%t", StartingNodeHash, FinishedLoop)
			resp, err := h.Requester.SendRequest(h.NodeInfo.SuccessorArray[i], requestEndpoint, http.MethodGet, constants.REQUEST_TIMEOUT)
			time.Sleep(5 * time.Second)	// DEBUG test heartbeat
			if err == nil {
				fmt.Printf("[Debug], received healthCheck response from child!")
				// Case: the node responds
				w.WriteHeader(resp.StatusCode)
				return
			}

			// Otherwise: try the next descendent
		}
	}

	fmt.Printf("[Error] all descendants have failed\n")
	w.WriteHeader(http.StatusInternalServerError)
}
