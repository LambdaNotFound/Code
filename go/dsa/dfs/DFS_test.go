package dfs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_floodFill(t *testing.T) {
	tests := []struct {
		name     string
		image    [][]int
		sr, sc   int
		newColor int
		expected [][]int
	}{
		{
			name: "basic example",
			image: [][]int{
				{1, 1, 1},
				{1, 1, 0},
				{1, 0, 1},
			},
			sr: 1, sc: 1, newColor: 2,
			expected: [][]int{
				{2, 2, 2},
				{2, 2, 0},
				{2, 0, 1},
			},
		},
		{
			name: "no change needed (already new color)",
			image: [][]int{
				{0, 0, 0},
				{0, 1, 1},
			},
			sr: 1, sc: 1, newColor: 1,
			expected: [][]int{
				{0, 0, 0},
				{0, 1, 1},
			},
		},
		{
			name: "single pixel image",
			image: [][]int{
				{0},
			},
			sr: 0, sc: 0, newColor: 2,
			expected: [][]int{
				{2},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := floodFill(tc.image, tc.sr, tc.sc, tc.newColor)
			assert.Equal(t, tc.expected, got)
		})
	}
}
