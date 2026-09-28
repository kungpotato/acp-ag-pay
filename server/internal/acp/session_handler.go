package acp

import (
	"encoding/json"
	"net/http"
)

type SessionHandler struct {
	sessions *SessionStore
}

func NewSessionHandler(sessions *SessionStore) *SessionHandler {
	return &SessionHandler{sessions: sessions}
}

func (h *SessionHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/acp/checkout_sessions", h.create)
	mux.HandleFunc("GET /api/acp/checkout_sessions/{id}", h.get)
}

type createSessionRequest struct {
	Items []RequestedItem `json:"items"`
}

func (h *SessionHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	session, err := h.sessions.Create(req.Items)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, session)
}

func (h *SessionHandler) get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	session, ok := h.sessions.Get(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "checkout session not found"})
		return
	}
	writeJSON(w, http.StatusOK, session)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
