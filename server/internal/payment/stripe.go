// Package payment settles ACP checkout sessions through Stripe. This is the
// one place in the workshop that talks to a real external payment
// processor, always in Stripe test mode.
package payment

import (
	"errors"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/customer"
	"github.com/stripe/stripe-go/v81/paymentintent"
	"github.com/stripe/stripe-go/v81/setupintent"
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

// LinkSetup holds what the shopper needs to complete one-time payment
// linking (Lesson 9): a Stripe Customer to attach the payment method to,
// and a SetupIntent client secret for the one-time consent UI.
type LinkSetup struct {
	CustomerID    string
	SetupIntentID string
	ClientSecret  string
}

// CreateSetupIntent starts the ONE-TIME linking flow that Lesson 9 builds
// on top of Lesson 5's per-order settle button. A shopper completes this
// once (e.g. via Stripe Link's own authentication), after which the agent
// holds a durable payment method it can charge off_session, with no further
// checkout UI. This function only creates the Customer + SetupIntent shell;
// the payment method is attached when the shopper confirms client-side, and
// we only trust that a mandate is active once the webhook in Lesson 7/9
// reports `setup_intent.succeeded`.
func (s *StripeSettler) CreateSetupIntent() (LinkSetup, error) {
	if !s.configured {
		return LinkSetup{}, ErrNotConfigured
	}

	cust, err := customer.New(&stripe.CustomerParams{})
	if err != nil {
		return LinkSetup{}, err
	}

	si, err := setupintent.New(&stripe.SetupIntentParams{
		Customer: stripe.String(cust.ID),
		AutomaticPaymentMethods: &stripe.SetupIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	})
	if err != nil {
		return LinkSetup{}, err
	}

	return LinkSetup{
		CustomerID:    cust.ID,
		SetupIntentID: si.ID,
		ClientSecret:  si.ClientSecret,
	}, nil
}

// CreateOffSessionPaymentIntent is the payment agentic commerce is actually
// about: the agent charges a previously linked payment method directly,
// with Confirm+OffSession set, and no shopper-facing checkout step at all.
// Contrast with CreatePaymentIntent (Lesson 5), which still needs a human to
// confirm on a client with the returned client secret.
func (s *StripeSettler) CreateOffSessionPaymentIntent(customerID, paymentMethodID string, amountCents int64, currency, sessionID string) (Settlement, error) {
	if !s.configured {
		return Settlement{}, ErrNotConfigured
	}

	pi, err := paymentintent.New(&stripe.PaymentIntentParams{
		Amount:        stripe.Int64(amountCents),
		Currency:      stripe.String(currency),
		Customer:      stripe.String(customerID),
		PaymentMethod: stripe.String(paymentMethodID),
		Confirm:       stripe.Bool(true),
		OffSession:    stripe.Bool(true),
		Metadata: map[string]string{
			"acp_checkout_session_id": sessionID,
		},
	})
	if err != nil {
		return Settlement{}, err
	}

	return Settlement{
		PaymentIntentID: pi.ID,
		ClientSecret:    pi.ClientSecret,
	}, nil
}
