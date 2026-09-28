package acp

import (
	"testing"

	"acp-ag-pay/server/internal/catalog"
)

func TestSessionCreateComputesTotal(t *testing.T) {
	catalogStore := catalog.NewStore()
	sessions := NewSessionStore(catalogStore)

	book, _ := catalogStore.Get("bk-001")
	session, err := sessions.Create([]RequestedItem{{ProductID: "bk-001", Quantity: 2}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if session.TotalCents != book.PriceCents*2 {
		t.Fatalf("expected total %d, got %d", book.PriceCents*2, session.TotalCents)
	}
	if session.Status != StatusPending {
		t.Fatalf("expected pending status, got %s", session.Status)
	}
}

func TestSessionCreateReservesStock(t *testing.T) {
	catalogStore := catalog.NewStore()
	sessions := NewSessionStore(catalogStore)

	before, _ := catalogStore.Get("bk-001")
	if _, err := sessions.Create([]RequestedItem{{ProductID: "bk-001", Quantity: 3}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	after, _ := catalogStore.Get("bk-001")
	if after.Stock != before.Stock-3 {
		t.Fatalf("expected stock to drop by 3, got %d -> %d", before.Stock, after.Stock)
	}
}

func TestSessionCreateRollsBackOnInsufficientStock(t *testing.T) {
	catalogStore := catalog.NewStore()
	sessions := NewSessionStore(catalogStore)

	before, _ := catalogStore.Get("bk-001")
	_, err := sessions.Create([]RequestedItem{
		{ProductID: "bk-001", Quantity: 1},
		{ProductID: "bk-003", Quantity: 999}, // impossible quantity
	})
	if err == nil {
		t.Fatal("expected error for impossible quantity")
	}
	after, _ := catalogStore.Get("bk-001")
	if after.Stock != before.Stock {
		t.Fatalf("expected bk-001 stock to be rolled back, before=%d after=%d", before.Stock, after.Stock)
	}
}

func TestSessionCreateRejectsUnknownProduct(t *testing.T) {
	sessions := NewSessionStore(catalog.NewStore())
	if _, err := sessions.Create([]RequestedItem{{ProductID: "does-not-exist", Quantity: 1}}); err == nil {
		t.Fatal("expected error for unknown product")
	}
}

func TestSessionGetRoundTrip(t *testing.T) {
	sessions := NewSessionStore(catalog.NewStore())
	created, err := sessions.Create([]RequestedItem{{ProductID: "bk-002", Quantity: 1}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := sessions.Get(created.ID)
	if !ok || got.ID != created.ID {
		t.Fatalf("expected to find session %s", created.ID)
	}
}
