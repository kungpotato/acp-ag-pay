package acp

import (
	"encoding/json"
	"net/http"
)

// DevHandler exposes workshop-only helpers that have no place in a real
// deployment. Lessons 6, 7 and 9 need to demonstrate order status polling,
// webhook confirmation, and an already-linked payment mandate, but this
// workshop never has a real Stripe account wired up — without a genuine
// sk_test_ key there is no way to reach these states through the real
// /api/agent/pay or /api/agent/link flows. This handler lets the workshop
// seed that state directly so the later lessons can be demoed honestly. It
// is registered unconditionally here for simplicity, but a real deployment
// would delete this file entirely or gate it behind a dev-only build tag.
type DevHandler struct {
	orders *OrderStore
	links  *LinkStore
}

func NewDevHandler(orders *OrderStore, links *LinkStore) *DevHandler {
	return &DevHandler{orders: orders, links: links}
}

func (h *DevHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/dev/seed_order", h.seedOrder)
	mux.HandleFunc("POST /api/dev/seed_link", h.seedLink)
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

type seedLinkRequest struct {
	MaxAmountCents  int64  `json:"max_amount_cents"`
	Currency        string `json:"currency"`
	CustomerID      string `json:"customer_id"`
	PaymentMethodID string `json:"payment_method_id"`
}

// seedLink creates an already-ACTIVE link, skipping the real Stripe
// Customer/SetupIntent round trip. This stands in for a shopper who
// completed the one-time linking UI in Lesson 9's real flow.
func (h *DevHandler) seedLink(w http.ResponseWriter, r *http.Request) {
	var req seedLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Currency == "" {
		req.Currency = "thb"
	}
	if req.CustomerID == "" {
		req.CustomerID = "cus_seeded_" + newLinkID()
	}
	if req.PaymentMethodID == "" {
		req.PaymentMethodID = "pm_seeded_" + newLinkID()
	}
	setupIntentID := "seti_seeded_" + newLinkID()

	h.links.Create(req.CustomerID, setupIntentID, "", req.MaxAmountCents, req.Currency)
	activated, _ := h.links.Activate(setupIntentID, req.PaymentMethodID)
	writeJSON(w, http.StatusCreated, activated)
}
