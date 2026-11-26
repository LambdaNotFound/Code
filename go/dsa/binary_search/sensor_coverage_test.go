package binarysearch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_minSensorRange(t *testing.T) {
	tests := []struct {
		name      string
		crossings []int
		towers    []int
		want      int
	}{
		{
			// given example: crossing 12 is 6 away from the nearest tower (6)
			name:      "given_example",
			crossings: []int{1, 2, 3, 4, 6, 10, 12},
			towers:    []int{1, 4, 6},
			want:      6,
		},
		{
			// single tower in the middle; farthest crossings (1 and 5) are dist 2
			name:      "single_tower_centered",
			crossings: []int{1, 2, 3, 4, 5},
			towers:    []int{3},
			want:      2,
		},
		{
			// two towers evenly spaced; every crossing is dist 1 from a tower
			name:      "two_towers_all_adjacent",
			crossings: []int{1, 3, 5, 7},
			towers:    []int{2, 6},
			want:      1,
		},
		{
			// single tower, single crossing — exact position match
			name:      "crossing_at_tower",
			crossings: []int{5},
			towers:    []int{5},
			want:      0,
		},
		{
			// single tower between two crossings; farthest is crossing 10 (dist 5)
			name:      "single_tower_two_crossings",
			crossings: []int{1, 10},
			towers:    []int{5},
			want:      5,
		},
		{
			// towers on both ends; middle crossing (5) is dist 4 from tower 1, dist 5 from tower 10
			name:      "towers_at_ends",
			crossings: []int{1, 5, 10},
			towers:    []int{1, 10},
			want:      4,
		},
		{
			// crossing far from lone tower
			name:      "single_crossing_far_tower",
			crossings: []int{10},
			towers:    []int{1},
			want:      9,
		},
		{
			// unsorted inputs — function must handle any order
			name:      "unsorted_inputs",
			crossings: []int{12, 1, 6, 4, 10, 2, 3},
			towers:    []int{6, 1, 4},
			want:      6,
		},
		{
			// all crossings coincide with towers — range 0 suffices
			name:      "all_crossings_at_towers",
			crossings: []int{1, 4, 6},
			towers:    []int{1, 4, 6},
			want:      0,
		},
		{
			// single crossing, single tower
			name:      "single_crossing_single_tower",
			crossings: []int{7},
			towers:    []int{3},
			want:      4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, minSensorRange(tt.crossings, tt.towers))
		})
	}
}

func Test_findRadius(t *testing.T) {
	tests := []struct {
		name    string
		houses  []int
		heaters []int
		want    int
	}{
		{
			// houses 1 and 3 each need radius 1 from heater 2; house 4 needs radius 1 from heater 4
			name:    "leetcode_example_1",
			houses:  []int{1, 2, 3},
			heaters: []int{2},
			want:    1,
		},
		{
			// heaters at 1 and 3; house 2 is distance 1 from both; house 4 is distance 1 from 3
			name:    "leetcode_example_2",
			houses:  []int{1, 2, 3, 4},
			heaters: []int{1, 4},
			want:    1,
		},
		{
			// house and heater co-located — radius 0 suffices
			name:    "house_at_heater",
			houses:  []int{5},
			heaters: []int{5},
			want:    0,
		},
		{
			// single house, single heater far apart
			name:    "single_house_far_heater",
			houses:  []int{1},
			heaters: []int{10},
			want:    9,
		},
		{
			// heater to the left of all houses
			name:    "heater_left_of_all",
			houses:  []int{3, 7, 10},
			heaters: []int{1},
			want:    9,
		},
		{
			// heater to the right of all houses
			name:    "heater_right_of_all",
			houses:  []int{1, 4, 8},
			heaters: []int{10},
			want:    9,
		},
		{
			// unsorted inputs — function must sort before searching
			// sorted heaters [2,7]; house 4 → min(|4-2|,|4-7|)=2; house 1 → 1; house 3 → 1
			name:    "unsorted_inputs",
			houses:  []int{4, 1, 3},
			heaters: []int{7, 2},
			want:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, findRadius(tt.houses, tt.heaters))
		})
	}
}
