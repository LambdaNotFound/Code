package binarysearch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_bsearchRotatedSortedArray(t *testing.T) {
    testCases := []struct {
        name     string
        nums     []int
        target   int
        expected int
    }{
        {"target_in_right_half", []int{4, 5, 6, 7, 0, 1, 2}, 0, 4},
        {"target_not_found", []int{4, 5, 6, 7, 0, 1, 2}, 3, -1},
        {"single_element_miss", []int{1}, 0, -1},
        {"target_in_left_half", []int{4, 5, 6, 7, 0, 1, 2}, 6, 2},
        {"target_at_pivot", []int{3, 1}, 1, 1},
        {"not_rotated", []int{1, 2, 3, 4, 5}, 3, 2},
        {"single_element_hit", []int{5}, 5, 0},
        {"target_is_first", []int{4, 5, 6, 7, 0, 1, 2}, 4, 0},
        {"target_is_last", []int{4, 5, 6, 7, 0, 1, 2}, 2, 6},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := searchRotatedSortedArray(tc.nums, tc.target)
            assert.Equal(t, tc.expected, result)
        })
    }
}

func Test_searchRange(t *testing.T) {
    testCases := []struct {
        name     string
        nums     []int
        target   int
        expected []int
    }{
        {"found_range", []int{5, 7, 7, 8, 8, 10}, 8, []int{3, 4}},
        {"not_found", []int{5, 7, 7, 8, 8, 10}, 6, []int{-1, -1}},
        {"empty_array", []int{}, 8, []int{-1, -1}},
        {"single_occurrence", []int{1, 2, 3, 4, 5}, 3, []int{2, 2}},
        {"all_same", []int{2, 2, 2, 2}, 2, []int{0, 3}},
        {"target_at_start", []int{1, 1, 2, 3}, 1, []int{0, 1}},
        {"target_at_end", []int{1, 2, 3, 3}, 3, []int{2, 3}},
        {"single_element_hit", []int{5}, 5, []int{0, 0}},
        {"single_element_miss", []int{5}, 6, []int{-1, -1}},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := searchRange(tc.nums, tc.target)
            assert.Equal(t, tc.expected, result)
        })
    }
}

func Test_binarySearch(t *testing.T) {
    testCases := []struct {
        name     string
        nums     []int
        target   int
        expected int
    }{
        {"found_middle", []int{-1, 0, 3, 5, 9, 12}, 9, 4},
        {"not_found", []int{-1, 0, 3, 5, 9, 12}, 2, -1},
        {"target_first", []int{1, 2, 3, 4, 5}, 1, 0},
        {"target_last", []int{1, 2, 3, 4, 5}, 5, 4},
        {"single_element_hit", []int{7}, 7, 0},
        {"single_element_miss", []int{7}, 3, -1},
        {"negative_numbers", []int{-5, -3, -1, 0, 2}, -3, 1},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := binarySearch(tc.nums, tc.target)
            assert.Equal(t, tc.expected, result)
        })
    }
}
