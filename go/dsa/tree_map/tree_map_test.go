package treemap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTreeMap() *TreeMap {
	return &TreeMap{m: make(map[int]int)}
}

func Test_TreeMap_Put(t *testing.T) {
	testCases := []struct {
		name         string
		puts         [][2]int // {key, value}
		expectedKeys []int
	}{
		{
			name:         "single put",
			puts:         [][2]int{{5, 10}},
			expectedKeys: []int{5},
		},
		{
			name:         "keys inserted in order stay sorted",
			puts:         [][2]int{{1, 1}, {3, 3}, {5, 5}},
			expectedKeys: []int{1, 3, 5},
		},
		{
			name:         "keys inserted out of order become sorted",
			puts:         [][2]int{{5, 5}, {1, 1}, {3, 3}},
			expectedKeys: []int{1, 3, 5},
		},
		{
			name:         "duplicate key updates value not keys slice",
			puts:         [][2]int{{2, 10}, {4, 20}, {2, 99}},
			expectedKeys: []int{2, 4},
		},
		{
			name:         "negative keys sorted correctly",
			puts:         [][2]int{{3, 0}, {-1, 0}, {0, 0}},
			expectedKeys: []int{-1, 0, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tm := newTreeMap()
			for _, p := range tc.puts {
				tm.Put(p[0], p[1])
			}
			assert.Equal(t, tc.expectedKeys, tm.keys)
		})
	}

	t.Run("duplicate key value is updated", func(t *testing.T) {
		tm := newTreeMap()
		tm.Put(5, 10)
		tm.Put(5, 99)
		assert.Equal(t, 99, tm.m[5])
		assert.Equal(t, []int{5}, tm.keys)
	})
}

func Test_TreeMap_Floor(t *testing.T) {
	testCases := []struct {
		name          string
		puts          [][2]int
		query         int
		expectedKey   int
		expectedFound bool
	}{
		{
			name:          "exact match returns that key",
			puts:          [][2]int{{1, 0}, {3, 0}, {5, 0}},
			query:         3,
			expectedKey:   3,
			expectedFound: true,
		},
		{
			name:          "between keys returns smaller",
			puts:          [][2]int{{1, 0}, {3, 0}, {5, 0}},
			query:         4,
			expectedKey:   3,
			expectedFound: true,
		},
		{
			name:          "above all keys returns largest",
			puts:          [][2]int{{1, 0}, {3, 0}, {5, 0}},
			query:         10,
			expectedKey:   5,
			expectedFound: true,
		},
		{
			name:          "below all keys returns not found",
			puts:          [][2]int{{1, 0}, {3, 0}, {5, 0}},
			query:         0,
			expectedKey:   0,
			expectedFound: false,
		},
		{
			name:          "empty map returns not found",
			puts:          [][2]int{},
			query:         5,
			expectedKey:   0,
			expectedFound: false,
		},
		{
			name:          "query equals minimum key",
			puts:          [][2]int{{2, 0}, {4, 0}, {6, 0}},
			query:         2,
			expectedKey:   2,
			expectedFound: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tm := newTreeMap()
			for _, p := range tc.puts {
				tm.Put(p[0], p[1])
			}
			key, found := tm.Floor(tc.query)
			assert.Equal(t, tc.expectedFound, found)
			if tc.expectedFound {
				assert.Equal(t, tc.expectedKey, key)
			}
		})
	}
}

func Test_TreeMap_Ceil(t *testing.T) {
	testCases := []struct {
		name          string
		puts          [][2]int
		query         int
		expectedKey   int
		expectedFound bool
	}{
		{
			name:          "exact match returns that key",
			puts:          [][2]int{{1, 0}, {3, 0}, {5, 0}},
			query:         3,
			expectedKey:   3,
			expectedFound: true,
		},
		{
			name:          "between keys returns larger",
			puts:          [][2]int{{1, 0}, {3, 0}, {5, 0}},
			query:         4,
			expectedKey:   5,
			expectedFound: true,
		},
		{
			name:          "below all keys returns smallest",
			puts:          [][2]int{{1, 0}, {3, 0}, {5, 0}},
			query:         0,
			expectedKey:   1,
			expectedFound: true,
		},
		{
			name:          "above all keys returns not found",
			puts:          [][2]int{{1, 0}, {3, 0}, {5, 0}},
			query:         6,
			expectedKey:   0,
			expectedFound: false,
		},
		{
			name:          "empty map returns not found",
			puts:          [][2]int{},
			query:         5,
			expectedKey:   0,
			expectedFound: false,
		},
		{
			name:          "query equals maximum key",
			puts:          [][2]int{{2, 0}, {4, 0}, {6, 0}},
			query:         6,
			expectedKey:   6,
			expectedFound: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tm := newTreeMap()
			for _, p := range tc.puts {
				tm.Put(p[0], p[1])
			}
			key, found := tm.Ceil(tc.query)
			assert.Equal(t, tc.expectedFound, found)
			if tc.expectedFound {
				assert.Equal(t, tc.expectedKey, key)
			}
		})
	}
}

func Test_longestSubarray(t *testing.T) {
	testCases := []struct {
		name     string
		nums     []int
		limit    int
		expected int
	}{
		{
			name:     "leetcode example 1",
			nums:     []int{8, 2, 4, 7},
			limit:    4,
			expected: 2,
		},
		{
			name:     "leetcode example 2",
			nums:     []int{10, 1, 2, 4, 7, 2},
			limit:    5,
			expected: 4,
		},
		{
			name:     "leetcode example 3",
			nums:     []int{4, 2, 2, 2, 4, 4, 2, 2},
			limit:    0,
			expected: 3,
		},
		{
			name:     "single element always valid",
			nums:     []int{5},
			limit:    0,
			expected: 1,
		},
		{
			name:     "all same values",
			nums:     []int{3, 3, 3, 3},
			limit:    0,
			expected: 4,
		},
		{
			name:     "strictly increasing within limit",
			nums:     []int{1, 2, 3, 4, 5},
			limit:    4,
			expected: 5,
		},
		{
			name:     "limit zero only consecutive equal",
			nums:     []int{1, 2, 1},
			limit:    0,
			expected: 1,
		},
		{
			name:     "large limit includes whole array",
			nums:     []int{1, 100, 50, 75},
			limit:    99,
			expected: 4,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := longestSubarray(tc.nums, tc.limit)
			assert.Equal(t, tc.expected, result)
		})
	}
}
