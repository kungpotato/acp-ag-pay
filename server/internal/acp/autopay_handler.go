package acp

import (
	"encoding/json"
	"errors"
	"net/http"

	"acp-ag-pay/server/internal/payment"
)

// AutopayHandler is Lesson 9's real point: the agent settles a checkout
// session using a previously linked payment method, with NO shopper-facing
// checkout step. Compare to PayHandler (Lesson 5), which returns a client
// secret a human still has to confirm.
type AutopayHandler struct {
	sessions *SessionStore
	orders   *OrderStore
	links    *LinkStore
	settler  *payment.StripeSettler
}

func NewAutopayHandler(sessions *SessionStore, orders *OrderStore, links *LinkStore, settler *payment.StripeSettler) *AutopayHandler {
	return &AutopayHandler{sessions: sessions, orders: orders, links: links, settler: settler}
}

func (h *AutopayHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/agent/autopay", h.autopay)
}

type autopayRequest struct {
	SessionID string `json:"session_id"`
	LinkID    string `json:"link_id"`
}

func (h *AutopayHandler) autopay(w http.ResponseWriter, r *http.Request) {
	var req autopayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	session, ok := h.sessions.Get(req.SessionID)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "checkout session not found"})
		return
	}
	if session.Status != StatusPending {
		writeJSON(w, http.StatusConflict, map[string]string{
			"error": "checkout session is not pending, current status: " + string(session.Status),
		})
		return
	}

	// This is the mandate check: the agent can only spend what the
	// shopper pre-authorized when they linked their payment method, on
	// the same currency they authorized it in.
	link, err := h.links.Authorize(req.LinkID, session.TotalCents)
	if err != nil {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
		return
	}

	settlement, err := h.settler.CreateOffSessionPaymentIntent(
		link.CustomerID, link.PaymentMethodID, session.TotalCents, session.Currency, session.ID,
	)
	if errors.Is(err, payment.ErrNotConfigured) {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"error": "payment settlement is not configured: set STRIPE_SECRET_KEY",
		})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "stripe error: " + err.Error()})
		return
	}

	updated, err := h.sessions.Update(session.ID, func(s *CheckoutSession) {
		s.Status = StatusReadyForPayment
		s.PaymentIntentID = settlement.PaymentIntentID
		s.ClientSecret = settlement.ClientSecret
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	order := h.orders.CreateFromSession(updated)
	writeJSON(w, http.StatusOK, payResponse{Session: updated, Order: order})
}
