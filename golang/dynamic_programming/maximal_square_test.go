package dynamic_programming

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_maximalSquare(t *testing.T) {
	tests := []struct {
		name     string
		matrix   [][]byte
		expected int
	}{
		{
			name: "example1",
			matrix: [][]byte{
				{'1', '0', '1', '0', '0'},
				{'1', '0', '1', '1', '1'},
				{'1', '1', '1', '1', '1'},
				{'1', '0', '0', '1', '0'},
			},
			expected: 4,
		},
		{
			name:     "example2",
			matrix:   [][]byte{{'0', '1'}, {'1', '0'}},
			expected: 1,
		},
		{
			name:     "all_zeros",
			matrix:   [][]byte{{'0'}},
			expected: 0,
		},
		{
			name:     "all_ones_1x1",
			matrix:   [][]byte{{'1'}},
			expected: 1,
		},
		{
			name: "all_ones_3x3",
			matrix: [][]byte{
				{'1', '1', '1'},
				{'1', '1', '1'},
				{'1', '1', '1'},
			},
			expected: 9,
		},
		{
			name: "single_row",
			matrix: [][]byte{{'1', '1', '1', '1'}},
			expected: 1,
		},
		{
			name: "single_col",
			matrix: [][]byte{{'1'}, {'1'}, {'1'}, {'1'}},
			expected: 1,
		},
		{
			name: "L_shape",
			matrix: [][]byte{
				{'1', '1', '0'},
				{'1', '1', '0'},
				{'0', '0', '0'},
			},
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, maximalSquare(tt.matrix))
		})
	}
}
