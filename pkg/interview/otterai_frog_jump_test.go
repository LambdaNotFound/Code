package interview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_canCross(t *testing.T) {
	tests := []struct {
		name   string
		stones []int
		want   bool
	}{
		{
			name:   "leetcode_example_1_can_cross",
			stones: []int{0, 1, 3, 5, 6, 8, 12, 17},
			want:   true,
		},
		{
			name:   "leetcode_example_2_cannot_cross",
			stones: []int{0, 1, 2, 3, 4, 8, 9, 11},
			want:   false,
		},
		{
			name:   "single_stone_already_at_target",
			stones: []int{0},
			want:   true,
		},
		{
			name:   "second_stone_not_at_distance_one",
			stones: []int{0, 2},
			want:   false,
		},
		{
			name:   "second_stone_at_distance_one",
			stones: []int{0, 1},
			want:   true,
		},
		{
			name:   "consecutive_integers",
			stones: []int{0, 1, 2, 3, 4, 5},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, canCross(tt.stones))
		})
	}
}
