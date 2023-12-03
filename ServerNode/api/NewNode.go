package api

import (
	"ServerNode/structs"
	"ServerNode/util"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type NewNodeReq struct {
	// Which node the reuest "started" from, since it's passed down the ring
	Origin string
	// List of ndoes that have viewed this message
	ViewList []string
	// The ip of the new node
	NewNode string
}

func (h *Handler) NewNode(w http.ResponseWriter, r *http.Request) {
	var reqBody = NewNodeReq{}
	err := util.ReadBody(r.Body, &reqBody)
	if err != nil {
		fmt.Println("Error decoding request")
		return
	}

	fmt.Printf("NewNode called! Originally from %s, the new node is %s\n", reqBody.Origin, reqBody.NewNode)

	// Determine if the message made a full loop, and stop if so.
	// This is the case if the view list is non-empty, and we're the original target.
	// (If we're the origin, we would've done all the work we needed
	// to do when we first received this message, so we just return)
	if reqBody.Origin == h.NodeInfo.NodeUrl && len(reqBody.ViewList) > 0 {
		writeOK(w)
		return
	}

	// Oops. The new node is us.
	if reqBody.NewNode == h.NodeInfo.NodeUrl {
		writeOK(w)
		return
	}

	// Add to view list
	reqBody.ViewList = append(reqBody.ViewList, h.NodeInfo.NodeUrl)

	// If we have fewer than #StoredNbrs successors, always try to insert
	// ... otherwise, we need to be a little smarter about inserting
	insertionPoint := 0

	for i := 0; i < len(h.NodeInfo.SuccessorArray); i += 1 {
		if util.Sha256String(h.NodeInfo.SuccessorArray[i]) <
			util.Sha256String(reqBody.NewNode) && !((util.Sha256String(h.NodeInfo.NodeUrl) < util.Sha256String(reqBody.NewNode)) &&
			(util.Sha256String(h.NodeInfo.SuccessorArray[i]) < util.Sha256String(h.NodeInfo.NodeUrl))) {
			insertionPoint = i + 1
		}
	}

	origSuccArr := make([]string, 0)
	origSuccArr = append(origSuccArr, h.NodeInfo.SuccessorArray...)

	// Innocent for loop. Surely, nothing could go wrong here.
	for len(h.NodeInfo.SuccessorArray) > 0 {
		// Pass the message along, and expect a response back
		data, _ := json.Marshal(reqBody)
		_, err := http.Post(
			h.NodeInfo.SuccessorArray[0]+"/api/join",
			"application/json",
			bytes.NewBuffer(data),
		)
		if err == nil {
			fmt.Println("Successfully continuing")
			break
		}
		// If we don't get a response from our successor...
		// Assume it's dead, remove it from our list, and try the next
		// This works for up to #StoredNbrs failures in the general case.
		h.NodeInfo.SuccessorArray = h.NodeInfo.SuccessorArray[1:]
		fmt.Println("Successor failed...")

		// We don't have a great way to recover from this one if we run out of succ...
		if len(h.NodeInfo.SuccessorArray) == 0 {
			fmt.Println("uhhh")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte{})
			return
		}
	}

	h.NodeInfo.SuccessorArray =
		append(
			h.NodeInfo.SuccessorArray[:insertionPoint],
			append(
				[]string{reqBody.NewNode}, h.NodeInfo.SuccessorArray[insertionPoint:]...,
			)...,
		)

	// Limit length of succ array
	h.NodeInfo.SuccessorArray = h.NodeInfo.SuccessorArray[:min(h.NodeInfo.StoredNbrs, len(h.NodeInfo.SuccessorArray))]

	// If this node is the true direct predecessor to the new node...
	if h.NodeInfo.SuccessorArray[0] == reqBody.NewNode {
		fmt.Println("We're the predecessor!")
		// also share its pre-modification successor array
		if len(origSuccArr) < h.NodeInfo.StoredNbrs {
			origSuccArr = append(origSuccArr, h.NodeInfo.NodeUrl)
		}
		data, _ := json.Marshal(
			structs.SuccessorsResponse{
				Successors: origSuccArr,
			})
		_, err = http.Post(
			reqBody.NewNode+"/api/successors",
			"application/json",
			bytes.NewBuffer(data),
		)
		if err != nil {
			fmt.Println(err)
			return
		}

		// And reassign entries!
		// This is perhaps the most naive way to do this but uh, it should work
		h.ReassignEntries(w, r)
	}

	fmt.Printf("New node join detected. New successors: %s\n", h.NodeInfo.SuccessorArray)

	// At last. Freedom.
	writeOK(w)
}

func writeOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}
