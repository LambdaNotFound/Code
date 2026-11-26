package matrix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_multiplyMatrix(t *testing.T) {
	tests := []struct {
		name     string
		mat1     [][]int
		mat2     [][]int
		expected [][]int
	}{
		{
			name:     "leetcode_example",
			mat1:     [][]int{{1, 0, 0}, {-1, 0, 3}},
			mat2:     [][]int{{7, 0, 0}, {0, 0, 0}, {0, 0, 1}},
			expected: [][]int{{7, 0, 0}, {-7, 0, 3}},
		},
		{
			name:     "identity",
			mat1:     [][]int{{1, 0}, {0, 1}},
			mat2:     [][]int{{1, 0}, {0, 1}},
			expected: [][]int{{1, 0}, {0, 1}},
		},
		{
			name:     "all_zeros_mat1",
			mat1:     [][]int{{0, 0}, {0, 0}},
			mat2:     [][]int{{1, 2}, {3, 4}},
			expected: [][]int{{0, 0}, {0, 0}},
		},
		{
			name:     "1x1",
			mat1:     [][]int{{3}},
			mat2:     [][]int{{4}},
			expected: [][]int{{12}},
		},
		{
			name:     "1x2_times_2x1",
			mat1:     [][]int{{1, 2}},
			mat2:     [][]int{{3}, {4}},
			expected: [][]int{{11}},
		},
		{
			name:     "all_ones",
			mat1:     [][]int{{1, 1}, {1, 1}},
			mat2:     [][]int{{1, 1}, {1, 1}},
			expected: [][]int{{2, 2}, {2, 2}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := multiply(tt.mat1, tt.mat2)
			assert.Equal(t, tt.expected, got, "multiply")

			got = multiplyCompressed(tt.mat1, tt.mat2)
			assert.Equal(t, tt.expected, got, "multiplyCompressed")
		})
	}
}
