package acp

import (
	"encoding/json"
	"errors"
	"net/http"

	"acp-ag-pay/server/internal/payment"
)

// PayHandler is the ACP "settle" step: it turns a pending checkout session
// into a Stripe PaymentIntent. This is deliberately a separate endpoint
// from session creation (Lesson 4) — a shopping agent decides to settle
// only after the shopper confirms the cart.
type PayHandler struct {
	sessions *SessionStore
	orders   *OrderStore
	settler  *payment.StripeSettler
}

func NewPayHandler(sessions *SessionStore, orders *OrderStore, settler *payment.StripeSettler) *PayHandler {
	return &PayHandler{sessions: sessions, orders: orders, settler: settler}
}

type payResponse struct {
	Session CheckoutSession `json:"session"`
	Order   Order           `json:"order"`
}

func (h *PayHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/agent/pay", h.pay)
}

type payRequest struct {
	SessionID string `json:"session_id"`
}

func (h *PayHandler) pay(w http.ResponseWriter, r *http.Request) {
	var req payRequest
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

	settlement, err := h.settler.CreatePaymentIntent(session.TotalCents, session.Currency, session.ID)
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
