package api

import (
	"ServerNode/constants"
	"ServerNode/util"
	"fmt"
	"net/http"

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
	fmt.Printf("given StartingNodeHash: %s, FinishedLoop: %t\n", StartingNodeHash, FinishedLoop)
	fmt.Printf("current node hash: %s\n", h.NodeInfo.NodeHash)

	// Case: we reach an equal or greater node from the start after having looped
	if FinishedLoop && h.NodeInfo.NodeHash >= StartingNodeHash {
		fmt.Printf("[Debug] Finished the loop\n")
		w.WriteHeader(http.StatusOK)
		return

	} else {
		fmt.Printf("[Debug] continuing on the loop\n")
		for i := 0; i < min(h.NodeInfo.StoredNbrs, len(h.NodeInfo.SuccessorArray)); i++ {
			// Case: we've looped
			if util.Sha256String(h.NodeInfo.SuccessorArray[i]) <= StartingNodeHash {
				FinishedLoop = true
			}

			// Case: check the next descendant with timeout
			requestUrl := fmt.Sprintf("%s/%s/%t", h.NodeInfo.SuccessorArray[i], StartingNodeHash, FinishedLoop)
			fmt.Printf("[Debug] sending request to nbr %d: %s\n", i, requestUrl)
			resp, err := h.Requester.SendRequest(requestUrl, constants.REQUEST_TIMEOUT)
			if err == nil {
				// Case: the node responds
				w.WriteHeader(resp.StatusCode)
				break
			}

			// Otherwise: try the next descendent
		}
	}

	w.WriteHeader(http.StatusInternalServerError)
}
