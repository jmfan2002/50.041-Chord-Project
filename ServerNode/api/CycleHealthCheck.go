package api

import (
	"ServerNode/constants"
	"ServerNode/util"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type CycleHealthResponse struct {
	CycleSize int
}

func (h *Handler) CycleHealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[Debug] CycleHealthCheck called\n")

	// Parse variables from url -------------------------------------------
	PreviousNodeHash := mux.Vars(r)["PreviousNodeHash"]
	StartNodeHash := mux.Vars(r)["StartNodeHash"]
	fmt.Printf("[Debug] given PreviousNodeHash: %s, StartNodeHash: %s\n", PreviousNodeHash, StartNodeHash)
	fmt.Printf("[Debug] current node hash: %s\n", h.NodeInfo.NodeHash)

	// We've cycled back, return -------------------------------------------
	if h.NodeInfo.NodeHash == StartNodeHash {
		fmt.Printf("[Debug] cycle complete, we've reached the start node! \n")
		util.WriteResponse(w, CycleHealthResponse{CycleSize: 0}, http.StatusOK)
		return
	}
	if (StartNodeHash == "nil") {
		StartNodeHash = h.NodeInfo.NodeHash;
	}

	// Continue the cycle -------------------------------------------
	fmt.Printf("[Debug] continuing on the loop\n")
	for i := 0; i < min(h.NodeInfo.StoredNbrs, len(h.NodeInfo.SuccessorArray)); i++ {
		fmt.Printf("[Debug] sending msg to %s\n", h.NodeInfo.SuccessorArray[i])

		// Check the next descendant
		requestEndpoint := fmt.Sprintf("/api/cycleHealth/%s/%s", h.NodeInfo.NodeHash, StartNodeHash)
		resp, err := h.Requester.SendRequest(h.NodeInfo.SuccessorArray[i], requestEndpoint, http.MethodGet, nil, constants.REQUEST_TIMEOUT)

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
			err := util.ReadBody(resp.Body, healthResp)
			// time.Sleep(1 * time.Second)

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

	fmt.Printf("[Error] all descendants have failed for current NodeUrl |%s|\n", h.NodeInfo.NodeUrl)
	w.WriteHeader(http.StatusInternalServerError)
}
