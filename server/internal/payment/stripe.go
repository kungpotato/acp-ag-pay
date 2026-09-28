// Package payment settles ACP checkout sessions through Stripe. This is the
// one place in the workshop that talks to a real external payment
// processor, always in Stripe test mode.
package payment

import (
	"errors"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
)

// ErrNotConfigured is returned when STRIPE_SECRET_KEY is unset. Callers
// should turn this into an HTTP 503 rather than a 500: the server is
// working correctly, a required credential is simply missing (see
// .env.example).
var ErrNotConfigured = errors.New("stripe is not configured: set STRIPE_SECRET_KEY")

// Settlement holds what the checkout session needs to remember about the
// Stripe side of a payment.
type Settlement struct {
	PaymentIntentID string
	ClientSecret    string
}

// StripeSettler creates PaymentIntents against the Stripe API. It is safe
// for concurrent use; the stripe-go SDK is stateless beyond the package
// -level API key.
type StripeSettler struct {
	configured bool
}

// NewStripeSettler configures the stripe-go SDK with the given secret key.
// An empty key produces a settler that always returns ErrNotConfigured,
// matching the graceful-degradation behavior promised in .env.example.
func NewStripeSettler(secretKey string) *StripeSettler {
	if secretKey == "" {
		return &StripeSettler{configured: false}
	}
	stripe.Key = secretKey
	return &StripeSettler{configured: true}
}

func (s *StripeSettler) Configured() bool { return s.configured }

// CreatePaymentIntent settles a checkout session's total by creating a
// Stripe PaymentIntent for that amount. This is the ACP "settle" step: the
// checkout session already carries a server-verified total (Lesson 4), and
// this function is the only place that total ever touches Stripe.
func (s *StripeSettler) CreatePaymentIntent(amountCents int64, currency string, sessionID string) (Settlement, error) {
	if !s.configured {
		return Settlement{}, ErrNotConfigured
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountCents),
		Currency: stripe.String(currency),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		Metadata: map[string]string{
			"acp_checkout_session_id": sessionID,
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return Settlement{}, err
	}

	return Settlement{
		PaymentIntentID: pi.ID,
		ClientSecret:    pi.ClientSecret,
	}, nil
}
