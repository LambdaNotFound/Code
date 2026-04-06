package monoqueue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_maxSlidingWindow(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected []int
	}{
		{"leetcode_1", []int{1, 3, -1, -3, 5, 3, 6, 7}, 3, []int{3, 3, 5, 5, 6, 7}},
		{"leetcode_2", []int{1}, 1, []int{1}},
		{"k_equals_len", []int{4, 2, 5, 1}, 4, []int{5}},
		{"k_one", []int{3, 1, 2}, 1, []int{3, 1, 2}},
		{"all_same", []int{5, 5, 5, 5}, 2, []int{5, 5, 5}},
		{"decreasing", []int{10, 9, 8, 7}, 2, []int{10, 9, 8}},
		{"increasing", []int{1, 2, 3, 4}, 2, []int{2, 3, 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input1 := append([]int(nil), tt.nums...)
			assert.Equal(t, tt.expected, maxSlidingWindow(input1, tt.k), "maxSlidingWindow")
		})
	}
}
