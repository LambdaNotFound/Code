package oodesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 * Factory Method pattern
 *
 * Defines a method for creating an object, but lets each concrete type
 * decide which product to instantiate — the caller works only against the
 * creator's interface and never names a concrete product type directly.
 *
 * Contrast with Abstract Factory (abstract_factory_test.go): Factory Method
 * produces a single product per creator; Abstract Factory produces a whole
 * family of related products guaranteed to match each other.
 */

// Transport is the product interface — what every concrete transport can do.
type Transport interface {
	Deliver() string
}

type Truck struct{}

func (Truck) Deliver() string { return "delivering by truck on the road" }

type Ship struct{}

func (Ship) Deliver() string { return "delivering by ship on the sea" }

// Logistics declares the factory method (CreateTransport) plus a method
// that uses it — PlanDelivery never names Truck or Ship, so a new Logistics
// variant can introduce a new Transport without any change to this code.
type Logistics interface {
	CreateTransport() Transport
	PlanDelivery() string
}

type RoadLogistics struct{}

func (RoadLogistics) CreateTransport() Transport { return Truck{} }

func (l RoadLogistics) PlanDelivery() string {
	return l.CreateTransport().Deliver()
}

type SeaLogistics struct{}

func (SeaLogistics) CreateTransport() Transport { return Ship{} }

func (l SeaLogistics) PlanDelivery() string {
	return l.CreateTransport().Deliver()
}

func Test_RoadLogistics_CreateTransport(t *testing.T) {
	var logistics Logistics = RoadLogistics{}

	transport := logistics.CreateTransport()

	assert.IsType(t, Truck{}, transport)
	assert.Equal(t, "delivering by truck on the road", transport.Deliver())
}

func Test_SeaLogistics_CreateTransport(t *testing.T) {
	var logistics Logistics = SeaLogistics{}

	transport := logistics.CreateTransport()

	assert.IsType(t, Ship{}, transport)
	assert.Equal(t, "delivering by ship on the sea", transport.Deliver())
}

func Test_PlanDelivery_usesTheFactoryMethod(t *testing.T) {
	assert.Equal(t, "delivering by truck on the road", RoadLogistics{}.PlanDelivery())
	assert.Equal(t, "delivering by ship on the sea", SeaLogistics{}.PlanDelivery())
}

func Test_Logistics_callerNeverNamesAConcreteTransport(t *testing.T) {
	// planRoute only knows about the Logistics interface — swapping which
	// concrete Logistics it receives changes the whole delivery mechanism,
	// with no change to planRoute itself.
	planRoute := func(l Logistics) string {
		return l.PlanDelivery()
	}

	assert.Equal(t, "delivering by truck on the road", planRoute(RoadLogistics{}))
	assert.Equal(t, "delivering by ship on the sea", planRoute(SeaLogistics{}))
}
