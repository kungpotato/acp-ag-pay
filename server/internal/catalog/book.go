package catalog

// Book is a single title sold in the bookstore.
type Book struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Genre       string `json:"genre"`
	PriceCents  int64  `json:"price_cents"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
	Stock       int    `json:"stock"`
}
