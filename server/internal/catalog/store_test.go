package catalog

import "testing"

func TestSearchMatchesTitleCaseInsensitive(t *testing.T) {
	s := NewStore()
	got := s.Search("clean code")
	if len(got) != 1 || got[0].ID != "bk-003" {
		t.Fatalf("expected to find bk-003, got %+v", got)
	}
}

func TestSearchEmptyQueryReturnsAll(t *testing.T) {
	s := NewStore()
	if got := s.Search(""); len(got) != len(s.List()) {
		t.Fatalf("expected all books, got %d", len(got))
	}
}

func TestReserveStockFailsWhenInsufficient(t *testing.T) {
	s := NewStore()
	if err := s.ReserveStock("bk-003", 100); err == nil {
		t.Fatal("expected error for insufficient stock")
	}
}

func TestReserveStockSucceeds(t *testing.T) {
	s := NewStore()
	before, _ := s.Get("bk-001")
	if err := s.ReserveStock("bk-001", 2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	after, _ := s.Get("bk-001")
	if after.Stock != before.Stock-2 {
		t.Fatalf("expected stock to drop by 2, got %d -> %d", before.Stock, after.Stock)
	}
}
