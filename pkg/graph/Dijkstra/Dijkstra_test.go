package dijkstra

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
			src:     0, dst: 3, k: 1,
			expected: 700,
		},
		{
			name:    "leetcode_example2",
			n:       3,
			flights: [][]int{{0, 1, 100}, {1, 2, 100}, {0, 2, 500}},
			src:     0, dst: 2, k: 1,
			expected: 200,
		},
		{
			name:    "leetcode_example3_direct_only",
			n:       3,
			flights: [][]int{{0, 1, 100}, {1, 2, 100}, {0, 2, 500}},
			src:     0, dst: 2, k: 0,
			expected: 500,
		},
		{
			name:    "no_path",
			n:       3,
			flights: [][]int{{0, 1, 100}},
			src:     0, dst: 2, k: 1,
			expected: -1,
		},
		{
			name:    "direct_flight",
			n:       2,
			flights: [][]int{{0, 1, 50}},
			src:     0, dst: 1, k: 0,
			expected: 50,
		},
		{
			name:    "k_limits_cheaper_route",
			n:       4,
			flights: [][]int{{0, 1, 1}, {1, 2, 1}, {2, 3, 1}, {0, 3, 100}},
			src:     0, dst: 3, k: 1,
			// 0→1→2→3 needs 2 stops (k=1 allows 1 stop), only 0→3=100 works
			expected: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, findCheapestPrice(tt.n, tt.flights, tt.src, tt.dst, tt.k))
		})
	}
}
