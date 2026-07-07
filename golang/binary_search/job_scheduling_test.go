package binarysearch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_jobScheduling(t *testing.T) {
	testCases := []struct {
		name      string
		startTime []int
		endTime   []int
		profit    []int
		expected  int
	}{
		{
			name:      "leetcode_example_1",
			startTime: []int{1, 2, 3, 3},
			endTime:   []int{3, 4, 5, 6},
			profit:    []int{50, 10, 40, 70},
			expected:  120, // job[0] (1-3, $50) + job[3] (3-6, $70)
		},
		{
			name:      "leetcode_example_2",
			startTime: []int{1, 2, 3, 4, 6},
			endTime:   []int{3, 5, 10, 6, 9},
			profit:    []int{20, 20, 100, 70, 60},
			expected:  150, // job[0] + job[3] + job[4]
		},
		{
			name:      "leetcode_example_3",
			startTime: []int{1, 1, 1},
			endTime:   []int{2, 3, 4},
			profit:    []int{5, 6, 4},
			expected:  6, // best single job
		},
		{
			name:      "single_job",
			startTime: []int{1},
			endTime:   []int{5},
			profit:    []int{100},
			expected:  100,
		},
		{
			name:      "non_overlapping_jobs_take_all",
			startTime: []int{1, 3, 5},
			endTime:   []int{3, 5, 7},
			profit:    []int{10, 20, 30},
			expected:  60,
		},
		{
			name:      "all_overlapping_take_best",
			startTime: []int{1, 1, 1},
			endTime:   []int{4, 4, 4},
			profit:    []int{10, 30, 20},
			expected:  30,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := jobScheduling(tc.startTime, tc.endTime, tc.profit)
			assert.Equal(t, tc.expected, result)

			// top-down must agree with bottom-up
			result = jobSchedulingTopDown(tc.startTime, tc.endTime, tc.profit)
			assert.Equal(t, tc.expected, result)
		})
	}
}
