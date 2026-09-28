// Package acp implements a simplified, educational subset of the Agentic
// Commerce Protocol (ACP): a product feed for discovery (this file) and a
// checkout session lifecycle (session.go). It is not a certified ACP
// implementation — field names and flow are trimmed down to the concepts a
// workshop needs to build real intuition: how a shopping agent *discovers*
// products it never hardcoded, and how a *session* carries a cart through
// checkout to settlement.
package acp

import "acp-ag-pay/server/internal/catalog"

// Product is one entry in the ACP product feed. An external shopping agent
// (or our own, from Lesson 2) discovers the whole catalog through this feed
// without needing bespoke integration code per store.
type Product struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	PriceCents   int64  `json:"price_cents"`
	Currency     string `json:"currency"`
	Availability string `json:"availability"` // "in_stock" | "out_of_stock"
	URL          string `json:"url"`
}

func productFromBook(b catalog.Book) Product {
	availability := "in_stock"
	if b.Stock <= 0 {
		availability = "out_of_stock"
	}
	return Product{
		ID:           b.ID,
		Title:        b.Title,
		Description:  b.Description,
		PriceCents:   b.PriceCents,
		Currency:     b.Currency,
		Availability: availability,
		URL:          "/api/books/" + b.ID,
	}
}

// Feed converts the full catalog into an ACP-style product feed.
func Feed(store *catalog.Store) []Product {
	books := store.List()
	feed := make([]Product, 0, len(books))
	for _, b := range books {
		feed = append(feed, productFromBook(b))
	}
	return feed
}
