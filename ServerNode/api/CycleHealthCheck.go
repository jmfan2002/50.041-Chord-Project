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
	StartingNodeHash := mux.Vars(r)["StartingNodeHash"]
	FinishedLoop, err := strconv.ParseBool(mux.Vars(r)["FinishedLoop"])
	if err != nil {
		fmt.Printf("[Error] FinishedLoop %s must be a boolean\n", mux.Vars(r)["FinishedLoop"])
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Case: we return an equal or greater node after looping
	if FinishedLoop && h.NodeInfo.NodeHash >= StartingNodeHash {
		fmt.Printf("[Debug] Finished the loop\n")
		w.WriteHeader(http.StatusOK)
		return

	} else {
		for i := 0; i < h.NodeInfo.StoredNbrs; i++ {
			// Case: we've looped
			if util.Sha256String(h.NodeInfo.SuccessorArray[i]) <= StartingNodeHash {
				FinishedLoop = true
			}

			// Case: check the next descendant with timeout
			requestUrl := fmt.Sprintf("%s/%s/%t", h.NodeInfo.SuccessorArray[i], StartingNodeHash, FinishedLoop)
			fmt.Printf("[Debug] sending request to nbr %d: %s\n", i, requestUrl)
			resp, err := util.SendRequest(requestUrl, constants.REQUEST_TIMEOUT)
			if err == nil {
				// Case: the node responds
				w.WriteHeader(resp.StatusCode)
				break
			}

			// Otherwise: try the next descendent
		}
	}
	fmt.Printf("given StartingNodeHash: %s, FinishedLoop: %t\n", StartingNodeHash, FinishedLoop)
}
