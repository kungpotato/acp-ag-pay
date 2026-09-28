package agent

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	agent *Agent
}

func NewHandler(agent *Agent) *Handler {
	return &Handler{agent: agent}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/agent/chat", h.chat)
}

type chatRequest struct {
	Message string `json:"message"`
}

func (h *Handler) chat(w http.ResponseWriter, r *http.Request) {
	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	reply := h.agent.Handle(r.Context(), req.Message)
	writeJSON(w, http.StatusOK, reply)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
