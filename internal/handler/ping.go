package handler

import (
	"net/http"

	"desent-api-quest/internal/httpjson"
)

type PingHandler struct{}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (h *PingHandler) Handle(w http.ResponseWriter, _ *http.Request) {
	httpjson.WriteJSON(w, http.StatusOK, map[string]bool{"success": true})
}
