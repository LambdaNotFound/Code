package graph

import (
	. "gocode/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_diameterOfBinaryTree(t *testing.T) {
    tests := []struct {
        name     string
        root     *TreeNode
        expected int
    }{
        {
            name:     "empty tree",
            root:     nil,
            expected: 0,
        },
        {
            name:     "single node",
            root:     &TreeNode{Val: 1},
            expected: 0,
        },
        {
            name: "simple tree",
            //    1
            //   / \
            //  2   3
            // / \
            //4   5
            root: &TreeNode{
                Val: 1,
                Left: &TreeNode{
                    Val:   2,
                    Left:  &TreeNode{Val: 4},
                    Right: &TreeNode{Val: 5},
                },
                Right: &TreeNode{Val: 3},
            },
            expected: 3, // path [4 → 2 → 1 → 3]
        },
        {
            name: "skewed tree",
            // 1
            //  \
            //   2
            //    \
            //     3
            //      \
            //       4
            root: &TreeNode{
                Val: 1,
                Right: &TreeNode{
                    Val: 2,
                    Right: &TreeNode{
                        Val:   3,
                        Right: &TreeNode{Val: 4},
                    },
                },
            },
            expected: 3, // edges along the chain
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got := diameterOfBinaryTree(tc.root)
            assert.Equal(t, tc.expected, got)
        })
    }
}

func TestNumIslands(t *testing.T) {
    tests := []struct {
        name     string
        grid     [][]byte
        expected int
    }{
        {
            name:     "empty grid",
            grid:     [][]byte{},
            expected: 0,
        },
        {
            name: "single island",
            grid: [][]byte{
                {'1', '1', '0', '0', '0'},
                {'1', '1', '0', '0', '0'},
                {'0', '0', '0', '1', '1'},
                {'0', '0', '0', '1', '1'},
            },
            expected: 2,
        },
        {
            name: "all water",
            grid: [][]byte{
                {'0', '0'},
                {'0', '0'},
            },
            expected: 0,
        },
        {
            name: "all land",
            grid: [][]byte{
                {'1', '1'},
                {'1', '1'},
            },
            expected: 1,
        },
        {
            name: "diagonal not connected",
            grid: [][]byte{
                {'1', '0'},
                {'0', '1'},
            },
            expected: 2, // diagonals don’t connect
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got := numIslands(tc.grid)
            assert.Equal(t, tc.expected, got)
        })
    }
}

func TestMaxDepth(t *testing.T) {
    tests := []struct {
        name     string
        root     *TreeNode
        expected int
    }{
        {
            name:     "empty tree",
            root:     nil,
            expected: 0,
        },
        {
            name: "single node",
            root: &TreeNode{Val: 1},
            expected: 1,
        },
        {
            name: "balanced tree",
            //    1
            //   / \
            //  2   3
            // / \
            //4   5
            root: &TreeNode{
                Val: 1,
                Left: &TreeNode{
                    Val: 2,
                    Left:  &TreeNode{Val: 4},
                    Right: &TreeNode{Val: 5},
                },
                Right: &TreeNode{Val: 3},
            },
            expected: 3,
        },
        {
            name: "skewed tree",
            // 1
            //  \
            //   2
            //    \
            //     3
            //      \
            //       4
            root: &TreeNode{
                Val: 1,
                Right: &TreeNode{
                    Val: 2,
                    Right: &TreeNode{
                        Val: 3,
                        Right: &TreeNode{Val: 4},
                    },
                },
            },
            expected: 4,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got := maxDepth(tc.root)
            assert.Equal(t, tc.expected, got)
        })
    }
}


func TestFloodFill(t *testing.T) {
    tests := []struct {
        name     string
        image    [][]int
        sr, sc   int
        newColor int
        expected [][]int
    }{
        {
            name: "basic example",
            image: [][]int{
                {1, 1, 1},
                {1, 1, 0},
                {1, 0, 1},
            },
            sr: 1, sc: 1, newColor: 2,
            expected: [][]int{
                {2, 2, 2},
                {2, 2, 0},
                {2, 0, 1},
            },
        },
        {
            name: "no change needed (already new color)",
            image: [][]int{
                {0, 0, 0},
                {0, 1, 1},
            },
            sr: 1, sc: 1, newColor: 1,
            expected: [][]int{
                {0, 0, 0},
                {0, 1, 1},
            },
        },
        {
            name: "single pixel image",
            image: [][]int{
                {0},
            },
            sr: 0, sc: 0, newColor: 2,
            expected: [][]int{
                {2},
            },
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got := floodFill(tc.image, tc.sr, tc.sc, tc.newColor)
            assert.Equal(t, tc.expected, got)
        })
    }
}