package apidesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type lruOp struct {
	action   string // "get" or "put"
	key, val int
	want     int // only checked for "get"
}

func Test_LRUCache(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
		ops      []lruOp
	}{
		{
			name:     "miss on empty cache returns -1",
			capacity: 2,
			ops:      []lruOp{{"get", 1, 0, -1}},
		},
		{
			name:     "hit returns correct value",
			capacity: 2,
			ops: []lruOp{
				{"put", 1, 10, 0},
				{"get", 1, 0, 10},
			},
		},
		{
			name:     "miss on absent key returns -1",
			capacity: 2,
			ops: []lruOp{
				{"put", 1, 10, 0},
				{"get", 99, 0, -1},
			},
		},
		{
			name:     "update existing key changes value without eviction",
			capacity: 2,
			ops: []lruOp{
				{"put", 1, 10, 0},
				{"put", 2, 20, 0},
				{"put", 1, 100, 0},
				{"get", 1, 0, 100},
				{"get", 2, 0, 20},
			},
		},
		{
			name:     "update promotes key to most recently used",
			capacity: 2,
			ops: []lruOp{
				{"put", 1, 1, 0},
				{"put", 2, 2, 0},
				{"put", 1, 11, 0}, // key 1 becomes MRU; key 2 is now LRU
				{"put", 3, 3, 0},  // evicts key 2
				{"get", 1, 0, 11},
				{"get", 2, 0, -1},
				{"get", 3, 0, 3},
			},
		},
		{
			name:     "evicts least recently used on overflow",
			capacity: 3,
			ops: []lruOp{
				{"put", 1, 1, 0}, // MRU→LRU: [1]
				{"put", 2, 2, 0}, // [2, 1]
				{"get", 1, 0, 1}, // [1, 2]
				{"put", 3, 3, 0}, // [3, 1, 2]
				{"get", 2, 0, 2}, // [2, 3, 1]
				{"put", 4, 4, 0}, // evicts key 1 (LRU); [4, 2, 3]
				{"get", 1, 0, -1},
				{"get", 2, 0, 2},
				{"get", 3, 0, 3},
				{"get", 4, 0, 4},
			},
		},
		{
			name:     "capacity 1 evicts previous key on each put",
			capacity: 1,
			ops: []lruOp{
				{"put", 1, 1, 0},
				{"put", 2, 2, 0},
				{"get", 1, 0, -1},
				{"get", 2, 0, 2},
			},
		},
		{
			name:     "get refreshes key preventing its eviction",
			capacity: 2,
			ops: []lruOp{
				{"put", 1, 1, 0},
				{"put", 2, 2, 0},
				{"get", 1, 0, 1},  // key 1 becomes MRU; key 2 is LRU
				{"put", 3, 3, 0},  // evicts key 2, not key 1
				{"get", 1, 0, 1},
				{"get", 2, 0, -1},
				{"get", 3, 0, 3},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cache := ConstructorLRUCache(tc.capacity)
			for _, op := range tc.ops {
				if op.action == "put" {
					cache.Put(op.key, op.val)
				} else {
					assert.Equal(t, op.want, cache.Get(op.key))
				}
			}
		})
	}
}
