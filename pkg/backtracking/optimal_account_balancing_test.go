package backtracking

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_minTransfers(t *testing.T) {
	testCases := []struct {
		name         string
		transactions [][]int
		expected     int
	}{
		{
			name:         "leetcode_example_1",
			transactions: [][]int{{0, 1, 10}, {2, 0, 5}},
			expected:     2,
		},
		{
			name:         "leetcode_example_2",
			transactions: [][]int{{0, 1, 10}, {1, 0, 1}, {1, 2, 5}, {2, 0, 5}},
			expected:     1,
		},
		{
			name:         "already_balanced",
			transactions: [][]int{{0, 1, 5}, {1, 0, 5}},
			expected:     0,
		},
		{
			name:         "single_transaction",
			transactions: [][]int{{0, 1, 100}},
			expected:     1,
		},
		{
			name:         "three_person_chain",
			transactions: [][]int{{0, 1, 10}, {1, 2, 10}},
			expected:     1, // 0 pays 2 directly
		},
		{
			name:         "three_person_cycle_all_equal",
			transactions: [][]int{{0, 1, 10}, {1, 2, 10}, {2, 0, 10}},
			expected:     0, // net balances all zero
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, minTransfers(tc.transactions))

			// both implementations must agree
			assert.Equal(t, tc.expected, minTransfersDP(tc.transactions))
		})
	}
}
