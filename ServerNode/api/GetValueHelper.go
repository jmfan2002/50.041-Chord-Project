package api

import (
	"ServerNode/constants"
	"ServerNode/structs"
	"ServerNode/util"
	"fmt"
	"net/http"
)

func (h *Handler) GetValueHelper(w http.ResponseWriter, Key, Nonce, PreviousNodeHash string) {
	EntryHash := util.Sha256String(Key + Nonce)

	// We've reached the correct node -------------------------------------------
	// Case 1: standard case
	// Case 2: current node has looped back around to 0 and entry belongs in the node with lowest hash
	if h.NodeInfo.NodeHash > EntryHash || h.NodeInfo.NodeHash < PreviousNodeHash && EntryHash > PreviousNodeHash {
		fmt.Printf("[Debug] Node %s is the correct destination for hash %s \n", h.NodeInfo.NodeUrl, EntryHash)
		entry, ok := h.NodeInfo.NodeContents[EntryHash]
		if !ok {
			util.WriteResponse(w, nil, http.StatusNotFound)
			return
		} else {
			util.WriteResponse(w, entry, http.StatusOK)
			return
		}
	}

	// Not the correct node, keep searching -------------------------------------------
	fmt.Printf("[Debug] continuing on the loop\n")
	for i := 0; i < min(h.NodeInfo.StoredNbrs, len(h.NodeInfo.SuccessorArray)); i++ {
		fmt.Printf("[Debug] sending msg to %s\n", h.NodeInfo.SuccessorArray[i])

		// Check the next descendant
		requestEndpoint := fmt.Sprintf("/api/%s/%s/%s", Key, Nonce, h.NodeInfo.NodeHash)
		resp, err := h.Requester.SendRequest(h.NodeInfo.SuccessorArray[i], requestEndpoint, http.MethodGet, nil, constants.REQUEST_TIMEOUT)

		if err != nil {
			// Descendent is unresponsive
			fmt.Printf("[Warning] child %s is not healthy, trying next\n", h.NodeInfo.SuccessorArray[i])

		} else if resp.StatusCode != http.StatusOK {
			// Descendent returns a bad status code, return
			w.WriteHeader(resp.StatusCode)
			return

		} else {
			// Descendent responds pass response forward
			var nodeResponse = &structs.EntryResponse{}
			err := util.ReadBody(resp.Body, nodeResponse)
			// time.Sleep(1 * time.Second)

			if err != nil {
				fmt.Printf("[Error] failed to decode response body: %s\n", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			util.WriteResponse(w, nodeResponse, http.StatusOK)
			return
		}
	}

	fmt.Printf("[Error] all descendants have failed for current NodeUrl |%s|\n", h.NodeInfo.NodeUrl)
	w.WriteHeader(http.StatusInternalServerError)
}
