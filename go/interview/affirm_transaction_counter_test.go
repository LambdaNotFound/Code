package interview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Counter(t *testing.T) {
	t.Run("empty_counter_returns_zero", func(t *testing.T) {
		c := NewCounter()
		assert.Equal(t, 0.0, c.getTotalTransactionInLastOneHour(100))
	})

	t.Run("single_transaction_within_window", func(t *testing.T) {
		c := NewCounter()
		c.putTransaction(50.0, 1000)
		assert.Equal(t, 50.0, c.getTotalTransactionInLastOneHour(1001))
	})

	t.Run("transaction_outside_window_excluded", func(t *testing.T) {
		c := NewCounter()
		c.putTransaction(100.0, 0)
		assert.Equal(t, 0.0, c.getTotalTransactionInLastOneHour(3601))
	})

	t.Run("multiple_transactions_summed", func(t *testing.T) {
		c := NewCounter()
		c.putTransaction(10.0, 100)
		c.putTransaction(20.0, 200)
		c.putTransaction(30.0, 300)
		assert.Equal(t, 60.0, c.getTotalTransactionInLastOneHour(300))
	})

	t.Run("old_transactions_evicted", func(t *testing.T) {
		c := NewCounter()
		c.putTransaction(100.0, 1000)
		c.putTransaction(50.0, 5000)
		// at t=4601, t=1000 is > 3600s ago → evicted
		assert.Equal(t, 50.0, c.getTotalTransactionInLastOneHour(4601))
	})

	t.Run("same_timestamp_accumulates", func(t *testing.T) {
		c := NewCounter()
		c.putTransaction(10.0, 500)
		c.putTransaction(20.0, 500)
		assert.Equal(t, 30.0, c.getTotalTransactionInLastOneHour(500))
	})
}
