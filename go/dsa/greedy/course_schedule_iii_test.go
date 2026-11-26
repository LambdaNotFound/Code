package greedy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_scheduleCourse(t *testing.T) {
	tests := []struct {
		name     string
		courses  [][]int
		expected int
	}{
		{
			name:     "Example 1 (LeetCode)",
			courses:  [][]int{{100, 200}, {200, 1300}, {1000, 1250}, {2000, 3200}},
			expected: 3,
		},
		{
			// day=1 (take {1,2}), day=3 (take {2,3}), day=6 > 4 → drop longest (3), day=3; take 2
			name:     "All fit perfectly",
			courses:  [][]int{{1, 2}, {2, 3}, {3, 4}},
			expected: 2,
		},
		{
			name:     "None fit (all deadlines too short)",
			courses:  [][]int{{3, 2}, {4, 3}},
			expected: 0,
		},
		{
			name:     "Some fit (pick optimal subset)",
			courses:  [][]int{{5, 5}, {4, 6}, {2, 6}},
			expected: 2, // Can take {4,6} and {2,6} or {5,5} and {2,6}
		},
		{
			// only {2,4} and {1,100} fit; courses with deadlines 2 and 3 can't be kept
			name:     "Tight deadlines require optimal replacement",
			courses:  [][]int{{3, 2}, {4, 3}, {2, 4}, {1, 100}},
			expected: 2,
		},
		{
			name:     "Empty input",
			courses:  [][]int{},
			expected: 0,
		},
		{
			name:     "Single course fits",
			courses:  [][]int{{5, 10}},
			expected: 1,
		},
		{
			name:     "Single course doesn't fit",
			courses:  [][]int{{10, 5}},
			expected: 0,
		},
		{
			name:     "Courses with same deadlines but different durations",
			courses:  [][]int{{1, 5}, {2, 5}, {3, 5}, {4, 5}},
			expected: 2, // Best pick shortest ones: {1,5} + {2,5}
		},
		{
			// {2,2} fits, {3,4} pushed day to 5 > 4 so drop {3,4}, then {10,100} fits; total 2
			name:     "Large course late in schedule",
			courses:  [][]int{{2, 2}, {3, 4}, {10, 100}},
			expected: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := scheduleCourseHeap(tc.courses)
			assert.Equal(t, tc.expected, got)

			got = scheduleCourse(tc.courses)
			assert.Equal(t, tc.expected, got)
		})
	}
}
