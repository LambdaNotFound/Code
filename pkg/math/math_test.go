package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_productExceptSelf(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected []int
	}{
		{
			name:     "basic example",
			nums:     []int{1, 2, 3, 4},
			expected: []int{24, 12, 8, 6},
		},
		{
			name:     "with zero",
			nums:     []int{-1, 1, 0, -3, 3},
			expected: []int{0, 0, 9, 0, 0},
		},
		{
			name:     "single element",
			nums:     []int{5},
			expected: []int{1}, // usually problem guarantees len >= 2, but we handle it
		},
		{
			name:     "two elements",
			nums:     []int{2, 3},
			expected: []int{3, 2},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := productExceptSelf(tc.nums)
			assert.Equal(t, tc.expected, got)

			got = productExceptSelfClaude(tc.nums)
			assert.Equal(t, tc.expected, got)
		})
	}
}
