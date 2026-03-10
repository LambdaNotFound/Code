package solid_coding

import (
	. "gocode/types"
	"sort"
)

/*
 * Intervals => sort intervals by start time, such that intervals[i].Start < intervals[j].Start
 *
 * 1. overlap: intervals[i].End >= intervals[j].Start
 *  i:        [a, ..., b]
 *  j:            [x, ..., y]               merge: End = max(intervals[i].End, intervals[j].End)
 *
 * 2. overlap: intervals[i].End >= intervals[j].Start && intervals[i].End >= intervals[j].End
 *  i:        [a, ............, b]
 *  j:            [x, ..., y]               merge: End = max(intervals[i].End, intervals[j].End)
 *
 * 3. no overlap: intervals[i].End < intervals[j].Start
 *  i:        [a, ..., b]
 *  j:                      [x, ..., y]
 *
 */

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
		newInterval[0] = min(intervals[i][0], newInterval[0])
		newInterval[1] = max(intervals[i][1], newInterval[1])
	}

	res = append(res, newInterval)
	for i < len(intervals) {
		res = append(res, intervals[i])
		i++
	}
	return res
}

func insertWithSlice(intervals [][]int, newInterval []int) [][]int {
	before, after := make([][]int, 0), make([][]int, 0)
	for i := 0; i < len(intervals); i++ {
		cur := intervals[i]
		if cur[1] < newInterval[0] {
			before = append(before, cur)
		} else if newInterval[1] < cur[0] {
			after = append(after, cur)
		} else {
			newInterval[0] = min(newInterval[0], cur[0])
			newInterval[1] = max(newInterval[1], cur[1])
		}
	}
	res := append(before, newInterval)
	res = append(res, after...)
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

	result := [][]int{intervals[0]}
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
 * 435. Non-overlapping Intervals
 *
 * return the minimum number of intervals you need to remove to make
 *     the rest of the intervals non-overlapping.
 *
 * A greedy strategy works well here. After sorting intervals by their start time,
 *     we process them from left to right and always keep the interval that ends
 *     earlier when an overlap occurs.
 */
func eraseOverlapIntervals(intervals [][]int) int {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	pre, erased := intervals[0], 0
	for _, interval := range intervals[1:] {
		if pre[1] <= interval[0] { // no erase
			pre = interval
		} else if interval[1] < pre[1] { // erase pre
			pre = interval
			erased += 1
		} else { // interval[0] < pre[1] < interval[1]
			erased += 1
		}
	}

	return erased
}

/**
 * Meeting Rooms
 *
 */
func canAttendMeetings(intervals []Interval) bool {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].Start < intervals[j].Start
	})

	for i := 1; i < len(intervals); i++ {
		if intervals[i].Start < intervals[i-1].End {
			return false
		}
	}

	return true
}

/**
 * Meeting Rooms II
 *
 * whats next earliest start time, next earliest end time?
 */
func minMeetingRooms(intervals []Interval) int {
	start, end := make([]int, len(intervals)), make([]int, len(intervals))
	for i := range intervals {
		start[i] = intervals[i].Start
		end[i] = intervals[i].End
	}
	sort.Ints(start)
	sort.Ints(end)

	res := 0
	for i, j, cnt := 0, 0, 0; i < len(start) && j < len(end); {
		if start[i] < end[j] {
			i += 1
			cnt += 1
		} else {
			j += 1
			cnt -= 1
		}
		res = max(res, cnt)
	}
	return res
}

func minMeetingRoomsSweepLine(intervals []Interval) int {
	hashmap := make(map[int]int) // time <> start++, end--
	for _, i := range intervals {
		hashmap[i.Start]++
		hashmap[i.End]--
	}

	keys := make([]int, 0, len(hashmap))
	for k := range hashmap {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	prev, res := 0, 0
	for _, k := range keys {
		prev += hashmap[k]
		res = max(res, prev)
	}
	return res
}
