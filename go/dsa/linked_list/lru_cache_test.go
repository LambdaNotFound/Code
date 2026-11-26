package linked_list

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LRUCache(t *testing.T) {
	testCases := []struct {
		name     string
		capacity int
		ops      []struct {
			op  string // "get" or "put"
			key int
			val int // used only for "put"
		}
		expected []int // expected return values for "get" ops (-1 for miss); -1 sentinel for "put"
	}{
		{
			name:     "basic_get_miss",
			capacity: 2,
			ops: []struct {
				op  string
				key int
				val int
			}{
				{"get", 1, 0},
			},
			expected: []int{-1},
		},
		{
			name:     "put_then_get",
			capacity: 2,
			ops: []struct {
				op  string
				key int
				val int
			}{
				{"put", 1, 10},
				{"get", 1, 0},
			},
			expected: []int{-1, 10},
		},
		{
			name:     "evict_lru_on_capacity_exceeded",
			capacity: 2,
			ops: []struct {
				op  string
				key int
				val int
			}{
				{"put", 1, 1},
				{"put", 2, 2},
				{"put", 3, 3}, // evicts key=1 (LRU)
				{"get", 1, 0}, // miss
				{"get", 2, 0}, // hit
				{"get", 3, 0}, // hit
			},
			expected: []int{-1, -1, -1, -1, 2, 3},
		},
		{
			name:     "get_refreshes_recency",
			capacity: 2,
			ops: []struct {
				op  string
				key int
				val int
			}{
				{"put", 1, 1},
				{"put", 2, 2},
				{"get", 1, 0}, // refreshes key=1, so key=2 becomes LRU
				{"put", 3, 3}, // evicts key=2
				{"get", 1, 0}, // hit
				{"get", 2, 0}, // miss
				{"get", 3, 0}, // hit
			},
			expected: []int{-1, -1, 1, -1, 1, -1, 3},
		},
		{
			name:     "update_existing_key",
			capacity: 2,
			ops: []struct {
				op  string
				key int
				val int
			}{
				{"put", 1, 1},
				{"put", 1, 100}, // update value
				{"get", 1, 0},
			},
			expected: []int{-1, -1, 100},
		},
		{
			name:     "capacity_one",
			capacity: 1,
			ops: []struct {
				op  string
				key int
				val int
			}{
				{"put", 1, 1},
				{"put", 2, 2}, // evicts key=1
				{"get", 1, 0}, // miss
				{"get", 2, 0}, // hit
			},
			expected: []int{-1, -1, -1, 2},
		},
		{
			name:     "leetcode_example",
			capacity: 2,
			ops: []struct {
				op  string
				key int
				val int
			}{
				{"put", 1, 1},
				{"put", 2, 2},
				{"get", 1, 0}, // returns 1
				{"put", 3, 3}, // evicts key=2
				{"get", 2, 0}, // returns -1 (miss)
				{"put", 4, 4}, // evicts key=1
				{"get", 1, 0}, // returns -1 (miss)
				{"get", 3, 0}, // returns 3
				{"get", 4, 0}, // returns 4
			},
			expected: []int{-1, -1, 1, -1, -1, -1, -1, 3, 4},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache := Constructor(tc.capacity)
			ei := 0
			for _, op := range tc.ops {
				if op.op == "put" {
					cache.Put(op.key, op.val)
					assert.Equal(t, -1, tc.expected[ei], "put should have sentinel -1 in expected")
				} else {
					assert.Equal(t, tc.expected[ei], cache.Get(op.key))
				}
				ei++
			}
		})
	}
}
