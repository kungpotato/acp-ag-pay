package acp

import (
	"encoding/json"
	"net/http"

	"acp-ag-pay/server/internal/catalog"
)

type FeedHandler struct {
	store *catalog.Store
}

func NewFeedHandler(store *catalog.Store) *FeedHandler {
	return &FeedHandler{store: store}
}

func (h *FeedHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/acp/products", h.listProducts)
}

func (h *FeedHandler) listProducts(w http.ResponseWriter, r *http.Request) {
	feed := Feed(h.store)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"products": feed,
	})
}
