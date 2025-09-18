package tree

import (
	. "gocode/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostorderTraversal(t *testing.T) {
    tests := []struct {
        name     string
        root     *TreeNode
        expected []int
    }{
        {
            name:     "empty tree",
            root:     nil,
            expected: []int{},
        },
        {
            name:     "single node",
            root:     &TreeNode{Val: 1},
            expected: []int{1},
        },
        {
            name: "three nodes skewed right",
            root: &TreeNode{
                Val: 1,
                Right: &TreeNode{
                    Val:  2,
                    Left: &TreeNode{Val: 3},
                },
            },
            expected: []int{3, 2, 1}, // Left → Right → Root
        },
        {
            name: "full tree",
            root: &TreeNode{
                Val: 1,
                Left: &TreeNode{
                    Val:   2,
                    Left:  &TreeNode{Val: 4},
                    Right: &TreeNode{Val: 5},
                },
                Right: &TreeNode{
                    Val:   3,
                    Left:  &TreeNode{Val: 6},
                    Right: &TreeNode{Val: 7},
                },
            },
            expected: []int{4, 5, 2, 6, 7, 3, 1},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got := postorderTraversal(tc.root)

            assert.Equal(t, got, tc.expected)
        })
    }
}
