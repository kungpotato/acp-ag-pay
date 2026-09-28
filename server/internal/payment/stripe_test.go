package payment

import "testing"

func TestUnconfiguredSettlerReturnsErrNotConfigured(t *testing.T) {
	s := NewStripeSettler("")
	if s.Configured() {
		t.Fatal("expected settler to be unconfigured with empty key")
	}
	_, err := s.CreatePaymentIntent(1000, "thb", "cs_test")
	if err != ErrNotConfigured {
		t.Fatalf("expected ErrNotConfigured, got %v", err)
	}
}

func TestConfiguredSettlerReportsConfigured(t *testing.T) {
	s := NewStripeSettler("sk_test_dummy")
	if !s.Configured() {
		t.Fatal("expected settler to be configured with a non-empty key")
	}
}
