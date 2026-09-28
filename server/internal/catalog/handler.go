package catalog

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Handler exposes the catalog over HTTP for the Next.js frontend (Lesson 1)
// and later for the shopping agent (Lesson 2+).
type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/books", h.listBooks)
	mux.HandleFunc("GET /api/books/{id}", h.getBook)
}

func (h *Handler) listBooks(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	books := h.store.Search(q)
	writeJSON(w, http.StatusOK, books)
}

func (h *Handler) getBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	book, ok := h.store.Get(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "book not found"})
		return
	}
	writeJSON(w, http.StatusOK, book)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
