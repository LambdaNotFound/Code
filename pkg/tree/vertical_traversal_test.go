package tree

import (
	. "gocode/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerticalTraversal(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected [][]int
	}{
		{
			name:     "nil root",
			root:     nil,
			expected: nil,
		},
		{
			name:     "single node",
			root:     &TreeNode{Val: 1},
			expected: [][]int{{1}},
		},
		{
			// [3,9,20,null,null,15,7]
			//     3
			//    / \
			//   9  20
			//     /  \
			//    15   7
			name:     "example 1",
			root:     node(3, &TreeNode{Val: 9}, node(20, &TreeNode{Val: 15}, &TreeNode{Val: 7})),
			expected: [][]int{{9}, {3, 15}, {20}, {7}},
		},
		{
			// [1,2,3,4,5,6,7]
			//         1
			//        / \
			//       2   3
			//      / \ / \
			//     4  5 6  7
			// 5 and 6 share (col=0, row=2) → tie broken by value ascending
			name: "example 2 - tie broken by value",
			root: node(1,
				node(2, &TreeNode{Val: 4}, &TreeNode{Val: 5}),
				node(3, &TreeNode{Val: 6}, &TreeNode{Val: 7}),
			),
			expected: [][]int{{4}, {2}, {1, 5, 6}, {3}, {7}},
		},
		{
			// left-skewed with right subturn
			//   1
			//  /
			// 2
			//  \
			//   3   ← col=-1+1=0, row=2 → same col as root
			name:     "left child with right grandchild lands in root column",
			root:     node(1, node(2, nil, &TreeNode{Val: 3}), nil),
			expected: [][]int{{2}, {1, 3}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, verticalTraversal(tt.root))
		})
	}
}
