package acp

import (
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
	writeJSON(w, http.StatusOK, map[string]any{"products": Feed(h.store)})
}
