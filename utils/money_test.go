package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 * Money pattern
 *
 * Never model money with float64 — binary floating point can't represent
 * decimals like 0.1 exactly, so repeated addition silently drifts
 * (0.1 + 0.2 != 0.3). Instead store the amount as an integer count of the
 * currency's smallest unit (cents, pence, ...) and tag it with an ISO 4217
 * currency code. Arithmetic returns a new Money and refuses to mix
 * currencies.
 */

// Currency is an ISO 4217 alphabetic currency code.
type Currency string

const (
	USD Currency = "USD" // US Dollar — 2 minor-unit digits (1 unit = 100 cents)
	EUR Currency = "EUR" // Euro — 2 minor-unit digits (1 unit = 100 cents)
	JPY Currency = "JPY" // Japanese Yen — 0 minor-unit digits (no subunit in practice)
)

// Money is an amount in minor units (e.g. cents) of a given ISO 4217 currency.
type Money struct {
	amount   int64
	currency Currency
}

func NewMoney(amount int64, currency Currency) Money {
	return Money{amount: amount, currency: currency}
}

func (m Money) Amount() int64      { return m.amount }
func (m Money) Currency() Currency { return m.currency }

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("currency mismatch: %s vs %s", m.currency, other.currency)
	}
	return Money{amount: m.amount + other.amount, currency: m.currency}, nil
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("currency mismatch: %s vs %s", m.currency, other.currency)
	}
	return Money{amount: m.amount - other.amount, currency: m.currency}, nil
}

func (m Money) Multiply(factor int64) Money {
	return Money{amount: m.amount * factor, currency: m.currency}
}

/*
 * Rounding a fractional-cent value (e.g. amount * tax rate) back to whole
 * cents needs an explicit, documented policy — the two common ones disagree
 * exactly on .5 ties:
 *
 *   - round-half-up:      ties always round away from zero (2.50 -> 3)
 *   - round-half-to-even: ties round to the nearest even cent (2.50 -> 2),
 *     a.k.a. "banker's rounding" — avoids a systematic upward bias when
 *     rounding many values, e.g. aggregating tax across many line items.
 *
 * Both take an amount in hundredths-of-a-cent and return whole cents.
 */

func RoundHalfUp(hundredthsOfCent int64) int64 {
	return (hundredthsOfCent + 50) / 100
}

func RoundHalfToEven(hundredthsOfCent int64) int64 {
	cents := hundredthsOfCent / 100
	remainder := hundredthsOfCent % 100

	switch {
	case remainder > 50:
		return cents + 1
	case remainder < 50:
		return cents
	default: // exact .5 tie
		if cents%2 == 0 {
			return cents
		}
		return cents + 1
	}
}

func Test_Money_Add_sameCurrency(t *testing.T) {
	a := NewMoney(1050, USD) // $10.50
	b := NewMoney(250, USD)  // $2.50

	sum, err := a.Add(b)

	assert.NoError(t, err)
	assert.Equal(t, int64(1300), sum.Amount())
	assert.Equal(t, USD, sum.Currency())
}

func Test_Money_Add_currencyMismatch(t *testing.T) {
	usd := NewMoney(1000, USD)
	eur := NewMoney(1000, EUR)

	_, err := usd.Add(eur)

	assert.EqualError(t, err, "currency mismatch: USD vs EUR")
}

func Test_Money_Subtract_sameCurrency(t *testing.T) {
	a := NewMoney(500, USD)
	b := NewMoney(199, USD)

	diff, err := a.Subtract(b)

	assert.NoError(t, err)
	assert.Equal(t, int64(301), diff.Amount())
}

func Test_Money_Multiply(t *testing.T) {
	price := NewMoney(299, USD) // $2.99
	total := price.Multiply(3)

	assert.Equal(t, int64(897), total.Amount())
}

func Test_RoundHalfUp_exactTie(t *testing.T) {
	assert.Equal(t, int64(3), RoundHalfUp(250)) // 2.50 cents -> 3, ties away from zero
	assert.Equal(t, int64(4), RoundHalfUp(350)) // 3.50 cents -> 4
}

func Test_RoundHalfToEven_exactTie(t *testing.T) {
	assert.Equal(t, int64(2), RoundHalfToEven(250)) // 2.50 cents -> 2 (2 is even)
	assert.Equal(t, int64(4), RoundHalfToEven(350)) // 3.50 cents -> 4 (4 is even)
}

func Test_Rounding_agreesOffTie(t *testing.T) {
	assert.Equal(t, RoundHalfUp(260), RoundHalfToEven(260)) // 2.60 -> 3, no tie involved
	assert.Equal(t, RoundHalfUp(240), RoundHalfToEven(240)) // 2.40 -> 2, no tie involved
}
