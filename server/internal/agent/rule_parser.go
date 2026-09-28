package agent

import (
	"strings"
)

// RuleParser is a deterministic, dependency-free intent parser. It never
// calls an external API, so it always works even without OPENROUTER_API_KEY
// (see LLMParser for the optional upgrade).
type RuleParser struct{}

func NewRuleParser() *RuleParser { return &RuleParser{} }

var greetings = []string{"สวัสดี", "หวัดดี", "hello", "hi", "hey"}

// fillerPhrases are stripped from the message before treating the remainder
// as a search query. Order matters: longer phrases first so a short phrase
// doesn't leave a dangling fragment of a longer one.
var fillerPhrases = []string{
	"อยากได้หนังสือเกี่ยวกับ",
	"หาหนังสือเกี่ยวกับ",
	"ขอหนังสือเกี่ยวกับ",
	"อยากได้หนังสือ",
	"หาหนังสือ",
	"ขอหนังสือ",
	"อยากได้",
	"ช่วยหา",
	"หาให้หน่อย",
	"ให้หน่อย",
	"หน่อย",
	"หา",
	"ขอ",
	"looking for",
	"search for",
	"find me",
	"i want",
	"search",
	"find",
}

func (p *RuleParser) Parse(message string) Intent {
	trimmed := strings.TrimSpace(message)
	lower := strings.ToLower(trimmed)

	for _, g := range greetings {
		if lower == g || strings.Contains(lower, g) {
			return Intent{Kind: IntentGreet}
		}
	}

	query := lower
	for _, phrase := range fillerPhrases {
		query = strings.ReplaceAll(query, phrase, "")
	}
	query = strings.TrimSpace(strings.Trim(query, "?!."))

	if query == "" {
		return Intent{Kind: IntentUnknown}
	}
	return Intent{Kind: IntentSearch, Query: query}
}
