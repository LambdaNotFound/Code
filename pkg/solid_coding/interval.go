package solid_coding

import (
	"container/heap"
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
 *
 * Two intervals [a,b] and [c,d] overlap iff: a <= d AND c <= b
 *
 * They DON'T overlap iff:
 *   b < c  (first ends before second starts)
 * OR
 *   d < a  (second ends before first starts)
 *
 */

// merge intervals template
func mergeIntervals(intervals [][]int) [][]int {
	// step 1: sort by start time
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	merged := [][]int{intervals[0]}
	for i := 1; i < len(intervals); i++ {
		last := merged[len(merged)-1]
		interval := intervals[i]

		if interval[0] <= last[1] { // overlap: merge
			last[1] = max(last[1], interval[1])
		} else { // no overlap: append
			merged = append(merged, interval)
		}
	}
	return merged
}

/**
 * 57. Insert Interval
 */
func insert(intervals [][]int, newInterval []int) [][]int {
	before, after := make([][]int, 0), make([][]int, 0)
	for i := 0; i < len(intervals); i++ {
		interval := intervals[i]
		if interval[1] < newInterval[0] {
			before = append(before, interval)
		} else if newInterval[1] < interval[0] {
			after = append(after, interval)
		} else {
			newInterval[0] = min(newInterval[0], interval[0])
			newInterval[1] = max(newInterval[1], interval[1])
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
		last := result[len(result)-1]
		if last[1] < interval[0] {
			result = append(result, interval)
		} else {
			last[1] = max(last[1], interval[1])
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
	erased := 0
	for prev, j := 0, 1; j < len(intervals); j++ {
		if intervals[prev][1] <= intervals[j][0] {
			prev = j
		} else if intervals[prev][1] > intervals[j][1] {
			erased++
			prev = j
		} else if intervals[prev][1] <= intervals[j][1] {
			erased++
		}
	}
	return erased
}

func eraseOverlapIntervalsSortByEndTime(intervals [][]int) int {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][1] < intervals[j][1]
	})
	erased := 0
	for prev, j := 0, 1; j < len(intervals); j++ {
		if intervals[j][0] < intervals[prev][1] {
			erased++
		} else {
			prev = j
		}
	}
	return erased
}

// interval intersections template, Two pointers:
func intersect(listA, listB [][]int) [][]int {
	intersections := [][]int{}
	i, j := 0, 0

	for i < len(listA) && j < len(listB) {
		// find overlap
		overlapStart := max(listA[i][0], listB[j][0])
		overlapEnd := min(listA[i][1], listB[j][1])

		if overlapStart <= overlapEnd {
			intersections = append(intersections, []int{overlapStart, overlapEnd})
		}

		// advance pointer with smaller end
		if listA[i][1] < listB[j][1] {
			i++
		} else {
			j++
		}
	}
	return intersections
}

/**
 * 986. Interval List Intersections
 */
func intervalIntersection(firstList [][]int, secondList [][]int) [][]int {
	intersections := [][]int{}
	for i, j := 0, 0; i < len(firstList) && j < len(secondList); {
		start := max(firstList[i][0], secondList[j][0])
		end := min(firstList[i][1], secondList[j][1])
		if start <= end {
			intersections = append(intersections, []int{start, end})
		}

		if firstList[i][1] < secondList[j][1] {
			i++
		} else {
			j++
		}
	}

	return intersections
}

/**
 * Meeting Rooms
 *
 * check if there's overlap
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
 * find the minimum number of rooms required to schedule all meetings without any conflicts.
 *
 * whats next earliest start time, next earliest end time?
 *
 * 1. Sweep line
 */
func minMeetingRooms(intervals []Interval) int {
	timeline := make(map[int]int) // time point → room delta (start: +1, end: -1)
	for _, meeting := range intervals {
		timeline[meeting.Start]++
		timeline[meeting.End]--
	}

	timePoints := make([]int, 0, len(timeline))
	for timePoint := range timeline {
		timePoints = append(timePoints, timePoint)
	}
	sort.Ints(timePoints)

	rooms, maxRooms := 0, 0
	for _, timePoint := range timePoints {
		rooms += timeline[timePoint]
		maxRooms = max(maxRooms, rooms)
	}
	return maxRooms
}

func minMeetingRoomsMinHeap(intervals [][]int) int {
	// sort meetings by start time
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	// min-heap of end times — tracks when each room becomes free
	minHeap := &EndTimeHeap{}
	heap.Init(minHeap)

	for _, interval := range intervals {
		start, end := interval[0], interval[1]

		if minHeap.Len() > 0 && (*minHeap)[0] <= start {
			// earliest-ending room is free — reuse it
			heap.Pop(minHeap)
		}
		// assign meeting to room (new or reused)
		heap.Push(minHeap, end)
	}

	return minHeap.Len()
}

// min-heap ordered by end time
type EndTimeHeap []int

func (h EndTimeHeap) Len() int           { return len(h) }
func (h EndTimeHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h EndTimeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *EndTimeHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *EndTimeHeap) Pop() interface{} {
	old := *h
	item := old[len(old)-1]
	*h = old[:len(old)-1]
	return item
}

/**
 * 1094. Car Pooling
 */
func carPooling(trips [][]int, capacity int) bool {
	type event struct {
		pos, delta int
	}

	// each trip generates a board and alight event
	events := make([]event, 0, len(trips)*2)
	for _, trip := range trips {
		numPassengers, from, to := trip[0], trip[1], trip[2]
		events = append(events, event{pos: from, delta: +numPassengers}) // board
		events = append(events, event{pos: to, delta: -numPassengers})   // off-board
	}

	// sort by position; alight before board at same position
	// (passengers leave before new ones board at the same stop)
	sort.Slice(events, func(i, j int) bool {
		if events[i].pos != events[j].pos {
			return events[i].pos < events[j].pos
		}
		return events[i].delta < events[j].delta // negative delta first
	})

	currentPassengers := 0
	for _, e := range events {
		currentPassengers += e.delta
		if currentPassengers > capacity {
			return false
		}
	}
	return true
}

/**
 * Employee Free Time
 *
 * Write a function to find the common free time for all employees from a list called schedule.
 * Each employee's schedule is represented by a list of non-overlapping intervals sorted by start times.
 * The function should return a list of finite, non-zero length intervals where all employees are free, also sorted in order.
 *
 * Input: schedule = [[[2,4],[7,10]],[[1,5]],[[6,9]]]
 * Output: [(5,6)]
 * Explanation: The three employees collectively have only one common free time interval, which is from 5 to 6.
 */
func employeeFreeTime(schedule [][][]int) [][]int {
	// Step 1: flatten all intervals
	flattened := make([][]int, 0)
	for _, employee := range schedule {
		flattened = append(flattened, employee...)
	}

	// Step 2: sort by start time
	sort.Slice(flattened, func(i, j int) bool {
		return flattened[i][0] < flattened[j][0]
	})

	// Step 3: merge overlapping intervals
	merged := [][]int{flattened[0]} // ← seed with first interval
	for _, interval := range flattened[1:] {
		last := merged[len(merged)-1]
		if last[1] < interval[0] {
			merged = append(merged, interval)
		} else {
			last[1] = max(last[1], interval[1])
		}
	}

	// Step 4: gaps between merged intervals = free time
	freeTimes := make([][]int, 0)
	for i := 1; i < len(merged); i++ {
		freeTimes = append(freeTimes, []int{merged[i-1][1], merged[i][0]})
	}

	return freeTimes
}
