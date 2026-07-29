package oodesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 * Observer pattern
 *
 * Defines a one-to-many dependency between a Subject and any number of
 * Observers: whenever the Subject's state changes, every subscribed
 * Observer is notified automatically. The Subject never knows or cares what
 * a given Observer does with the notification, or how many there are.
 */

// Observer is notified whenever the Subject it's watching changes.
type Observer interface {
	Notify(price float64)
}

// Subject lets Observers subscribe and unsubscribe from its changes.
type Subject interface {
	Subscribe(o Observer)
	Unsubscribe(o Observer)
}

// StockTicker is a concrete Subject: it tracks a price and notifies every
// subscribed Observer whenever that price changes.
type StockTicker struct {
	symbol    string
	price     float64
	observers []Observer
}

func NewStockTicker(symbol string, price float64) *StockTicker {
	return &StockTicker{symbol: symbol, price: price}
}

func (t *StockTicker) Subscribe(o Observer) {
	t.observers = append(t.observers, o)
}

func (t *StockTicker) Unsubscribe(o Observer) {
	for i, existing := range t.observers {
		if existing == o {
			t.observers = append(t.observers[:i], t.observers[i+1:]...)
			return
		}
	}
}

func (t *StockTicker) SetPrice(price float64) {
	t.price = price
	for _, o := range t.observers {
		o.Notify(price)
	}
}

// PriceLogger is a concrete Observer that records every price it's notified of.
type PriceLogger struct {
	Prices []float64
}

func (l *PriceLogger) Notify(price float64) {
	l.Prices = append(l.Prices, price)
}

// PriceAlert is a second, unrelated Observer implementation — demonstrating
// that StockTicker works with any number of different Observer types at
// once, without ever needing to know which ones are subscribed.
type PriceAlert struct {
	Threshold float64
	Triggered bool
}

func (a *PriceAlert) Notify(price float64) {
	if price >= a.Threshold {
		a.Triggered = true
	}
}

func Test_StockTicker_notifiesAllSubscribedObservers(t *testing.T) {
	ticker := NewStockTicker("GOOG", 100)
	logger1 := &PriceLogger{}
	logger2 := &PriceLogger{}
	ticker.Subscribe(logger1)
	ticker.Subscribe(logger2)

	ticker.SetPrice(105)
	ticker.SetPrice(110)

	assert.Equal(t, []float64{105, 110}, logger1.Prices)
	assert.Equal(t, []float64{105, 110}, logger2.Prices)
}

func Test_StockTicker_unsubscribeStopsFurtherNotifications(t *testing.T) {
	ticker := NewStockTicker("GOOG", 100)
	logger := &PriceLogger{}
	ticker.Subscribe(logger)

	ticker.SetPrice(105)
	ticker.Unsubscribe(logger)
	ticker.SetPrice(110)

	assert.Equal(t, []float64{105}, logger.Prices)
}

func Test_StockTicker_supportsDifferentObserverTypesSimultaneously(t *testing.T) {
	ticker := NewStockTicker("GOOG", 100)
	logger := &PriceLogger{}
	alert := &PriceAlert{Threshold: 150}
	ticker.Subscribe(logger)
	ticker.Subscribe(alert)

	ticker.SetPrice(120)
	assert.False(t, alert.Triggered)

	ticker.SetPrice(160)
	assert.True(t, alert.Triggered)
	assert.Equal(t, []float64{120, 160}, logger.Prices)
}

func Test_StockTicker_noObserversIsSafe(t *testing.T) {
	ticker := NewStockTicker("GOOG", 100)

	assert.NotPanics(t, func() { ticker.SetPrice(200) })
}
