package apidesign

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
			name:         "already_settled",
			transactions: [][]int{{0, 1, 5}, {1, 0, 5}},
			expected:     0,
		},
		{
			name:         "three_person_cycle",
			transactions: [][]int{{0, 1, 10}, {1, 2, 10}, {2, 0, 10}},
			expected:     0, // all balances are zero
		},
		{
			name:         "single_transaction",
			transactions: [][]int{{0, 1, 100}},
			expected:     1,
		},
		{
			name: "chain_settles_in_two",
			// 0 owes 1 ten, 1 owes 2 ten → 0 pays 2 directly: 1 transfer
			transactions: [][]int{{0, 1, 10}, {1, 2, 10}},
			expected:     1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, minTransfers(tc.transactions))
		})
	}
}
