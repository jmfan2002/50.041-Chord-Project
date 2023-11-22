package api

import (
	"ServerNode/constants"
	"ServerNode/util"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) UpdateSuccessors(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[Debug] CycleHealthCheck called\n")

	// Parse variables from url -------------------------------------------
	StartingNodeHash := mux.Vars(r)["StartingNodeHash"]
	if StartingNodeHash == "nil" {
		StartingNodeHash = h.NodeInfo.NodeHash
	}

	FinishedLoop, err := strconv.ParseBool(mux.Vars(r)["FinishedLoop"])
	if err != nil {
		fmt.Printf("[Error] FinishedLoop %s must be a boolean\n", mux.Vars(r)["FinishedLoop"])
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Printf("[Debug] given StartingNodeHash: %s, FinishedLoop: %t\n", StartingNodeHash, FinishedLoop)
	fmt.Printf("[Debug] current node hash: %s\n", h.NodeInfo.NodeHash)

	// We've cycled back, return -------------------------------------------
	if FinishedLoop && h.NodeInfo.NodeHash >= StartingNodeHash {
		fmt.Printf("[Debug] cycle complete, returning!\n")
		util.WriteResponse(w, CycleHealthResponse{CycleSize: 1}, http.StatusOK)
		return

		// Continue the cycle -------------------------------------------
	} else {
		fmt.Printf("[Debug] continuing on the loop\n")
		for i := 0; i < min(h.NodeInfo.StoredNbrs, len(h.NodeInfo.SuccessorArray)); i++ {
			fmt.Printf("[Debug] sending msg to %s\n", h.NodeInfo.SuccessorArray[i])

			// Case: we've looped
			if util.Sha256String(h.NodeInfo.SuccessorArray[i]) <= StartingNodeHash {
				fmt.Printf("[Debug] setting finishedLoop to true, next node hash of %s <= starting hash of %s\n", util.Sha256String(h.NodeInfo.SuccessorArray[i]), StartingNodeHash)
				FinishedLoop = true
			}

			// Check the next descendant
			requestEndpoint := fmt.Sprintf("/api/cycleHealth/%s/%t", StartingNodeHash, FinishedLoop)
			resp, err := h.Requester.SendRequest(h.NodeInfo.SuccessorArray[i], requestEndpoint, http.MethodGet, constants.REQUEST_TIMEOUT)

			if err != nil {
				// Descendent is unresponsive
				fmt.Printf("[Debug] child %s is not healthy, trying next\n", h.NodeInfo.SuccessorArray[i])

			} else if resp.StatusCode != http.StatusOK {
				// Descendent returns a bad status code, return
				fmt.Printf("[Debug] next node is reporting break in cycle with status code: %d\n", resp.StatusCode)
				w.WriteHeader(resp.StatusCode)
				return

			} else {
				// Descendent responds ok, pass response forward
				var healthResp = &CycleHealthResponse{}
				err := util.ReadResponseBody(resp, healthResp)

				if err != nil {
					fmt.Printf("[Error] failed to decode response body: %s\n", err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				healthResp.CycleSize++
				fmt.Printf("[Debug] received healthCheck response from child! nodes: %d\n", healthResp.CycleSize)
				util.WriteResponse(w, healthResp, http.StatusOK)
				return
			}
		}
	}

	fmt.Printf("[Error] all descendants have failed for current NodeUrl |%s|\n", h.NodeInfo.NodeUrl)
	w.WriteHeader(http.StatusInternalServerError)
}
