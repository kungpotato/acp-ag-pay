package main

import (
	"log"
	"net/http"
	"os"

	"acp-ag-pay/server/internal/catalog"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	store := catalog.NewStore()
	mux := http.NewServeMux()
	catalog.NewHandler(store).Register(mux)

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	log.Printf("bookstore server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
