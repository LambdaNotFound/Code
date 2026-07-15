package oodesign

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 * Adapter pattern
 *
 * Converts the interface of an existing type into another interface that
 * client code expects, so incompatible types can work together without
 * modifying either side. Unlike Facade (which simplifies a whole
 * subsystem), Adapter's job is narrow: translate one call shape into
 * another — same behavior, different signature.
 */

// PaymentProcessor is the interface application code depends on.
type PaymentProcessor interface {
	Pay(amountCents int64) error
}

// modernGateway already satisfies PaymentProcessor natively — no adapter needed.
type modernGateway struct{ charged int64 }

func (g *modernGateway) Pay(amountCents int64) error {
	g.charged += amountCents
	return nil
}

// LegacyGateway is a pre-existing/third-party type with an incompatible
// signature: dollars instead of cents, and it returns a confirmation code
// instead of a plain error.
type LegacyGateway struct{ ChargedDollars float64 }

func (g *LegacyGateway) SubmitPaymentInDollars(dollars float64) (confirmationCode string, err error) {
	if dollars < 0 {
		return "", fmt.Errorf("cannot charge negative amount: %.2f", dollars)
	}
	g.ChargedDollars += dollars
	return fmt.Sprintf("LEGACY-%.0f", dollars*100), nil
}

// LegacyGatewayAdapter adapts LegacyGateway to the PaymentProcessor
// interface: converts cents to dollars and drops the confirmation code.
type LegacyGatewayAdapter struct{ legacy *LegacyGateway }

func NewLegacyGatewayAdapter(legacy *LegacyGateway) *LegacyGatewayAdapter {
	return &LegacyGatewayAdapter{legacy: legacy}
}

func (a *LegacyGatewayAdapter) Pay(amountCents int64) error {
	dollars := float64(amountCents) / 100
	_, err := a.legacy.SubmitPaymentInDollars(dollars)
	return err
}

// Checkout depends only on PaymentProcessor — it doesn't know or care
// whether it's talking to a modern gateway or an adapted legacy one.
func Checkout(processor PaymentProcessor, amountCents int64) error {
	return processor.Pay(amountCents)
}

func Test_LegacyGatewayAdapter_Pay(t *testing.T) {
	t.Run("converts cents to dollars for the legacy API", func(t *testing.T) {
		legacy := &LegacyGateway{}
		adapter := NewLegacyGatewayAdapter(legacy)

		err := adapter.Pay(2599) // $25.99

		assert.NoError(t, err)
		assert.InDelta(t, 25.99, legacy.ChargedDollars, 0.001)
	})

	t.Run("propagates errors from the legacy gateway", func(t *testing.T) {
		legacy := &LegacyGateway{}
		adapter := NewLegacyGatewayAdapter(legacy)

		err := adapter.Pay(-100)

		assert.EqualError(t, err, "cannot charge negative amount: -1.00")
	})
}

func Test_Checkout_acceptsEitherImplementation(t *testing.T) {
	t.Run("modern gateway satisfies PaymentProcessor directly", func(t *testing.T) {
		modern := &modernGateway{}

		err := Checkout(modern, 500)

		assert.NoError(t, err)
		assert.Equal(t, int64(500), modern.charged)
	})

	t.Run("legacy gateway works through the adapter", func(t *testing.T) {
		legacy := &LegacyGateway{}
		adapter := NewLegacyGatewayAdapter(legacy)

		err := Checkout(adapter, 500)

		assert.NoError(t, err)
		assert.InDelta(t, 5.00, legacy.ChargedDollars, 0.001)
	})
}
