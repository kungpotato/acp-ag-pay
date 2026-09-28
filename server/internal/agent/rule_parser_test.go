package agent

import "testing"

func TestRuleParserGreet(t *testing.T) {
	p := NewRuleParser()
	intent := p.Parse("สวัสดีค่ะ")
	if intent.Kind != IntentGreet {
		t.Fatalf("expected greet, got %+v", intent)
	}
}

func TestRuleParserSearchStripsFillerWords(t *testing.T) {
	p := NewRuleParser()
	intent := p.Parse("หาแฮร์รี่ พอตเตอร์ให้หน่อย")
	if intent.Kind != IntentSearch {
		t.Fatalf("expected search, got %+v", intent)
	}
	if intent.Query == "" {
		t.Fatal("expected non-empty query")
	}
}

func TestRuleParserEnglishSearch(t *testing.T) {
	p := NewRuleParser()
	intent := p.Parse("find me clean code")
	if intent.Kind != IntentSearch {
		t.Fatalf("expected search, got %+v", intent)
	}
}

func TestRuleParserUnknownOnEmptyAfterStripping(t *testing.T) {
	p := NewRuleParser()
	intent := p.Parse("หา")
	if intent.Kind != IntentUnknown {
		t.Fatalf("expected unknown, got %+v", intent)
	}
}
