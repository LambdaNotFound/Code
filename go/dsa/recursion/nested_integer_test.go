package recursion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// helpers to build NestedInteger trees concisely
func ni(v int) *NestedInteger {
	return &NestedInteger{isInt: true, val: v}
}

func nl(items ...*NestedInteger) *NestedInteger {
	return &NestedInteger{list: items}
}

// drainIterator collects all values from the iterator.
func drainIterator(it *NestedIterator) []int {
	var res []int
	for it.HasNext() {
		res = append(res, it.Next())
	}
	return res
}

// --- 339. Nested List Weight Sum ---

func Test_depthSum(t *testing.T) {
	tests := []struct {
		name       string
		nestedList []*NestedInteger
		want       int
	}{
		{
			// [[1,1],2,[1,1]] → 1*2 + 1*2 + 2*1 + 1*2 + 1*2 = 10
			name:       "example1",
			nestedList: []*NestedInteger{nl(ni(1), ni(1)), ni(2), nl(ni(1), ni(1))},
			want:       10,
		},
		{
			// [1,[4,[6]]] → 1*1 + 4*2 + 6*3 = 27
			name:       "example2",
			nestedList: []*NestedInteger{ni(1), nl(ni(4), nl(ni(6)))},
			want:       27,
		},
		{
			// single integer at depth 1
			name:       "single integer",
			nestedList: []*NestedInteger{ni(5)},
			want:       5,
		},
		{
			// deeply nested: [[[3]]] → 3*3 = 9
			name:       "triple nested",
			nestedList: []*NestedInteger{nl(nl(ni(3)))},
			want:       9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, depthSum(tt.nestedList))
			assert.Equal(t, tt.want, depthSumDFS(tt.nestedList))
		})
	}
}

// --- 341. Flatten Nested List Iterator ---

func Test_NestedIterator(t *testing.T) {
	tests := []struct {
		name       string
		nestedList []*NestedInteger
		want       []int
	}{
		{
			name:       "example1 [[1,1],2,[1,1]]",
			nestedList: []*NestedInteger{nl(ni(1), ni(1)), ni(2), nl(ni(1), ni(1))},
			want:       []int{1, 1, 2, 1, 1},
		},
		{
			name:       "example2 [1,[4,[6]]]",
			nestedList: []*NestedInteger{ni(1), nl(ni(4), nl(ni(6)))},
			want:       []int{1, 4, 6},
		},
		{
			name:       "single integer",
			nestedList: []*NestedInteger{ni(7)},
			want:       []int{7},
		},
		{
			name:       "all nested single level",
			nestedList: []*NestedInteger{nl(ni(1), ni(2)), nl(ni(3), ni(4))},
			want:       []int{1, 2, 3, 4},
		},
		{
			name:       "deeply nested [[[1,2]]]",
			nestedList: []*NestedInteger{nl(nl(ni(1), ni(2)))},
			want:       []int{1, 2},
		},
		{
			name:       "empty outer list",
			nestedList: []*NestedInteger{},
			want:       nil,
		},
		{
			name:       "list containing empty list",
			nestedList: []*NestedInteger{nl(), ni(3)},
			want:       []int{3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := Constructor(tt.nestedList)
			assert.Equal(t, tt.want, drainIterator(it))
		})
	}
}

func Test_NestedIterator_HasNext_idempotent(t *testing.T) {
	it := Constructor([]*NestedInteger{ni(1), ni(2)})
	// calling HasNext multiple times without Next should not advance
	assert.True(t, it.HasNext())
	assert.True(t, it.HasNext())
	assert.Equal(t, 1, it.Next())
	assert.True(t, it.HasNext())
	assert.Equal(t, 2, it.Next())
	assert.False(t, it.HasNext())
}
