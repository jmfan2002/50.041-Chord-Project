package api

import (
	"ServerNode/constants"
	"ServerNode/structs"
	"ServerNode/util"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) UpdateSuccessors(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("[Msg] UpdateSuccessors called\n")

	// Parse variables from url -------------------------------------------
	StartingNodeHash := mux.Vars(r)["StartingNodeHash"]

	CurrentOverlap, err := strconv.Atoi(mux.Vars(r)["CurrentOverlap"])
	if err != nil {
		fmt.Printf("[ERROR] cannot convert CurrentOverlap to string: %s", mux.Vars(r)["CurrentOverlap"])
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// We've fully overlapped, return -------------------------------------------
	if CurrentOverlap == h.NodeInfo.StoredNbrs-1 {
		// fmt.Printf("[Debug] overlap complete, returning! \n")
		util.WriteResponse(w, structs.SuccessorsResponse{Successors: []string{h.NodeInfo.NodeUrl}}, http.StatusOK)
		return
	}

	// Continue on overlap -------------------------------------------
	if CurrentOverlap > 0 {
		CurrentOverlap++
	}

	// We've cycled back, start on overlap -------------------------------------------
	if StartingNodeHash != "nil" && CurrentOverlap == 0 && h.NodeInfo.NodeHash == StartingNodeHash {
		// fmt.Printf("[Debug] cycle complete, starting on overlap \n")
		CurrentOverlap++
	}
	if StartingNodeHash == "nil" {
		StartingNodeHash = h.NodeInfo.NodeHash
	}

	if CurrentOverlap == 0 {
		// fmt.Printf("[Debug] successor array starts as: %s\n", h.NodeInfo.SuccessorArray)
	}

	// Check descendants -------------------------------------------
	for i := 0; i < min(h.NodeInfo.StoredNbrs, len(h.NodeInfo.SuccessorArray)); i++ {
		// fmt.Printf("[Debug] sending msg to %s\n", h.NodeInfo.SuccessorArray[i])

		// Request next descendent
		requestEndpoint := fmt.Sprintf("/api/successors/%s/%d", StartingNodeHash, CurrentOverlap)
		resp, err := h.Requester.SendRequest(h.NodeInfo.SuccessorArray[i], requestEndpoint, http.MethodPatch, nil, constants.REQUEST_TIMEOUT)

		if err != nil {
			// Descendent is unresponsive
			// fmt.Printf("[Debug] child %s is not healthy, trying next\n", h.NodeInfo.SuccessorArray[i])

		} else if resp.StatusCode != http.StatusOK {
			// Descendent returns a bad status code, return
			fmt.Printf("[Debug] next node is reporting break in cycle with status code: %d\n", resp.StatusCode)
			w.WriteHeader(resp.StatusCode)
			return

		} else {
			// Descendent responds ok, process response
			var updateResp = &structs.SuccessorsResponse{}
			err := util.ReadBody(resp.Body, updateResp)
			// time.Sleep(1 * time.Second)

			if err != nil {
				fmt.Printf("[Error] failed to decode response body: %s\n", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Filter out self references (loop back to original node)
			for idx := 0; idx < len(updateResp.Successors); idx++ {
				if updateResp.Successors[idx] == h.NodeInfo.NodeUrl {
					updateResp.Successors = updateResp.Successors[:idx]
					break
				}
			}

			// If it's not an overlap, update array
			if CurrentOverlap == 0 {
				h.NodeInfo.SuccessorArray = util.CopySliceString(updateResp.Successors)
				// fmt.Printf("[Debug] successor array is now: %s\n", h.NodeInfo.SuccessorArray)
			}

			// Regardless, return array of current node and previous k - 1
			updateResp.Successors = append([]string{h.NodeInfo.NodeUrl}, updateResp.Successors...)
			if len(updateResp.Successors) >= h.NodeInfo.StoredNbrs {
				updateResp.Successors = updateResp.Successors[:h.NodeInfo.StoredNbrs-1]
			}

			// fmt.Printf("[Debug] returning successor array: %s\n", updateResp.Successors)
			util.WriteResponse(w, updateResp, http.StatusOK)
			return
		}
	}

	fmt.Printf("[Error] all descendants have failed UpdateSuccessors for current NodeUrl |%s|\n", h.NodeInfo.NodeUrl)
	// util.WriteResponse(w, structs.SuccessorsResponse{Successors: []string{h.NodeInfo.NodeUrl}}, http.StatusOK)
	w.WriteHeader(http.StatusInternalServerError)
}
