package acp

import (
	"encoding/json"
	"net/http"
)

// DevHandler exposes workshop-only helpers that have no place in a real
// deployment. Lessons 6 and 7 need to demonstrate order status polling and
// webhook confirmation, but this workshop never has a real Stripe account
// wired up — without a genuine sk_test_ key there is no way to reach
// OrderPendingConfirmation through the real /api/agent/pay flow. This
// handler lets the workshop seed an order directly so the frontend and
// webhook lessons can be demoed honestly. It is registered unconditionally
// here for simplicity, but a real deployment would delete this file
// entirely or gate it behind a dev-only build tag.
type DevHandler struct {
	orders *OrderStore
}

func NewDevHandler(orders *OrderStore) *DevHandler {
	return &DevHandler{orders: orders}
}

func (h *DevHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/dev/seed_order", h.seedOrder)
}

type seedOrderRequest struct {
	TotalCents      int64  `json:"total_cents"`
	Currency        string `json:"currency"`
	PaymentIntentID string `json:"payment_intent_id"`
}

func (h *DevHandler) seedOrder(w http.ResponseWriter, r *http.Request) {
	var req seedOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Currency == "" {
		req.Currency = "thb"
	}
	if req.PaymentIntentID == "" {
		req.PaymentIntentID = "pi_seeded_" + newOrderID()
	}

	order := h.orders.CreateFromSession(CheckoutSession{
		ID:              "cs_seeded",
		TotalCents:      req.TotalCents,
		Currency:        req.Currency,
		PaymentIntentID: req.PaymentIntentID,
	})
	writeJSON(w, http.StatusCreated, order)
}
