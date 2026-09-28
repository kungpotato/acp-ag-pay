package acp

import "testing"

func TestOrderCreateFromSessionStartsPendingConfirmation(t *testing.T) {
	orders := NewOrderStore()
	session := CheckoutSession{
		ID: "cs_test", TotalCents: 1000, Currency: "thb",
		PaymentIntentID: "pi_test",
	}
	order := orders.CreateFromSession(session)
	if order.Status != OrderPendingConfirmation {
		t.Fatalf("expected pending_confirmation, got %s", order.Status)
	}
	if order.SessionID != session.ID {
		t.Fatalf("expected session id to carry over")
	}
}

func TestOrderGetByPaymentIntent(t *testing.T) {
	orders := NewOrderStore()
	session := CheckoutSession{ID: "cs_test", PaymentIntentID: "pi_abc"}
	created := orders.CreateFromSession(session)

	got, ok := orders.GetByPaymentIntent("pi_abc")
	if !ok || got.ID != created.ID {
		t.Fatalf("expected to find order by payment intent id")
	}
}

func TestOrderConfirmSetsStatusAndTimestamp(t *testing.T) {
	orders := NewOrderStore()
	created := orders.CreateFromSession(CheckoutSession{ID: "cs_test", PaymentIntentID: "pi_x"})

	confirmed, ok := orders.Confirm(created.ID, OrderPaid)
	if !ok {
		t.Fatal("expected order to be found")
	}
	if confirmed.Status != OrderPaid {
		t.Fatalf("expected paid, got %s", confirmed.Status)
	}
	if confirmed.ConfirmedAt == nil {
		t.Fatal("expected ConfirmedAt to be set")
	}
}
