package api

import (
	"ServerNode/constants"
	"ServerNode/structs"
	"fmt"
	"net/http"
)

func (h *Handler) ReassignEntries(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("[Msg] ReassignEntries called\n")
	// First, copy over values
	previousEntries := make([]structs.EntryResponse, 0)
	for _, val := range h.NodeInfo.NodeContents {
		previousEntries = append(previousEntries, val)
	}

	// Next, clear the map
	h.NodeInfo.NodeContents = make(map[string]structs.EntryResponse)

	// Finally, call insert on each entry
	for _, entry := range previousEntries {
		// fmt.Printf("[Debug] re-inserting entry %s\n", entry)

		// Check the next descendant
		for i := 0; i < min(h.NodeInfo.StoredNbrs, len(h.NodeInfo.SuccessorArray)); i++ {
			requestEndpoint := fmt.Sprintf("/api")
			requestBody := structs.NewSetValueReqBody(entry.Key, entry.Value, entry.Nonce)
			resp, err := h.Requester.SendRequest(h.NodeInfo.SuccessorArray[i], requestEndpoint, http.MethodPost, requestBody, constants.REQUEST_TIMEOUT)

			if err != nil {
				// Descendent is unresponsive
				fmt.Printf("[Warning] child %s is not healthy, trying next\n", h.NodeInfo.SuccessorArray[i])

			} else if resp.StatusCode != http.StatusOK {
				// Descendent returns a bad status code, return
				fmt.Printf("[ERROR] child %s returned an error\n", h.NodeInfo.SuccessorArray[i])
				w.WriteHeader(resp.StatusCode)
				return

			} else {
				// fmt.Printf("[Debug] successfully inserted %s\n", entry)
				break
			}
		}
	}
}
