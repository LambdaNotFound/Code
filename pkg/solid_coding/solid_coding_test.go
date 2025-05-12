package solid_coding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_spiralOrder(t *testing.T) {
	testCases := []struct {
		name     string
		matrix   [][]int
		expected []int
	}{
		{
			"case 1",
			[][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			[]int{1, 2, 3, 6, 9, 8, 7, 4, 5},
		},
		{
			"case 2",
			[][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}},
			[]int{1, 2, 3, 4, 8, 12, 11, 10, 9, 5, 6, 7},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := spiralOrder(tc.matrix)
			assert.Equal(t, tc.expected, result)
		})
	}
}
