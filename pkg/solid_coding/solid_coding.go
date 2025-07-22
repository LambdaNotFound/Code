package solid_coding

import (
	. "gocode/utils"
	"sort"
)

/**
 * 54. Spiral Matrix
 */
func spiralOrder(matrix [][]int) []int {
	nextDirection := map[string]string{
		"top":    "right",
		"right":  "bottom",
		"bottom": "left",
		"left":   "top",
	}

	m, n := len(matrix), len(matrix[0])
	res := []int{}
	rowTop, rowBottom := 0, m-1
	colLeft, colRight := 0, n-1
	direction := "top"

	for len(res) != m*n {
		if direction == "top" {
			for j := colLeft; j <= colRight; j++ {
				res = append(res, matrix[rowTop][j])
			}
			rowTop++
		} else if direction == "right" {
			for i := rowTop; i <= rowBottom; i++ {
				res = append(res, matrix[i][colRight])
			}
			colRight--
		} else if direction == "bottom" {
			for j := colRight; j >= colLeft; j-- {
				res = append(res, matrix[rowBottom][j])
			}
			rowBottom--
		} else if direction == "left" {
			for i := rowBottom; i >= rowTop; i-- {
				res = append(res, matrix[i][colLeft])
			}
			colLeft++
		}
		direction = nextDirection[direction]
	}

	return res
}

/**
 * 56. Merge Intervals
 *
 * [[1,4],[2,3]] => [[1,4]]
 */
func merge(intervals [][]int) [][]int {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	result := [][]int{}
	result = append(result, intervals[0])
	for _, interval := range intervals {
		current := result[len(result)-1]
		if current[1] < interval[0] {
			result = append(result, interval)
		} else {
			if current[1] < interval[1] {
				current[1] = interval[1]
			}
		}
	}

	return result
}

/**
 * 57. Insert Interval
 */
func insert(intervals [][]int, newInterval []int) [][]int {
	res := make([][]int, 0)

	i := 0
	for ; i < len(intervals) && intervals[i][1] < newInterval[0]; i++ {
		res = append(res, intervals[i])
	}
	for ; i < len(intervals) && intervals[i][0] <= newInterval[1]; i++ {
        // only update the newInterval[min, max]
		newInterval[0] = Min(intervals[i][0], newInterval[0])
		newInterval[1] = Max(intervals[i][1], newInterval[1])
	}

	res = append(res, newInterval)
	for i < len(intervals) {
		res = append(res, intervals[i])
		i++
	}
	return res
}
