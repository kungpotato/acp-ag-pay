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

func TestPayHandlerReturns503WhenStripeNotConfigured(t *testing.T) {
	catalogStore := catalog.NewStore()
	sessions := NewSessionStore(catalogStore)
	session, err := sessions.Create([]RequestedItem{{ProductID: "bk-001", Quantity: 1}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	settler := payment.NewStripeSettler("") // unconfigured, matches .env.example default
	mux := http.NewServeMux()
	NewPayHandler(sessions, settler).Register(mux)

	body, _ := json.Marshal(payRequest{SessionID: session.ID})
	req := httptest.NewRequest(http.MethodPost, "/api/agent/pay", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestPayHandlerRejectsUnknownSession(t *testing.T) {
	sessions := NewSessionStore(catalog.NewStore())
	settler := payment.NewStripeSettler("")
	mux := http.NewServeMux()
	NewPayHandler(sessions, settler).Register(mux)

	body, _ := json.Marshal(payRequest{SessionID: "cs_does_not_exist"})
	req := httptest.NewRequest(http.MethodPost, "/api/agent/pay", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}
