package agent

import (
	"context"
	"testing"

	"acp-ag-pay/server/internal/catalog"
)

func TestAgentHandleSearchFindsBooks(t *testing.T) {
	a := New(catalog.NewStore(), nil)
	reply := a.Handle(context.Background(), "หาแฮร์รี่ พอตเตอร์ให้หน่อย")
	if len(reply.Books) == 0 {
		t.Fatalf("expected books in reply, got %+v", reply)
	}
}

func TestAgentHandleGreet(t *testing.T) {
	a := New(catalog.NewStore(), nil)
	reply := a.Handle(context.Background(), "สวัสดี")
	if reply.Intent.Kind != IntentGreet {
		t.Fatalf("expected greet intent, got %+v", reply.Intent)
	}
}

func TestAgentHandleNoMatchesIsGraceful(t *testing.T) {
	a := New(catalog.NewStore(), nil)
	reply := a.Handle(context.Background(), "หาหนังสือเกี่ยวกับควอนตัมฟิสิกส์ระดับปริญญาเอก")
	if len(reply.Books) != 0 {
		t.Fatalf("expected no books, got %+v", reply.Books)
	}
	if reply.Message == "" {
		t.Fatal("expected a graceful message even with no matches")
	}
}
