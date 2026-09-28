package main

import (
	"log"
	"net/http"
	"os"

	"acp-ag-pay/server/internal/agent"
	"acp-ag-pay/server/internal/catalog"
	"acp-ag-pay/server/internal/httpx"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	store := catalog.NewStore()
	mux := http.NewServeMux()
	catalog.NewHandler(store).Register(mux)

	var llmParser *agent.LLMParser
	if key := os.Getenv("OPENROUTER_API_KEY"); key != "" {
		llmParser = agent.NewLLMParser(key)
		log.Print("shopping agent: LLM-backed intent parsing enabled (OpenRouter)")
	} else {
		log.Print("shopping agent: OPENROUTER_API_KEY not set, using rule-based parser")
	}
	shoppingAgent := agent.New(store, llmParser)
	agent.NewHandler(shoppingAgent).Register(mux)

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	log.Printf("bookstore server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, httpx.CORS(mux)); err != nil {
		log.Fatal(err)
	}
}
