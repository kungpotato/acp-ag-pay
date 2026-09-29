package acp

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// LinkStatus tracks the one-time payment-linking flow that makes real
// agentic checkout possible: a shopper authorizes a payment method ONCE
// (LinkPending -> LinkActive), and afterwards the agent can charge that
// payment method autonomously, with no further checkout UI per order.
// Contrast with SessionStatus/OrderStatus (Lessons 4-6), which describe a
// single order; a Link outlives many orders.
type LinkStatus string

const (
	LinkPending LinkStatus = "pending"
	LinkActive  LinkStatus = "active"
	LinkRevoked LinkStatus = "revoked"
)

type Link struct {
	ID              string     `json:"id"`
	Status          LinkStatus `json:"status"`
	MaxAmountCents  int64      `json:"max_amount_cents"`
	Currency        string     `json:"currency"`
	CustomerID      string     `json:"customer_id"`
	SetupIntentID   string     `json:"setup_intent_id"`
	ClientSecret    string     `json:"client_secret,omitempty"`
	PaymentMethodID string     `json:"payment_method_id,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	ActivatedAt     *time.Time `json:"activated_at,omitempty"`
}

type LinkStore struct {
	mu            sync.Mutex
	links         map[string]Link
	bySetupIntent map[string]string
}

func NewLinkStore() *LinkStore {
	return &LinkStore{
		links:         make(map[string]Link),
		bySetupIntent: make(map[string]string),
	}
}

func (s *LinkStore) Create(customerID, setupIntentID, clientSecret string, maxAmountCents int64, currency string) Link {
	s.mu.Lock()
	defer s.mu.Unlock()

	link := Link{
		ID:             newLinkID(),
		Status:         LinkPending,
		MaxAmountCents: maxAmountCents,
		Currency:       currency,
		CustomerID:     customerID,
		SetupIntentID:  setupIntentID,
		ClientSecret:   clientSecret,
		CreatedAt:      time.Now().UTC(),
	}
	s.links[link.ID] = link
	s.bySetupIntent[setupIntentID] = link.ID
	return link
}

func (s *LinkStore) Get(id string) (Link, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	l, ok := s.links[id]
	return l, ok
}

// Activate is called only from the webhook handler (Lesson 7/9) once
// Stripe confirms `setup_intent.succeeded` — never from a client request,
// since a client claiming "I linked successfully" is not trustworthy on its
// own.
func (s *LinkStore) Activate(setupIntentID, paymentMethodID string) (Link, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id, ok := s.bySetupIntent[setupIntentID]
	if !ok {
		return Link{}, false
	}
	link := s.links[id]
	link.Status = LinkActive
	link.PaymentMethodID = paymentMethodID
	now := time.Now().UTC()
	link.ActivatedAt = &now
	s.links[id] = link
	return link, true
}

// Authorize checks that a link is active and the requested amount is
// within the mandate the shopper originally consented to. This is what
// stands in for Stripe's Shared Payment Token scoping in this workshop's
// simplified model.
func (s *LinkStore) Authorize(id string, amountCents int64) (Link, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	link, ok := s.links[id]
	if !ok {
		return Link{}, fmt.Errorf("link %s not found", id)
	}
	if link.Status != LinkActive {
		return Link{}, fmt.Errorf("link %s is not active (status: %s)", id, link.Status)
	}
	if amountCents > link.MaxAmountCents {
		return Link{}, fmt.Errorf("amount %d exceeds link's authorized max %d", amountCents, link.MaxAmountCents)
	}
	return link, nil
}

func newLinkID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return "link_" + hex.EncodeToString(b)
}
