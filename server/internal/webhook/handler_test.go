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
	h := NewHandler(orders, "")
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
	h := NewHandler(orders, testSecret)
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

	h := NewHandler(orders, testSecret)
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

func TestWebhookMarksOrderFailedOnPaymentFailed(t *testing.T) {
	orders := acp.NewOrderStore()
	order := orders.CreateFromSession(acp.CheckoutSession{
		ID: "cs_test", PaymentIntentID: "pi_test_456",
	})

	h := NewHandler(orders, testSecret)
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
