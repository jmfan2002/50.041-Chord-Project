package api

import (
	"ServerNode/constants"
	"ServerNode/structs"
	"ServerNode/util"

	"fmt"
	"net/http"
)

func (h *Handler) SetValue(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("[Msg] Set value called")

	// Read values from request -------------------------------------------
	reqBody := &structs.SetValueReqBody{}
	err := util.ReadBody(r.Body, reqBody)
	if err != nil {
		fmt.Printf("[ERROR] failed to read request body: %s\n", error.Error)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if reqBody.PreviousNodeHash == "" {
		// fmt.Println("[Debug] Called with no PreviousNodeHash")
		reqBody.PreviousNodeHash = "nil"
	}

	EntryHash := util.Sha256String(reqBody.Key + reqBody.Nonce)
	// fmt.Printf("[Debug] entry hash: %s\n", EntryHash)
	// fmt.Printf("\tprevious node hash: %s\n", reqBody.PreviousNodeHash)
	// fmt.Printf("\tprevious current node hash: %s\n", h.NodeInfo.NodeHash)

	// We've reached the correct node -------------------------------------------
	// Case 1: standard case
	// Case 2: current node has looped back around to 0 and entry belongs in the node with lowest hash
	if h.NodeInfo.NodeHash >= EntryHash || h.NodeInfo.NodeHash < reqBody.PreviousNodeHash && EntryHash > reqBody.PreviousNodeHash {
		insertedEntry := structs.EntryResponse{Key: reqBody.Key, Value: reqBody.Value, Nonce: reqBody.Nonce}
		fmt.Printf("Node %s is the correct destination for entry %s, inserting \n", h.NodeInfo.NodeUrl, insertedEntry)
		h.NodeInfo.NodeContents[EntryHash] = insertedEntry
		util.WriteResponse(w, h.NodeInfo.NodeContents[EntryHash], http.StatusOK)
		return
	}

	// Not the correct node, keep searching -------------------------------------------
	// fmt.Printf("[Debug] continuing on the loop\n")
	for i := 0; i < min(h.NodeInfo.StoredNbrs, len(h.NodeInfo.SuccessorArray)); i++ {
		// fmt.Printf("[Debug] sending msg to %s\n", h.NodeInfo.SuccessorArray[i])

		// Check the next descendant
		requestEndpoint := fmt.Sprintf("/api")
		reqBody.PreviousNodeHash = h.NodeInfo.NodeHash
		resp, err := h.Requester.SendRequest(h.NodeInfo.SuccessorArray[i], requestEndpoint, http.MethodPost, reqBody, constants.REQUEST_TIMEOUT)

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

			if err != nil {
				fmt.Printf("[Error] failed to decode response body: %s\n", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			util.WriteResponse(w, nodeResponse, http.StatusOK)
			return
		}
	}

	fmt.Printf("[Error] all descendants have failed SetValue for current NodeUrl |%s|\n", h.NodeInfo.NodeUrl)
	w.WriteHeader(http.StatusInternalServerError)
}
