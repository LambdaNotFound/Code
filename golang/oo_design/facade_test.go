package oodesign

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 * Facade pattern
 *
 * Provides a single, simplified entry point in front of a complex subsystem
 * (inventory, payment, shipping) so callers don't need to know how those
 * pieces coordinate. The facade doesn't add new behavior — it just
 * orchestrates existing subsystem calls in the right order and stops at the
 * first failure.
 *
 * Note what it does NOT give you: sequencing isn't atomic. If Reserve
 * succeeds but Charge fails, the reservation is not automatically rolled
 * back — that requires a separate compensating-action (Saga) pattern.
 */

type inventoryService interface {
	Reserve(sku string, qty int) error
}

type paymentService interface {
	Charge(amountCents int64) error
}

type shippingService interface {
	Schedule(address string) error
}

// Inventory, Payment, Shipping are the independent subsystem implementations
// the facade coordinates.
type Inventory struct{ stock map[string]int }

func NewInventory(stock map[string]int) *Inventory {
	return &Inventory{stock: stock}
}

func (i *Inventory) Reserve(sku string, qty int) error {
	if i.stock[sku] < qty {
		return fmt.Errorf("insufficient stock for %s: have %d, want %d", sku, i.stock[sku], qty)
	}
	i.stock[sku] -= qty
	return nil
}

type Payment struct{ balanceCents int64 }

func NewPayment(balanceCents int64) *Payment {
	return &Payment{balanceCents: balanceCents}
}

func (p *Payment) Charge(amountCents int64) error {
	if p.balanceCents < amountCents {
		return fmt.Errorf("insufficient balance: have %d, want %d", p.balanceCents, amountCents)
	}
	p.balanceCents -= amountCents
	return nil
}

type Shipping struct{ scheduled []string }

func NewShipping() *Shipping {
	return &Shipping{}
}

func (s *Shipping) Schedule(address string) error {
	if address == "" {
		return fmt.Errorf("shipping address is empty")
	}
	s.scheduled = append(s.scheduled, address)
	return nil
}

// OrderFacade hides inventory/payment/shipping coordination behind one call.
type OrderFacade struct {
	inventory inventoryService
	payment   paymentService
	shipping  shippingService
}

func NewOrderFacade(inventory inventoryService, payment paymentService, shipping shippingService) *OrderFacade {
	return &OrderFacade{inventory: inventory, payment: payment, shipping: shipping}
}

func (f *OrderFacade) PlaceOrder(sku string, qty int, amountCents int64, address string) error {
	if err := f.inventory.Reserve(sku, qty); err != nil {
		return fmt.Errorf("reserve failed: %w", err)
	}
	if err := f.payment.Charge(amountCents); err != nil {
		return fmt.Errorf("charge failed: %w", err)
	}
	if err := f.shipping.Schedule(address); err != nil {
		return fmt.Errorf("shipping failed: %w", err)
	}
	return nil
}

func Test_OrderFacade_PlaceOrder(t *testing.T) {
	t.Run("success reserves, charges, and schedules", func(t *testing.T) {
		inventory := NewInventory(map[string]int{"widget": 5})
		payment := NewPayment(1000)
		shipping := NewShipping()
		facade := NewOrderFacade(inventory, payment, shipping)

		err := facade.PlaceOrder("widget", 2, 500, "1 Main St")

		assert.NoError(t, err)
		assert.Equal(t, 3, inventory.stock["widget"])
		assert.Equal(t, int64(500), payment.balanceCents)
		assert.Equal(t, []string{"1 Main St"}, shipping.scheduled)
	})

	t.Run("insufficient stock stops before payment and shipping", func(t *testing.T) {
		inventory := NewInventory(map[string]int{"widget": 1})
		payment := NewPayment(1000)
		shipping := NewShipping()
		facade := NewOrderFacade(inventory, payment, shipping)

		err := facade.PlaceOrder("widget", 2, 500, "1 Main St")

		assert.EqualError(t, err, "reserve failed: insufficient stock for widget: have 1, want 2")
		assert.Equal(t, int64(1000), payment.balanceCents) // never charged
		assert.Empty(t, shipping.scheduled)                // never scheduled
	})

	t.Run("insufficient balance stops before shipping", func(t *testing.T) {
		inventory := NewInventory(map[string]int{"widget": 5})
		payment := NewPayment(100)
		shipping := NewShipping()
		facade := NewOrderFacade(inventory, payment, shipping)

		err := facade.PlaceOrder("widget", 2, 500, "1 Main St")

		assert.EqualError(t, err, "charge failed: insufficient balance: have 100, want 500")
		assert.Equal(t, 3, inventory.stock["widget"]) // reservation is NOT rolled back
		assert.Empty(t, shipping.scheduled)
	})

	t.Run("shipping failure leaves earlier side effects in place", func(t *testing.T) {
		inventory := NewInventory(map[string]int{"widget": 5})
		payment := NewPayment(1000)
		shipping := NewShipping()
		facade := NewOrderFacade(inventory, payment, shipping)

		err := facade.PlaceOrder("widget", 2, 500, "")

		assert.EqualError(t, err, "shipping failed: shipping address is empty")
		assert.Equal(t, 3, inventory.stock["widget"])     // already reserved
		assert.Equal(t, int64(500), payment.balanceCents) // already charged
	})
}
