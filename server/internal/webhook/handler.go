// Package webhook receives asynchronous confirmation from Stripe: the last
// step of ACP settlement. A successful /api/agent/pay call (Lesson 5) only
// means Stripe *accepted* the PaymentIntent request — it does not mean the
// card actually charged. Only a verified webhook event is trustworthy
// enough to mark an order paid.
package webhook

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"

	"acp-ag-pay/server/internal/acp"
)

type Handler struct {
	orders        *acp.OrderStore
	signingSecret string
}

func NewHandler(orders *acp.OrderStore, signingSecret string) *Handler {
	return &Handler{orders: orders, signingSecret: signingSecret}
}

func (h *Handler) Configured() bool { return h.signingSecret != "" }

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/webhooks/stripe", h.handle)
}

func (h *Handler) handle(w http.ResponseWriter, r *http.Request) {
	if !h.Configured() {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"error": "webhook is not configured: set STRIPE_WEBHOOK_SECRET",
		})
		return
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "failed to read body"})
		return
	}

	// ConstructEvent verifies the Stripe-Signature header against our
	// signing secret. This is the only reason a webhook is trustworthy
	// input: anyone can POST a fake payment_intent.succeeded payload to a
	// public URL, but they can't forge this signature without the secret.
	event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"), h.signingSecret)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "signature verification failed"})
		return
	}

	switch event.Type {
	case "payment_intent.succeeded":
		h.confirmOrder(event, acp.OrderPaid)
	case "payment_intent.payment_failed":
		h.confirmOrder(event, acp.OrderFailed)
	default:
		log.Printf("webhook: ignoring unhandled event type %s", event.Type)
	}

	writeJSON(w, http.StatusOK, map[string]bool{"received": true})
}

func (h *Handler) confirmOrder(event stripe.Event, status acp.OrderStatus) {
	var pi stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
		log.Printf("webhook: failed to parse payment_intent from event %s: %v", event.ID, err)
		return
	}

	order, ok := h.orders.GetByPaymentIntent(pi.ID)
	if !ok {
		log.Printf("webhook: no order found for payment_intent %s", pi.ID)
		return
	}

	if _, ok := h.orders.Confirm(order.ID, status); ok {
		log.Printf("webhook: order %s confirmed as %s", order.ID, status)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
