package interview

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

var purchaseRecords = [][]string{
	{"A", "B", "C"},
	{"A", "B"},
	{"B", "C", "D"},
	{"A", "C"},
}

func Test_highCorrelationSellers(t *testing.T) {
	t.Run("threshold_2", func(t *testing.T) {
		result := highCorrelationSellers(purchaseRecords, 2)
		sort.Strings(result["A"])
		sort.Strings(result["B"])
		sort.Strings(result["C"])
		assert.Equal(t, []string{"B", "C"}, result["A"])
		assert.Equal(t, []string{"A", "C"}, result["B"])
		assert.Equal(t, []string{"A", "B"}, result["C"])
		assert.Equal(t, []string{}, result["D"])
	})

	t.Run("threshold_too_high_no_results", func(t *testing.T) {
		result := highCorrelationSellers(purchaseRecords, 10)
		for _, v := range result {
			assert.Empty(t, v)
		}
	})

	t.Run("single_session", func(t *testing.T) {
		result := highCorrelationSellers([][]string{{"X", "Y"}}, 1)
		assert.Equal(t, []string{"Y"}, result["X"])
		assert.Equal(t, []string{"X"}, result["Y"])
	})
}

func Test_topKCorrelatedSellers(t *testing.T) {
	t.Run("top_2", func(t *testing.T) {
		result := topKCorrelatedSellers(purchaseRecords, 2)
		assert.Len(t, result["A"], 2)
		assert.Len(t, result["B"], 2)
		assert.Len(t, result["C"], 2)
	})

	t.Run("k_larger_than_neighbors", func(t *testing.T) {
		result := topKCorrelatedSellers(purchaseRecords, 100)
		assert.Len(t, result["D"], 2) // D co-occurs with B and C in session {"B","C","D"}
	})

	t.Run("k_zero", func(t *testing.T) {
		result := topKCorrelatedSellers(purchaseRecords, 0)
		for _, v := range result {
			assert.Empty(t, v)
		}
	})
}

func Test_scaledCorrelation(t *testing.T) {
	sessions := [][]string{
		{"A", "B", "C"},
		{"A", "B"},
		{"B", "C", "D"},
		{"A", "C"},
	}
	ch := make(chan []string, len(sessions))
	for _, s := range sessions {
		ch <- s
	}
	close(ch)

	result := scaledCorrelation(ch, 2, 2)

	// A-B pair appears 2 times, A-C 2 times, B-C 2 times
	assert.Contains(t, result["A"], "B")
	assert.Contains(t, result["A"], "C")
	assert.Contains(t, result["B"], "A")
	assert.Contains(t, result["B"], "C")
}
