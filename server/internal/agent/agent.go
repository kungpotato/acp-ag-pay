package agent

import (
	"context"
	"fmt"

	"acp-ag-pay/server/internal/catalog"
)

// Reply is what the shopping agent sends back to the chat UI.
type Reply struct {
	Message string         `json:"message"`
	Books   []catalog.Book `json:"books,omitempty"`
	Intent  Intent         `json:"intent"`
}

// Agent turns a free-text shopper message into a catalog search and a
// natural-language reply. It prefers the LLM parser when one is configured,
// and always has the rule-based parser as a zero-dependency fallback.
type Agent struct {
	store *catalog.Store
	rules *RuleParser
	llm   *LLMParser
}

func New(store *catalog.Store, llm *LLMParser) *Agent {
	return &Agent{store: store, rules: NewRuleParser(), llm: llm}
}

func (a *Agent) Handle(ctx context.Context, message string) Reply {
	intent := a.resolveIntent(ctx, message)

	switch intent.Kind {
	case IntentGreet:
		return Reply{
			Message: "สวัสดีค่ะ ยินดีต้อนรับสู่ร้านขายหนังสือ อยากได้หนังสือแนวไหนบอกได้เลยค่ะ",
			Intent:  intent,
		}
	case IntentSearch:
		books := a.store.Search(intent.Query)
		if len(books) == 0 {
			return Reply{
				Message: fmt.Sprintf("ขออภัยค่ะ ไม่พบหนังสือที่ตรงกับ \"%s\"", intent.Query),
				Intent:  intent,
			}
		}
		return Reply{
			Message: fmt.Sprintf("พบหนังสือ %d เล่มที่เกี่ยวกับ \"%s\"", len(books), intent.Query),
			Books:   books,
			Intent:  intent,
		}
	default:
		return Reply{
			Message: "ลองบอกชื่อหนังสือ ผู้เขียน หรือแนวที่สนใจได้เลยค่ะ เช่น \"หาแฮร์รี่ พอตเตอร์\"",
			Intent:  intent,
		}
	}
}

// resolveIntent tries the LLM parser first (when configured) and falls back
// to the deterministic rule-based parser on any error or when no LLM is
// configured at all. This mirrors the guarantee in .env.example: the agent
// never breaks just because OPENROUTER_API_KEY is unset.
func (a *Agent) resolveIntent(ctx context.Context, message string) Intent {
	if a.llm != nil {
		if intent, err := a.llm.Parse(ctx, message); err == nil {
			return intent
		}
	}
	return a.rules.Parse(message)
}
