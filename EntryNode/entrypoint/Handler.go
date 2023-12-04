package entrypoint

import (
	"EntryNode/util"
)

type Handler struct {
	EntryPoint EntryPoint
	Requester  util.Requester
}

func NewHandler(k int) *Handler {
	return &Handler{
		EntryPoint: *New(k),
		Requester:  &util.HeartbeatRequester{},
	}
}
