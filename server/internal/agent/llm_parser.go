package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// LLMParser upgrades intent parsing to an actual language model via
// OpenRouter. It is entirely optional: the shopping agent falls back to
// RuleParser automatically when no API key is configured (see NewParser).
type LLMParser struct {
	apiKey string
	model  string
	client *http.Client
}

func NewLLMParser(apiKey string) *LLMParser {
	return &LLMParser{
		apiKey: apiKey,
		model:  "openai/gpt-4o-mini",
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

const systemPrompt = `You are an intent parser for a bookstore shopping agent.
Classify the user's Thai or English message into JSON with fields:
"kind" (one of "search_book", "greet", "unknown") and "query" (the book,
author, or genre being searched for, empty string if not applicable).
Respond with ONLY the JSON object, no prose.`

type openRouterRequest struct {
	Model    string              `json:"model"`
	Messages []openRouterMessage `json:"messages"`
}

type openRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openRouterResponse struct {
	Choices []struct {
		Message openRouterMessage `json:"message"`
	} `json:"choices"`
}

func (p *LLMParser) Parse(ctx context.Context, message string) (Intent, error) {
	body := openRouterRequest{
		Model: p.model,
		Messages: []openRouterMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: message},
		},
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return Intent{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return Intent{}, err
	}
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return Intent{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Intent{}, fmt.Errorf("openrouter returned status %d", resp.StatusCode)
	}

	var out openRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return Intent{}, err
	}
	if len(out.Choices) == 0 {
		return Intent{}, fmt.Errorf("openrouter returned no choices")
	}

	var intent Intent
	content := strings.TrimSpace(out.Choices[0].Message.Content)
	if err := json.Unmarshal([]byte(content), &intent); err != nil {
		return Intent{}, fmt.Errorf("failed to parse model output as intent JSON: %w", err)
	}
	return intent, nil
}
