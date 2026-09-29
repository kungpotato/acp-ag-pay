package webhook

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stripe/stripe-go/v81"
	stripewebhook "github.com/stripe/stripe-go/v81/webhook"

	"acp-ag-pay/server/internal/acp"
)

const testSecret = "whsec_test_secret"

func signedRequest(t *testing.T, eventJSON string) *http.Request {
	t.Helper()
	signed := stripewebhook.GenerateTestSignedPayload(&stripewebhook.UnsignedPayload{
		Payload:   []byte(eventJSON),
		Secret:    testSecret,
		Timestamp: time.Now(),
	})
	req := httptest.NewRequest(http.MethodPost, "/api/webhooks/stripe", strings.NewReader(string(signed.Payload)))
	req.Header.Set("Stripe-Signature", signed.Header)
	return req
}

func paymentIntentEventJSON(eventType, paymentIntentID string) string {
	return `{
		"id": "evt_test",
		"object": "event",
		"api_version": "` + stripe.APIVersion + `",
		"type": "` + eventType + `",
		"data": {
			"object": {
				"id": "` + paymentIntentID + `",
				"object": "payment_intent"
			}
		}
	}`
}

func TestWebhookReturns503WhenNotConfigured(t *testing.T) {
	orders := acp.NewOrderStore()
	h := NewHandler(orders, acp.NewLinkStore(), "")
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest(http.MethodPost, "/api/webhooks/stripe", strings.NewReader("{}"))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}
}

func TestWebhookRejectsInvalidSignature(t *testing.T) {
	orders := acp.NewOrderStore()
	h := NewHandler(orders, acp.NewLinkStore(), testSecret)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest(http.MethodPost, "/api/webhooks/stripe", strings.NewReader("{}"))
	req.Header.Set("Stripe-Signature", "t=1,v1=deadbeef")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestWebhookConfirmsOrderOnPaymentSucceeded(t *testing.T) {
	orders := acp.NewOrderStore()
	order := orders.CreateFromSession(acp.CheckoutSession{
		ID: "cs_test", PaymentIntentID: "pi_test_123", TotalCents: 1000, Currency: "thb",
	})

	h := NewHandler(orders, acp.NewLinkStore(), testSecret)
	mux := http.NewServeMux()
	h.Register(mux)

	req := signedRequest(t, paymentIntentEventJSON("payment_intent.succeeded", "pi_test_123"))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	updated, ok := orders.Get(order.ID)
	if !ok {
		t.Fatal("expected order to still exist")
	}
	if updated.Status != acp.OrderPaid {
		t.Fatalf("expected order to be paid, got %s", updated.Status)
	}
}

func setupIntentEventJSON(setupIntentID, paymentMethodID string) string {
	return `{
		"id": "evt_test",
		"object": "event",
		"api_version": "` + stripe.APIVersion + `",
		"type": "setup_intent.succeeded",
		"data": {
			"object": {
				"id": "` + setupIntentID + `",
				"object": "setup_intent",
				"payment_method": "` + paymentMethodID + `"
			}
		}
	}`
}

func TestWebhookActivatesLinkOnSetupIntentSucceeded(t *testing.T) {
	orders := acp.NewOrderStore()
	links := acp.NewLinkStore()
	link := links.Create("cus_test", "seti_test_1", "secret", 100000, "thb")

	h := NewHandler(orders, links, testSecret)
	mux := http.NewServeMux()
	h.Register(mux)

	req := signedRequest(t, setupIntentEventJSON("seti_test_1", "pm_test_1"))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	updated, ok := links.Get(link.ID)
	if !ok {
		t.Fatal("expected link to still exist")
	}
	if updated.Status != acp.LinkActive {
		t.Fatalf("expected link to be active, got %s", updated.Status)
	}
	if updated.PaymentMethodID != "pm_test_1" {
		t.Fatalf("expected payment method to be recorded, got %q", updated.PaymentMethodID)
	}
}

func TestWebhookMarksOrderFailedOnPaymentFailed(t *testing.T) {
	orders := acp.NewOrderStore()
	order := orders.CreateFromSession(acp.CheckoutSession{
		ID: "cs_test", PaymentIntentID: "pi_test_456",
	})

	h := NewHandler(orders, acp.NewLinkStore(), testSecret)
	mux := http.NewServeMux()
	h.Register(mux)

	req := signedRequest(t, paymentIntentEventJSON("payment_intent.payment_failed", "pi_test_456"))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	updated, _ := orders.Get(order.ID)
	if updated.Status != acp.OrderFailed {
		t.Fatalf("expected order to be failed, got %s", updated.Status)
	}
}
