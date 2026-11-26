package dynamic_programming

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_findCheapestPrice(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		flights  [][]int
		src, dst int
		k        int
		expected int
	}{
		{
			name:    "leetcode_example1",
			n:       4,
			flights: [][]int{{0, 1, 100}, {1, 2, 100}, {2, 0, 100}, {1, 3, 600}, {2, 3, 200}},
			src: 0, dst: 3, k: 1,
			expected: 700,
		},
		{
			name:    "leetcode_example2",
			n:       3,
			flights: [][]int{{0, 1, 100}, {1, 2, 100}, {0, 2, 500}},
			src: 0, dst: 2, k: 1,
			expected: 200,
		},
		{
			name:    "leetcode_example3_direct_only",
			n:       3,
			flights: [][]int{{0, 1, 100}, {1, 2, 100}, {0, 2, 500}},
			src: 0, dst: 2, k: 0,
			expected: 500,
		},
		{
			name:    "unreachable_returns_minus_one",
			n:       3,
			flights: [][]int{{0, 1, 100}},
			src: 0, dst: 2, k: 1,
			expected: -1,
		},
		{
			name:    "src_equals_dst",
			n:       2,
			flights: [][]int{{0, 1, 50}},
			src: 0, dst: 0, k: 0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, findCheapestPrice(tt.n, tt.flights, tt.src, tt.dst, tt.k))
		})
	}
}
