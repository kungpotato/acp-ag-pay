package acp

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// OrderStatus tracks the shopper-facing state of an order, which is
// deliberately a smaller, more honest vocabulary than SessionStatus:
// shoppers don't need to know about "ready_for_payment", they need to know
// whether their money actually went through.
type OrderStatus string

const (
	// OrderPendingConfirmation is the status the moment a PaymentIntent is
	// created. Stripe has *accepted* the request, but settlement isn't
	// final until the webhook in Lesson 7 confirms it — a PaymentIntent
	// can still fail (card declined, 3DS abandoned, etc.) after this point.
	OrderPendingConfirmation OrderStatus = "pending_confirmation"
	OrderPaid                OrderStatus = "paid"
	OrderFailed              OrderStatus = "failed"
)

type Order struct {
	ID              string      `json:"id"`
	SessionID       string      `json:"session_id"`
	Status          OrderStatus `json:"status"`
	TotalCents      int64       `json:"total_cents"`
	Currency        string      `json:"currency"`
	PaymentIntentID string      `json:"payment_intent_id"`
	CreatedAt       time.Time   `json:"created_at"`
	ConfirmedAt     *time.Time  `json:"confirmed_at,omitempty"`
}

type OrderStore struct {
	mu     sync.Mutex
	orders map[string]Order
	// byPaymentIntent lets the webhook handler (Lesson 7) find an order by
	// the Stripe PaymentIntent id carried in the webhook payload, without
	// needing our own order id.
	byPaymentIntent map[string]string
}

func NewOrderStore() *OrderStore {
	return &OrderStore{
		orders:          make(map[string]Order),
		byPaymentIntent: make(map[string]string),
	}
}

func (s *OrderStore) CreateFromSession(session CheckoutSession) Order {
	s.mu.Lock()
	defer s.mu.Unlock()

	order := Order{
		ID:              newOrderID(),
		SessionID:       session.ID,
		Status:          OrderPendingConfirmation,
		TotalCents:      session.TotalCents,
		Currency:        session.Currency,
		PaymentIntentID: session.PaymentIntentID,
		CreatedAt:       time.Now().UTC(),
	}
	s.orders[order.ID] = order
	s.byPaymentIntent[session.PaymentIntentID] = order.ID
	return order
}

func (s *OrderStore) Get(id string) (Order, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.orders[id]
	return o, ok
}

func (s *OrderStore) GetByPaymentIntent(paymentIntentID string) (Order, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, ok := s.byPaymentIntent[paymentIntentID]
	if !ok {
		return Order{}, false
	}
	o, ok := s.orders[id]
	return o, ok
}

// Confirm marks an order paid or failed. Lesson 7's webhook handler is the
// only caller: order status only ever changes because Stripe told us so,
// never because the checkout response looked successful.
func (s *OrderStore) Confirm(id string, status OrderStatus) (Order, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.orders[id]
	if !ok {
		return Order{}, false
	}
	o.Status = status
	now := time.Now().UTC()
	o.ConfirmedAt = &now
	s.orders[id] = o
	return o, true
}

func newOrderID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return "ord_" + hex.EncodeToString(b)
}
