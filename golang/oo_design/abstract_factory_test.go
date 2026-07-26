package oodesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 * Abstract Factory pattern
 *
 * Provides an interface for creating families of related objects without
 * specifying their concrete types. Where Factory Method (factory_method_test.go)
 * produces a single product, an Abstract Factory produces a whole set of
 * products designed to be used together — and guarantees they match, since
 * one factory implementation is responsible for the entire family.
 */

// Button and Checkbox are the abstract products — every theme's family
// implements both.
type Button interface {
	Render() string
}

type Checkbox interface {
	Check() string
}

// Light theme family.
type LightButton struct{}

func (LightButton) Render() string { return "[light button]" }

type LightCheckbox struct{}

func (LightCheckbox) Check() string { return "[light checkbox checked]" }

// Dark theme family.
type DarkButton struct{}

func (DarkButton) Render() string { return "[dark button]" }

type DarkCheckbox struct{}

func (DarkCheckbox) Check() string { return "[dark checkbox checked]" }

// UIFactory is the abstract factory: one implementation produces every
// product in a theme's family, so a Button and Checkbox from the same
// factory are guaranteed to match.
type UIFactory interface {
	CreateButton() Button
	CreateCheckbox() Checkbox
}

type LightUIFactory struct{}

func (LightUIFactory) CreateButton() Button     { return LightButton{} }
func (LightUIFactory) CreateCheckbox() Checkbox { return LightCheckbox{} }

type DarkUIFactory struct{}

func (DarkUIFactory) CreateButton() Button     { return DarkButton{} }
func (DarkUIFactory) CreateCheckbox() Checkbox { return DarkCheckbox{} }

func Test_LightUIFactory_producesTheLightFamily(t *testing.T) {
	factory := LightUIFactory{}

	assert.Equal(t, "[light button]", factory.CreateButton().Render())
	assert.Equal(t, "[light checkbox checked]", factory.CreateCheckbox().Check())
}

func Test_DarkUIFactory_producesTheDarkFamily(t *testing.T) {
	factory := DarkUIFactory{}

	assert.Equal(t, "[dark button]", factory.CreateButton().Render())
	assert.Equal(t, "[dark checkbox checked]", factory.CreateCheckbox().Check())
}

func Test_UIFactory_callerNeverNamesAConcreteTheme(t *testing.T) {
	// renderForm only knows about the UIFactory interface — it builds a
	// matching button+checkbox pair without ever knowing which theme it got.
	renderForm := func(factory UIFactory) string {
		return factory.CreateButton().Render() + " " + factory.CreateCheckbox().Check()
	}

	assert.Equal(t, "[light button] [light checkbox checked]", renderForm(LightUIFactory{}))
	assert.Equal(t, "[dark button] [dark checkbox checked]", renderForm(DarkUIFactory{}))
}
