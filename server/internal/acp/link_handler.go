package acp

import (
	"encoding/json"
	"net/http"

	"acp-ag-pay/server/internal/payment"
)

// LinkHandler exposes the one-time payment-linking step. This is the only
// place in Lesson 9 where the shopper interacts with a checkout-like UI —
// once, ever, not per order.
type LinkHandler struct {
	links   *LinkStore
	settler *payment.StripeSettler
}

func NewLinkHandler(links *LinkStore, settler *payment.StripeSettler) *LinkHandler {
	return &LinkHandler{links: links, settler: settler}
}

func (h *LinkHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/agent/link", h.create)
	mux.HandleFunc("GET /api/agent/links/{id}", h.get)
}

type createLinkRequest struct {
	MaxAmountCents int64  `json:"max_amount_cents"`
	Currency       string `json:"currency"`
}

func (h *LinkHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Currency == "" {
		req.Currency = "thb"
	}
	if req.MaxAmountCents <= 0 {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "max_amount_cents must be positive"})
		return
	}

	setup, err := h.settler.CreateSetupIntent()
	if err == payment.ErrNotConfigured {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"error": "payment linking is not configured: set STRIPE_SECRET_KEY",
		})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "stripe error: " + err.Error()})
		return
	}

	link := h.links.Create(setup.CustomerID, setup.SetupIntentID, setup.ClientSecret, req.MaxAmountCents, req.Currency)
	writeJSON(w, http.StatusCreated, link)
}

func (h *LinkHandler) get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	link, ok := h.links.Get(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "link not found"})
		return
	}
	writeJSON(w, http.StatusOK, link)
}
