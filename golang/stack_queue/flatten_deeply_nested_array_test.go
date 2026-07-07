package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_flat(t *testing.T) {
	tests := []struct {
		name  string
		arr   []any
		n     int
		want  []any
	}{
		{
			name: "empty array",
			arr:  []any{},
			n:    1,
			want: []any{},
		},
		{
			name: "already flat, depth 1",
			arr:  []any{1, 2, 3},
			n:    1,
			want: []any{1, 2, 3},
		},
		{
			name: "single nesting, depth 1 flattens fully",
			arr:  []any{[]any{1, 2}, 3},
			n:    1,
			want: []any{1, 2, 3},
		},
		{
			name: "double nesting, depth 1 only unwraps one level",
			arr:  []any{[]any{1, []any{2}}, 3},
			n:    1,
			want: []any{1, []any{2}, 3},
		},
		{
			name: "double nesting, depth 2 fully flattens",
			arr:  []any{[]any{1, []any{2}}, 3},
			n:    2,
			want: []any{1, 2, 3},
		},
		{
			name: "depth 0 does not flatten",
			arr:  []any{1, []any{2, []any{3}}},
			n:    0,
			want: []any{1, []any{2, []any{3}}},
		},
		{
			name: "deeply nested, depth exactly matches nesting level",
			arr:  []any{1, []any{2, []any{3, []any{4}}}},
			n:    3,
			want: []any{1, 2, 3, 4},
		},
		{
			name: "depth larger than nesting level fully flattens",
			arr:  []any{1, []any{2, []any{3}}},
			n:    100,
			want: []any{1, 2, 3},
		},
		{
			name: "multiple nested subarrays at same level",
			arr:  []any{[]any{1, 2}, []any{3, 4}, []any{5}},
			n:    1,
			want: []any{1, 2, 3, 4, 5},
		},
		{
			name: "leetcode example 1: arr=[[1,2],3,[4,[5,6]]], n=1",
			arr:  []any{[]any{1, 2}, 3, []any{4, []any{5, 6}}},
			n:    1,
			want: []any{1, 2, 3, 4, []any{5, 6}},
		},
		{
			name: "leetcode example 2: arr=[1,[2,[3,[4,[5]]]]], n=2",
			arr:  []any{1, []any{2, []any{3, []any{4, []any{5}}}}},
			n:    2,
			want: []any{1, 2, 3, []any{4, []any{5}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flat(tt.arr, tt.n)
			assert.Equal(t, tt.want, got)
		})
	}
}
