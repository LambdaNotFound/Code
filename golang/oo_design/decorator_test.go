package oodesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 * Decorator pattern
 *
 * Wraps an object in another object that implements the same interface,
 * adding behavior before/after delegating to the wrapped object. Unlike
 * subclassing, decorators compose at runtime and stack in any order and any
 * combination, without an explosion of subclasses for every combination
 * (e.g. MilkSugarEspresso, SugarMilkEspresso, ...).
 */

// Beverage is the component interface both the base drink and every
// decorator implement.
type Beverage interface {
	Description() string
	CostCents() int64
}

// Espresso is the concrete component being decorated.
type Espresso struct{}

func (Espresso) Description() string { return "Espresso" }
func (Espresso) CostCents() int64    { return 250 }

// MilkDecorator, SugarDecorator, WhipDecorator each wrap a Beverage and add
// their own cost/description on top of whatever they wrap.
type MilkDecorator struct{ Beverage Beverage }

func (d MilkDecorator) Description() string { return d.Beverage.Description() + " + Milk" }
func (d MilkDecorator) CostCents() int64    { return d.Beverage.CostCents() + 50 }

type SugarDecorator struct{ Beverage Beverage }

func (d SugarDecorator) Description() string { return d.Beverage.Description() + " + Sugar" }
func (d SugarDecorator) CostCents() int64    { return d.Beverage.CostCents() + 25 }

type WhipDecorator struct{ Beverage Beverage }

func (d WhipDecorator) Description() string { return d.Beverage.Description() + " + Whip" }
func (d WhipDecorator) CostCents() int64    { return d.Beverage.CostCents() + 75 }

func Test_Espresso_alone(t *testing.T) {
	var drink Beverage = Espresso{}

	assert.Equal(t, "Espresso", drink.Description())
	assert.Equal(t, int64(250), drink.CostCents())
}

func Test_Decorator_singleWrap(t *testing.T) {
	var drink Beverage = MilkDecorator{Beverage: Espresso{}}

	assert.Equal(t, "Espresso + Milk", drink.Description())
	assert.Equal(t, int64(300), drink.CostCents())
}

func Test_Decorator_stackedInAnyOrder(t *testing.T) {
	t.Run("milk then sugar", func(t *testing.T) {
		var drink Beverage = SugarDecorator{Beverage: MilkDecorator{Beverage: Espresso{}}}

		assert.Equal(t, "Espresso + Milk + Sugar", drink.Description())
		assert.Equal(t, int64(325), drink.CostCents())
	})

	t.Run("sugar then milk - same total cost, different description order", func(t *testing.T) {
		var drink Beverage = MilkDecorator{Beverage: SugarDecorator{Beverage: Espresso{}}}

		assert.Equal(t, "Espresso + Sugar + Milk", drink.Description())
		assert.Equal(t, int64(325), drink.CostCents())
	})

	t.Run("three decorators stacked", func(t *testing.T) {
		var drink Beverage = WhipDecorator{Beverage: SugarDecorator{Beverage: MilkDecorator{Beverage: Espresso{}}}}

		assert.Equal(t, "Espresso + Milk + Sugar + Whip", drink.Description())
		assert.Equal(t, int64(400), drink.CostCents())
	})
}

func Test_Decorator_sameDecoratorAppliedTwice(t *testing.T) {
	// decorators compose freely — nothing stops wrapping with the same one
	// more than once, e.g. double milk.
	var drink Beverage = MilkDecorator{Beverage: MilkDecorator{Beverage: Espresso{}}}

	assert.Equal(t, "Espresso + Milk + Milk", drink.Description())
	assert.Equal(t, int64(350), drink.CostCents())
}
