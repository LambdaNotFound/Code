package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_numBusesToDestination(t *testing.T) {
	tests := []struct {
		name     string
		routes   [][]int
		source   int
		target   int
		expected int
	}{
		{
			name:     "leetcode_example1",
			routes:   [][]int{{1, 2, 7}, {3, 6, 7}},
			source:   1, target: 6,
			expected: 2,
		},
		{
			name:     "leetcode_example2",
			routes:   [][]int{{7, 12}, {4, 5, 15}, {6}},
			source:   15, target: 12,
			expected: -1,
		},
		{
			name:     "source_equals_target",
			routes:   [][]int{{1, 2, 3}},
			source:   1, target: 1,
			expected: 0,
		},
		{
			name:     "direct_route",
			routes:   [][]int{{1, 2, 3, 4}},
			source:   1, target: 4,
			expected: 1,
		},
		{
			name:     "two_buses_needed",
			routes:   [][]int{{1, 2}, {2, 3}},
			source:   1, target: 3,
			expected: 2,
		},
		{
			name:     "unreachable_target",
			routes:   [][]int{{1, 2, 3}},
			source:   1, target: 5,
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, numBusesToDestinationClaude(tt.routes, tt.source, tt.target))
		})
	}
}
