package agent

// IntentKind classifies what the shopper is trying to do.
type IntentKind string

const (
	IntentSearch  IntentKind = "search_book"
	IntentGreet   IntentKind = "greet"
	IntentUnknown IntentKind = "unknown"
)

// Intent is the structured meaning extracted from a free-text user message.
type Intent struct {
	Kind  IntentKind `json:"kind"`
	Query string     `json:"query,omitempty"`
}
