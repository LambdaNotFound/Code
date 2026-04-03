package binarysearch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_searchRotatedSortedArrayII(t *testing.T) {
	testCases := []struct {
		name     string
		nums     []int
		target   int
		expected bool
	}{
		{"found_in_right_half", []int{2, 5, 6, 0, 0, 1, 2}, 0, true},
		{"not_found", []int{2, 5, 6, 0, 0, 1, 2}, 3, false},
		{"found_with_duplicates_at_boundary", []int{1, 0, 1, 1, 1}, 0, true},
		{"all_same_found", []int{2, 2, 2, 2}, 2, true},
		{"all_same_not_found", []int{2, 2, 2, 2}, 3, false},
		{"single_element_hit", []int{5}, 5, true},
		{"single_element_miss", []int{5}, 3, false},
		{"not_rotated", []int{1, 2, 3, 4, 5}, 4, true},
		{"duplicate_at_pivot", []int{3, 1, 1}, 3, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, searchRotatedSortedArrayII(tc.nums, tc.target))
		})
	}
}

func Test_findMinInRotatedSortedArray(t *testing.T) {
	testCases := []struct {
		name     string
		nums     []int
		expected int
	}{
		{"rotated", []int{3, 4, 5, 1, 2}, 1},
		{"rotated_more", []int{4, 5, 6, 7, 0, 1, 2}, 0},
		{"not_rotated", []int{11, 13, 15, 17}, 11},
		{"single_element", []int{1}, 1},
		{"two_elements_rotated", []int{2, 1}, 1},
		{"two_elements_sorted", []int{1, 2}, 1},
		{"rotation_at_last", []int{2, 3, 4, 5, 1}, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := findMinInRotatedSortedArray(tc.nums)
			assert.Equal(t, tc.expected, result)

			// both implementations must agree
			result = findMinInRotatedSortedArrayII(tc.nums)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_findMinInRotatedSortedArrayII(t *testing.T) {
	testCases := []struct {
		name     string
		nums     []int
		expected int
	}{
		{"with_duplicates", []int{2, 2, 2, 0, 1}, 0},
		{"all_same", []int{1, 1, 1, 1}, 1},
		{"duplicates_at_boundary", []int{1, 3, 3}, 1},
		{"rotated_with_duplicates", []int{3, 1, 3, 3, 3}, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, findMinInRotatedSortedArrayII(tc.nums))
		})
	}
}

func Test_searchMatrix(t *testing.T) {
	testCases := []struct {
		name     string
		matrix   [][]int
		target   int
		expected bool
	}{
		{
			name:     "found",
			matrix:   [][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}},
			target:   3,
			expected: true,
		},
		{
			name:     "not_found",
			matrix:   [][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}},
			target:   13,
			expected: false,
		},
		{
			name:     "single_cell_hit",
			matrix:   [][]int{{1}},
			target:   1,
			expected: true,
		},
		{
			name:     "single_cell_miss",
			matrix:   [][]int{{1}},
			target:   2,
			expected: false,
		},
		{
			name:     "target_at_last_cell",
			matrix:   [][]int{{1, 3, 5}, {7, 9, 11}},
			target:   11,
			expected: true,
		},
		{
			name:     "target_at_first_cell",
			matrix:   [][]int{{1, 3, 5}, {7, 9, 11}},
			target:   1,
			expected: true,
		},
		{
			name:     "target_less_than_min",
			matrix:   [][]int{{5, 10}, {15, 20}},
			target:   3,
			expected: false,
		},
		{
			name:     "target_greater_than_max",
			matrix:   [][]int{{5, 10}, {15, 20}},
			target:   25,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, searchMatrix(tc.matrix, tc.target))
		})
	}
}

func Test_searchMatrix2(t *testing.T) {
	testCases := []struct {
		name     string
		matrix   [][]int
		target   int
		expected bool
	}{
		{
			name: "found",
			matrix: [][]int{
				{1, 4, 7, 11, 15},
				{2, 5, 8, 12, 19},
				{3, 6, 9, 16, 22},
				{10, 13, 14, 17, 24},
				{18, 21, 23, 26, 30},
			},
			target:   5,
			expected: true,
		},
		{
			name: "not_found",
			matrix: [][]int{
				{1, 4, 7, 11, 15},
				{2, 5, 8, 12, 19},
				{3, 6, 9, 16, 22},
				{10, 13, 14, 17, 24},
				{18, 21, 23, 26, 30},
			},
			target:   20,
			expected: false,
		},
		{
			name:     "single_cell_hit",
			matrix:   [][]int{{5}},
			target:   5,
			expected: true,
		},
		{
			name:     "single_cell_miss",
			matrix:   [][]int{{5}},
			target:   3,
			expected: false,
		},
		{
			name:     "target_at_top_right",
			matrix:   [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			target:   3,
			expected: true,
		},
		{
			name:     "target_at_bottom_left",
			matrix:   [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			target:   7,
			expected: true,
		},
		{
			name:     "empty_matrix",
			matrix:   [][]int{},
			target:   1,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, searchMatrix2(tc.matrix, tc.target))
		})
	}
}
