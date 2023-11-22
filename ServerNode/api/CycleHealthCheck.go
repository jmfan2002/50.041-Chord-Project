package api

import (
	"ServerNode/constants"
	"ServerNode/util"
	"fmt"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CycleHealthResponse struct {
	CycleSize int
}

func (h *Handler) CycleHealthCheck(w http.ResponseWriter, r *http.Request) {
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
		WriteResponse(w, CycleHealthResponse{CycleSize: 1}, http.StatusOK)
		return

	} else {
		fmt.Printf("[Debug] continuing on the loop\n")
		for i := 0; i < min(h.NodeInfo.StoredNbrs, len(h.NodeInfo.SuccessorArray)); i++ {
			fmt.Printf("[Debug] sending msg to %s\n", h.NodeInfo.SuccessorArray[i])
			// Case: we've looped
			if util.Sha256String(h.NodeInfo.SuccessorArray[i]) <= StartingNodeHash {
				fmt.Printf("[Debug] setting finishedLoop to true, next node hash of %s <= starting hash of %s\n", util.Sha256String(h.NodeInfo.SuccessorArray[i]), StartingNodeHash)
				FinishedLoop = true
			}

			// Case: check the next descendant with timeout
			requestEndpoint := fmt.Sprintf("/api/cycleHealth/%s/%t", StartingNodeHash, FinishedLoop)
			resp, err := h.Requester.SendRequest(h.NodeInfo.SuccessorArray[i], requestEndpoint, http.MethodGet, constants.REQUEST_TIMEOUT)

			if err != nil {
				fmt.Printf("[Debug] child %s is not healthy, trying next\n", h.NodeInfo.SuccessorArray[i])

			} else if resp.StatusCode != http.StatusOK {
				fmt.Printf("[Debug] next node is reporting break in cycle with status code: %d\n", resp.StatusCode)
				w.WriteHeader(resp.StatusCode)
				return

			} else {
				healthResp, err := readResponseBody(resp)
				if (err != nil) {
					fmt.Printf("[Error] failed to decode response body: %s\n", err.Error())
				}

				healthResp.CycleSize++
				fmt.Printf("[Debug] received healthCheck response from child! nodes: %d\n", healthResp.CycleSize)
				WriteResponse(w, healthResp, http.StatusOK)
				return
			}
		}
	}

	fmt.Printf("[Error] all descendants have failed for current NodeUrl |%s|\n", h.NodeInfo.NodeUrl)
	w.WriteHeader(http.StatusInternalServerError)
}

func readResponseBody(resp *http.Response) (*CycleHealthResponse, error) {
	defer resp.Body.Close()
	var cycleHealthResponse CycleHealthResponse

	err := json.NewDecoder(resp.Body).Decode(&cycleHealthResponse)
	if err != nil {
		return nil, err
	}

	return &cycleHealthResponse, nil
}

func WriteResponse(w http.ResponseWriter, response interface{}, httpCode int) (error) {
	marshalledResp, err := json.Marshal(response)
	if err != nil {
		fmt.Println("error marshalling data")
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(marshalledResp)
	return nil
}