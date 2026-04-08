package solid_coding

import (
	"testing"

	. "gocode/types"

	"github.com/stretchr/testify/assert"
)

func Test_merge(t *testing.T) {
	testCases := []struct {
		name      string
		intervals [][]int
		expected  [][]int
	}{
		{
			"case 1",
			[][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}},
			[][]int{{1, 6}, {8, 10}, {15, 18}},
		},
		{
			"case 2",
			[][]int{{1, 4}, {4, 5}},
			[][]int{{1, 5}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := merge(tc.intervals)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_insert(t *testing.T) {
	testCases := []struct {
		name        string
		intervals   [][]int
		newInterval []int
		expected    [][]int
	}{
		{
			"case 1",
			[][]int{{1, 3}, {6, 9}},
			[]int{2, 5},
			[][]int{{1, 5}, {6, 9}},
		},
		{
			"case 2",
			[][]int{{1, 2}, {3, 5}, {6, 7}, {8, 10}, {12, 16}},
			[]int{4, 8},
			[][]int{{1, 2}, {3, 10}, {12, 16}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := insert(tc.intervals, tc.newInterval)
			assert.Equal(t, tc.expected, result)
			result = insertWithSlice(tc.intervals, tc.newInterval)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_eraseOverlapIntervals(t *testing.T) {
	tests := []struct {
		name      string
		intervals [][]int
		expected  int
	}{
		{"leetcode_1", [][]int{{1, 2}, {2, 3}, {3, 4}, {1, 3}}, 1},
		{"leetcode_2", [][]int{{1, 2}, {1, 2}, {1, 2}}, 2},
		{"leetcode_3", [][]int{{1, 2}, {2, 3}}, 0},
		{"single", [][]int{{1, 5}}, 0},
		{"no_overlap", [][]int{{1, 2}, {3, 4}, {5, 6}}, 0},
		{"all_overlap", [][]int{{1, 10}, {2, 5}, {3, 4}}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in1 := deepCopyMatrix(tt.intervals)
			in2 := deepCopyMatrix(tt.intervals)
			assert.Equal(t, tt.expected, eraseOverlapIntervals(in1), "sort_by_start")
			assert.Equal(t, tt.expected, eraseOverlapIntervalsSortByEndTime(in2), "sort_by_end")
		})
	}
}

func Test_canAttendMeetings(t *testing.T) {
	tests := []struct {
		name      string
		intervals []Interval
		expected  bool
	}{
		{"can_attend", []Interval{{Start: 0, End: 30}, {Start: 35, End: 50}, {Start: 55, End: 60}}, true},
		{"overlap", []Interval{{Start: 0, End: 30}, {Start: 5, End: 10}}, false},
		{"adjacent", []Interval{{Start: 0, End: 10}, {Start: 10, End: 20}}, true},
		{"single", []Interval{{Start: 5, End: 10}}, true},
		{"exact_overlap", []Interval{{Start: 1, End: 5}, {Start: 1, End: 5}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intervals := append([]Interval(nil), tt.intervals...)
			assert.Equal(t, tt.expected, canAttendMeetings(intervals))
		})
	}
}

func Test_minMeetingRooms(t *testing.T) {
	tests := []struct {
		name      string
		intervals []Interval
		expected  int
	}{
		{"leetcode_1", []Interval{{Start: 0, End: 30}, {Start: 5, End: 10}, {Start: 15, End: 20}}, 2},
		{"no_overlap", []Interval{{Start: 0, End: 10}, {Start: 10, End: 20}}, 1},
		{"all_overlap", []Interval{{Start: 0, End: 10}, {Start: 1, End: 9}, {Start: 2, End: 8}}, 3},
		{"single", []Interval{{Start: 5, End: 10}}, 1},
		{"two_rooms", []Interval{{Start: 0, End: 5}, {Start: 0, End: 5}}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in1 := append([]Interval(nil), tt.intervals...)
			in2 := append([]Interval(nil), tt.intervals...)
			assert.Equal(t, tt.expected, minMeetingRooms(in1), "two_pointer")
			assert.Equal(t, tt.expected, minMeetingRoomsSweepLine(in2), "sweep_line")
		})
	}
}

func Test_minMeetingRoomsMinHeap(t *testing.T) {
	tests := []struct {
		name      string
		intervals [][]int
		expected  int
	}{
		{"leetcode_1", [][]int{{0, 30}, {5, 10}, {15, 20}}, 2},
		{"no_overlap", [][]int{{0, 10}, {10, 20}}, 1},
		{"all_overlap", [][]int{{0, 10}, {1, 9}, {2, 8}}, 3},
		{"single", [][]int{{5, 10}}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := deepCopyMatrix(tt.intervals)
			assert.Equal(t, tt.expected, minMeetingRoomsMinHeap(in))
		})
	}
}

func Test_employeeFreeTime(t *testing.T) {
	tests := []struct {
		name     string
		schedule [][][]int
		expected [][]int
	}{
		{
			"leetcode_1",
			[][][]int{{{1, 3}, {6, 7}}, {{2, 4}}, {{2, 5}, {9, 12}}},
			[][]int{{5, 6}, {7, 9}},
		},
		{
			"single_gap",
			[][][]int{{{1, 3}}, {{5, 7}}},
			[][]int{{3, 5}},
		},
		{
			"no_gap",
			[][][]int{{{1, 5}}, {{2, 4}}},
			[][]int{},
		},
		{
			"three_gaps",
			[][][]int{{{1, 2}}, {{4, 5}}, {{7, 8}}},
			[][]int{{2, 4}, {5, 7}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := employeeFreeTime(tt.schedule)
			if len(tt.expected) == 0 {
				assert.Empty(t, got)
			} else {
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func Test_mergeIntervals(t *testing.T) {
	tests := []struct {
		name      string
		intervals [][]int
		expected  [][]int
	}{
		{
			name:      "leetcode_example1",
			intervals: [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}},
			expected:  [][]int{{1, 6}, {8, 10}, {15, 18}},
		},
		{
			name:      "leetcode_example2",
			intervals: [][]int{{1, 4}, {4, 5}},
			expected:  [][]int{{1, 5}},
		},
		{
			name:      "no_overlap",
			intervals: [][]int{{1, 2}, {3, 4}, {5, 6}},
			expected:  [][]int{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			name:      "all_merge",
			intervals: [][]int{{1, 10}, {2, 3}, {4, 8}},
			expected:  [][]int{{1, 10}},
		},
		{
			name:      "single_interval",
			intervals: [][]int{{1, 5}},
			expected:  [][]int{{1, 5}},
		},
		{
			name:      "unordered_input",
			intervals: [][]int{{4, 5}, {1, 3}, {2, 4}},
			expected:  [][]int{{1, 5}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeIntervals(deepCopyMatrix(tt.intervals))
			assert.Equal(t, tt.expected, got)
		})
	}
}

func Test_intervalIntersection(t *testing.T) {
	tests := []struct {
		name       string
		firstList  [][]int
		secondList [][]int
		expected   [][]int
	}{
		{
			name:       "leetcode_example",
			firstList:  [][]int{{0, 2}, {5, 10}, {13, 23}, {24, 25}},
			secondList: [][]int{{1, 5}, {8, 12}, {15, 24}, {25, 26}},
			expected:   [][]int{{1, 2}, {5, 5}, {8, 10}, {15, 23}, {24, 24}, {25, 25}},
		},
		{
			name:       "no_intersection",
			firstList:  [][]int{{1, 2}, {5, 6}},
			secondList: [][]int{{3, 4}, {7, 8}},
			expected:   [][]int{},
		},
		{
			name:       "one_list_empty",
			firstList:  [][]int{{1, 5}},
			secondList: [][]int{},
			expected:   [][]int{},
		},
		{
			name:       "full_overlap",
			firstList:  [][]int{{1, 10}},
			secondList: [][]int{{2, 4}, {6, 8}},
			expected:   [][]int{{2, 4}, {6, 8}},
		},
		{
			name:       "exact_touch",
			firstList:  [][]int{{1, 3}},
			secondList: [][]int{{3, 5}},
			expected:   [][]int{{3, 3}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := intersect(tt.firstList, tt.secondList)
			assert.Equal(t, tt.expected, got)

			got2 := intervalIntersection(tt.firstList, tt.secondList)
			assert.Equal(t, tt.expected, got2)
		})
	}
}
