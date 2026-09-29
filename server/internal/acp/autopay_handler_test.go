package acp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"acp-ag-pay/server/internal/catalog"
	"acp-ag-pay/server/internal/payment"
)

func TestAutopayReturns503WhenStripeNotConfigured(t *testing.T) {
	catalogStore := catalog.NewStore()
	sessions := NewSessionStore(catalogStore)
	orders := NewOrderStore()
	links := NewLinkStore()

	session, err := sessions.Create([]RequestedItem{{ProductID: "bk-001", Quantity: 1}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	links.Create("cus_1", "seti_1", "", 1_000_000, "thb")
	link, _ := links.Activate("seti_1", "pm_1")

	settler := payment.NewStripeSettler("")
	mux := http.NewServeMux()
	NewAutopayHandler(sessions, orders, links, settler).Register(mux)

	body, _ := json.Marshal(autopayRequest{SessionID: session.ID, LinkID: link.ID})
	req := httptest.NewRequest(http.MethodPost, "/api/agent/autopay", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestAutopayRejectsAmountAboveMandate(t *testing.T) {
	catalogStore := catalog.NewStore()
	sessions := NewSessionStore(catalogStore)
	orders := NewOrderStore()
	links := NewLinkStore()

	book, _ := catalogStore.Get("bk-002") // more expensive book
	session, err := sessions.Create([]RequestedItem{{ProductID: book.ID, Quantity: 1}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	links.Create("cus_1", "seti_1", "", 1, "thb") // mandate capped far below price
	link, _ := links.Activate("seti_1", "pm_1")

	settler := payment.NewStripeSettler("")
	mux := http.NewServeMux()
	NewAutopayHandler(sessions, orders, links, settler).Register(mux)

	body, _ := json.Marshal(autopayRequest{SessionID: session.ID, LinkID: link.ID})
	req := httptest.NewRequest(http.MethodPost, "/api/agent/autopay", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestAutopayRejectsUnknownSession(t *testing.T) {
	sessions := NewSessionStore(catalog.NewStore())
	orders := NewOrderStore()
	links := NewLinkStore()
	settler := payment.NewStripeSettler("")
	mux := http.NewServeMux()
	NewAutopayHandler(sessions, orders, links, settler).Register(mux)

	body, _ := json.Marshal(autopayRequest{SessionID: "cs_missing", LinkID: "link_missing"})
	req := httptest.NewRequest(http.MethodPost, "/api/agent/autopay", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}
