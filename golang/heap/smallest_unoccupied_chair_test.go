package heap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_smallestChair(t *testing.T) {
	tests := []struct {
		name         string
		times        [][]int
		targetFriend int
		want         int
	}{
		{
			// Friend 1 arrives at 2, friends 0 and 2 arrive at 1 and 4.
			// Friend 0 takes chair 0. Friend 1 arrives at 2 (chair 0 still occupied until 3), takes chair 1.
			name:         "leetcode_example_1",
			times:        [][]int{{1, 4}, {2, 3}, {4, 6}},
			targetFriend: 1,
			want:         1,
		},
		{
			// Friend 1 arrives at 2; friend 0 leaves at 1 so chair 0 is free.
			name:         "leetcode_example_2",
			times:        [][]int{{3, 10}, {1, 5}, {2, 6}},
			targetFriend: 0,
			want:         2,
		},
		{
			// Only one friend — always gets chair 0.
			name:         "single_friend",
			times:        [][]int{{1, 5}},
			targetFriend: 0,
			want:         0,
		},
		{
			// Friends arrive sequentially, each after the previous leaves.
			// All reuse chair 0 in turn; target is the last one.
			name:         "sequential_reuse_chair_0",
			times:        [][]int{{1, 2}, {3, 4}, {5, 6}},
			targetFriend: 2,
			want:         0,
		},
		{
			// Friends 0 and 1 overlap (0→[5,10], 1→[1,20]).
			// After sorting: friend 1 arrives at 1 (chair 0), friend 0 arrives at 5 (chair 1).
			// Target is friend 0 (arrival 5) → chair 1.
			name:         "target_gets_second_chair",
			times:        [][]int{{5, 10}, {1, 20}},
			targetFriend: 0,
			want:         1,
		},
		{
			// Friends 0,1,2 all arrive before target and all leave before target arrives.
			// Sorted: [1,5],[2,6],[3,7],[10,15]
			// At target arrival (10) all three chairs (0,1,2) pile into availableHeap together,
			// forcing the heap comparator to run for the 2nd and 3rd pushes before picking chair 0.
			name:         "multiple_freed_seats_comparator_invoked",
			times:        [][]int{{1, 5}, {2, 6}, {3, 7}, {10, 15}},
			targetFriend: 3,
			want:         0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, smallestChair(tt.times, tt.targetFriend))
		})
	}
}
