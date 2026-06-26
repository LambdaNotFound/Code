package bitmask

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_tspOpen(t *testing.T) {
	tests := []struct {
		name string
		dist [][]int
		want int
	}{
		{
			name: "single city",
			dist: [][]int{{0}},
			want: 0,
		},
		{
			name: "two cities",
			dist: [][]int{
				{0, 3},
				{3, 0},
			},
			want: 3,
		},
		{
			name: "three cities symmetric",
			// 0->1->2 = 1+2=3, 0->2->1 = 4+2=6
			dist: [][]int{
				{0, 1, 4},
				{1, 0, 2},
				{4, 2, 0},
			},
			want: 3,
		},
		{
			name: "four cities",
			dist: [][]int{
				{0, 10, 15, 20},
				{10, 0, 35, 25},
				{15, 35, 0, 30},
				{20, 25, 30, 0},
			},
			// Best open path: 0->1->3->2 = 10+25+30=65
			want: 65,
		},
		{
			name: "asymmetric distances",
			dist: [][]int{
				{0, 1, 100},
				{100, 0, 1},
				{1, 100, 0},
			},
			// 0->1->2 = 1+1=2
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tspOpen(tt.dist))
		})
	}
}

func Test_tspClosed(t *testing.T) {
	tests := []struct {
		name string
		dist [][]int
		want int
	}{
		{
			name: "single city",
			dist: [][]int{{0}},
			want: 0,
		},
		{
			name: "two cities",
			dist: [][]int{
				{0, 3},
				{3, 0},
			},
			// 0->1->0 = 3+3=6
			want: 6,
		},
		{
			name: "three cities symmetric",
			dist: [][]int{
				{0, 1, 4},
				{1, 0, 2},
				{4, 2, 0},
			},
			// 0->1->2->0 = 1+2+4=7, 0->2->1->0 = 4+2+1=7
			want: 7,
		},
		{
			name: "four cities classic",
			dist: [][]int{
				{0, 10, 15, 20},
				{10, 0, 35, 25},
				{15, 35, 0, 30},
				{20, 25, 30, 0},
			},
			// 0->1->3->2->0 = 10+25+30+15=80
			want: 80,
		},
		{
			name: "asymmetric distances",
			dist: [][]int{
				{0, 1, 100},
				{100, 0, 1},
				{1, 100, 0},
			},
			// 0->1->2->0 = 1+1+1=3
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tspClosed(tt.dist))
		})
	}
}
