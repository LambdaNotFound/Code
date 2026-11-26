package interview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_trapRainWater(t *testing.T) {
	tests := []struct {
		name      string
		heightMap [][]int
		want      int
	}{
		{
			name: "leetcode_example_1",
			heightMap: [][]int{
				{1, 4, 3, 1, 3, 2},
				{3, 2, 1, 3, 2, 4},
				{2, 3, 3, 2, 3, 1},
			},
			want: 4,
		},
		{
			name: "leetcode_example_2",
			heightMap: [][]int{
				{3, 3, 3, 3, 3},
				{3, 2, 2, 2, 3},
				{3, 2, 1, 2, 3},
				{3, 2, 2, 2, 3},
				{3, 3, 3, 3, 3},
			},
			want: 10,
		},
		{
			name: "flat_surface_no_water",
			heightMap: [][]int{
				{1, 1, 1},
				{1, 1, 1},
				{1, 1, 1},
			},
			want: 0,
		},
		{
			name: "single_row_no_water",
			heightMap: [][]int{
				{1, 2, 3, 4, 5},
			},
			want: 0,
		},
		{
			name: "single_column_no_water",
			heightMap: [][]int{
				{1}, {2}, {3}, {4}, {5},
			},
			want: 0,
		},
		{
			name: "too_small_no_interior",
			heightMap: [][]int{
				{5, 5},
				{5, 5},
			},
			want: 0,
		},
		{
			name: "single_cell",
			heightMap: [][]int{
				{7},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, trapRainWater(tt.heightMap))
		})
	}
}
