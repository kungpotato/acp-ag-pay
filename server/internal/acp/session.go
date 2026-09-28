package acp

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"acp-ag-pay/server/internal/catalog"
)

// SessionStatus tracks where a checkout session is in its ACP-inspired
// lifecycle: pending (cart assembled, stock reserved) -> ready_for_payment
// (Lesson 5 attaches a Stripe PaymentIntent) -> completed (Lesson 7 webhook
// confirms settlement) or canceled.
type SessionStatus string

const (
	StatusPending         SessionStatus = "pending"
	StatusReadyForPayment SessionStatus = "ready_for_payment"
	StatusCompleted       SessionStatus = "completed"
	StatusCanceled        SessionStatus = "canceled"
)

type LineItem struct {
	ProductID      string `json:"product_id"`
	Title          string `json:"title"`
	Quantity       int    `json:"quantity"`
	UnitPriceCents int64  `json:"unit_price_cents"`
}

type CheckoutSession struct {
	ID              string        `json:"id"`
	Status          SessionStatus `json:"status"`
	LineItems       []LineItem    `json:"line_items"`
	TotalCents      int64         `json:"total_cents"`
	Currency        string        `json:"currency"`
	CreatedAt       time.Time     `json:"created_at"`
	PaymentIntentID string        `json:"payment_intent_id,omitempty"`
	ClientSecret    string        `json:"client_secret,omitempty"`
}

// RequestedItem is what the client sends when creating a session: just a
// product id and quantity. Price and title are always resolved server-side
// from the catalog, never trusted from the client.
type RequestedItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type SessionStore struct {
	mu       sync.Mutex
	sessions map[string]CheckoutSession
	catalog  *catalog.Store
}

func NewSessionStore(catalogStore *catalog.Store) *SessionStore {
	return &SessionStore{
		sessions: make(map[string]CheckoutSession),
		catalog:  catalogStore,
	}
}

// Create builds a checkout session from requested items, resolving price
// and title from the catalog and reserving stock so two shoppers can't
// oversell the last copy of a book while both are mid-checkout.
func (s *SessionStore) Create(items []RequestedItem) (CheckoutSession, error) {
	if len(items) == 0 {
		return CheckoutSession{}, fmt.Errorf("checkout session needs at least one item")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	lineItems := make([]LineItem, 0, len(items))
	var total int64
	var currency string
	reserved := make([]RequestedItem, 0, len(items))

	rollback := func() {
		for _, r := range reserved {
			// best-effort release; ReserveStock has no inverse in Lesson 1,
			// so Lesson 4 adds it here since checkout is the first place we
			// need to give stock back.
			_ = s.catalog.ReleaseStock(r.ProductID, r.Quantity)
		}
	}

	for _, item := range items {
		if item.Quantity <= 0 {
			rollback()
			return CheckoutSession{}, fmt.Errorf("quantity must be positive for %s", item.ProductID)
		}
		book, ok := s.catalog.Get(item.ProductID)
		if !ok {
			rollback()
			return CheckoutSession{}, fmt.Errorf("unknown product %s", item.ProductID)
		}
		if err := s.catalog.ReserveStock(item.ProductID, item.Quantity); err != nil {
			rollback()
			return CheckoutSession{}, err
		}
		reserved = append(reserved, item)

		lineItems = append(lineItems, LineItem{
			ProductID:      book.ID,
			Title:          book.Title,
			Quantity:       item.Quantity,
			UnitPriceCents: book.PriceCents,
		})
		total += book.PriceCents * int64(item.Quantity)
		currency = book.Currency
	}

	session := CheckoutSession{
		ID:         newSessionID(),
		Status:     StatusPending,
		LineItems:  lineItems,
		TotalCents: total,
		Currency:   currency,
		CreatedAt:  time.Now().UTC(),
	}
	s.sessions[session.ID] = session
	return session, nil
}

func (s *SessionStore) Get(id string) (CheckoutSession, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session, ok := s.sessions[id]
	return session, ok
}

// Update applies fn to the stored session under lock, used by Lesson 5 to
// attach Stripe PaymentIntent details and by Lesson 7 to mark completion.
func (s *SessionStore) Update(id string, fn func(*CheckoutSession)) (CheckoutSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session, ok := s.sessions[id]
	if !ok {
		return CheckoutSession{}, fmt.Errorf("session %s not found", id)
	}
	fn(&session)
	s.sessions[id] = session
	return session, nil
}

func newSessionID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return "cs_" + hex.EncodeToString(b)
}
