package acp

import (
	"testing"

	"acp-ag-pay/server/internal/catalog"
)

func TestFeedMarksOutOfStock(t *testing.T) {
	store := catalog.NewStore()
	books := store.List()
	if len(books) == 0 {
		t.Fatal("expected seed books")
	}
	// Drain stock for the first book.
	first := books[0]
	if err := store.ReserveStock(first.ID, first.Stock); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	feed := Feed(store)
	var found bool
	for _, p := range feed {
		if p.ID == first.ID {
			found = true
			if p.Availability != "out_of_stock" {
				t.Fatalf("expected out_of_stock, got %s", p.Availability)
			}
		}
	}
	if !found {
		t.Fatal("expected product in feed")
	}
}

func TestFeedIncludesAllBooks(t *testing.T) {
	store := catalog.NewStore()
	if len(Feed(store)) != len(store.List()) {
		t.Fatal("expected feed to include every book")
	}
}
