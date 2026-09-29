package acp

import "testing"

func TestLinkStartsPending(t *testing.T) {
	links := NewLinkStore()
	link := links.Create("cus_1", "seti_1", "secret", 100000, "thb")
	if link.Status != LinkPending {
		t.Fatalf("expected pending, got %s", link.Status)
	}
}

func TestLinkActivateBySetupIntent(t *testing.T) {
	links := NewLinkStore()
	links.Create("cus_1", "seti_1", "secret", 100000, "thb")

	activated, ok := links.Activate("seti_1", "pm_1")
	if !ok {
		t.Fatal("expected link to be found and activated")
	}
	if activated.Status != LinkActive {
		t.Fatalf("expected active, got %s", activated.Status)
	}
	if activated.PaymentMethodID != "pm_1" {
		t.Fatalf("expected payment method to be recorded")
	}
}

func TestAuthorizeRejectsInactiveLink(t *testing.T) {
	links := NewLinkStore()
	link := links.Create("cus_1", "seti_1", "secret", 100000, "thb")

	if _, err := links.Authorize(link.ID, 5000); err == nil {
		t.Fatal("expected error for pending (not yet active) link")
	}
}

func TestAuthorizeRejectsAmountAboveMandate(t *testing.T) {
	links := NewLinkStore()
	links.Create("cus_1", "seti_1", "secret", 10000, "thb")
	link, _ := links.Activate("seti_1", "pm_1")

	if _, err := links.Authorize(link.ID, 20000); err == nil {
		t.Fatal("expected error for amount exceeding mandate")
	}
}

func TestAuthorizeAllowsAmountWithinMandate(t *testing.T) {
	links := NewLinkStore()
	links.Create("cus_1", "seti_1", "secret", 10000, "thb")
	link, _ := links.Activate("seti_1", "pm_1")

	if _, err := links.Authorize(link.ID, 5000); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAuthorizeRejectsUnknownLink(t *testing.T) {
	links := NewLinkStore()
	if _, err := links.Authorize("link_does_not_exist", 100); err == nil {
		t.Fatal("expected error for unknown link")
	}
}
