package solid_coding

import (
	. "gocode/types"
	"sort"
)

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
