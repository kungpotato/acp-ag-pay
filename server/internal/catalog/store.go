package catalog

import (
	"fmt"
	"strings"
	"sync"
)

// Store is an in-memory book catalog. A real bookstore would back this with
// a database, but an in-memory map keeps the workshop focused on the ACP
// and Stripe mechanics rather than persistence.
type Store struct {
	mu    sync.RWMutex
	books map[string]Book
}

func NewStore() *Store {
	s := &Store{books: make(map[string]Book)}
	for _, b := range seedBooks() {
		s.books[b.ID] = b
	}
	return s
}

func (s *Store) List() []Book {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Book, 0, len(s.books))
	for _, b := range s.books {
		out = append(out, b)
	}
	return out
}

func (s *Store) Get(id string) (Book, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.books[id]
	return b, ok
}

// Search does a simple case-insensitive substring match on title/author/genre.
// Lesson 2 builds the shopping agent on top of this.
func (s *Store) Search(query string) []Book {
	s.mu.RLock()
	defer s.mu.RUnlock()
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return s.List()
	}
	var out []Book
	for _, b := range s.books {
		hay := strings.ToLower(b.Title + " " + b.Author + " " + b.Genre)
		if strings.Contains(hay, q) {
			out = append(out, b)
		}
	}
	return out
}

// ReserveStock decrements stock for a checkout. Returns an error if not enough
// stock is available.
func (s *Store) ReserveStock(id string, qty int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.books[id]
	if !ok {
		return fmt.Errorf("book %s not found", id)
	}
	if b.Stock < qty {
		return fmt.Errorf("insufficient stock for %s: have %d, want %d", id, b.Stock, qty)
	}
	b.Stock -= qty
	s.books[id] = b
	return nil
}

func seedBooks() []Book {
	return []Book{
		{ID: "bk-001", Title: "สามก๊ก", Author: "หลอกว้านจง", Genre: "วรรณกรรมคลาสสิก", PriceCents: 29000, Currency: "thb", Description: "มหากาพย์การเมืองและสงครามยุคสามก๊ก", Stock: 12},
		{ID: "bk-002", Title: "The Go Programming Language", Author: "Alan Donovan & Brian Kernighan", Genre: "เทคโนโลยี", PriceCents: 89000, Currency: "thb", Description: "หนังสือเรียน Go ฉบับมาตรฐาน", Stock: 8},
		{ID: "bk-003", Title: "Clean Code", Author: "Robert C. Martin", Genre: "เทคโนโลยี", PriceCents: 75000, Currency: "thb", Description: "แนวทางเขียนโค้ดให้อ่านง่ายและดูแลง่าย", Stock: 5},
		{ID: "bk-004", Title: "แฮร์รี่ พอตเตอร์กับศิลาอาถรรพ์", Author: "J.K. Rowling", Genre: "แฟนตาซี", PriceCents: 32500, Currency: "thb", Description: "จุดเริ่มต้นของพ่อมดน้อยแฮร์รี่ พอตเตอร์", Stock: 20},
		{ID: "bk-005", Title: "Sapiens: A Brief History of Humankind", Author: "Yuval Noah Harari", Genre: "สารคดี", PriceCents: 45000, Currency: "thb", Description: "ประวัติศาสตร์มนุษยชาติโดยสังเขป", Stock: 15},
		{ID: "bk-006", Title: "Accelerate", Author: "Nicole Forsgren, Jez Humble, Gene Kim", Genre: "เทคโนโลยี", PriceCents: 68000, Currency: "thb", Description: "งานวิจัยเรื่อง high-performing software teams", Stock: 6},
	}
}
