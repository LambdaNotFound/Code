package graph

import (
	. "gocode/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func deepCopy2D(src [][]int) [][]int {
    dst := make([][]int, len(src))
    for i := range src {
        dst[i] = make([]int, len(src[i]))
        copy(dst[i], src[i])
    }
    return dst
}

func Test_orangesRotting(t *testing.T) {
    testCases := []struct {
        name     string
        grid     [][]int
        expected int
    }{
        {
            "case 1",
            [][]int{{2, 1, 1}, {1, 1, 0}, {0, 1, 1}},
            4,
        },
        {
            "case 2",
            [][]int{{2, 1, 1}, {0, 1, 1}, {1, 0, 1}},
            -1,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            grid := deepCopy2D(tc.grid)
            result := orangesRotting(grid)
            assert.Equal(t, tc.expected, result)

            grid = deepCopy2D(tc.grid)
            result = orangesRotting_slice(grid)
            assert.Equal(t, tc.expected, result)
        })
    }
}

func Test_updateMatrix(t *testing.T) {
    tests := []struct {
        name     string
        input    [][]int
        expected [][]int
    }{
        {
            name: "Simple 3x3 matrix",
            input: [][]int{
                {0, 0, 0},
                {0, 1, 0},
                {1, 1, 1},
            },
            expected: [][]int{
                {0, 0, 0},
                {0, 1, 0},
                {1, 2, 1},
            },
        },
        {
            name: "All zeros",
            input: [][]int{
                {0, 0},
                {0, 0},
            },
            expected: [][]int{
                {0, 0},
                {0, 0},
            },
        },
        {
            name: "All ones, single zero in corner",
            input: [][]int{
                {0, 1, 1},
                {1, 1, 1},
                {1, 1, 1},
            },
            expected: [][]int{
                {0, 1, 2},
                {1, 2, 3},
                {2, 3, 4},
            },
        },
        {
            name: "Single cell zero",
            input: [][]int{
                {0},
            },
            expected: [][]int{
                {0},
            },
        },
        {
            name: "Rectangle 2x4",
            input: [][]int{
                {0, 0, 1, 1},
                {1, 1, 1, 0},
            },
            expected: [][]int{
                {0, 0, 1, 1},
                {1, 1, 1, 0},
            },
        },
        {
            name: "Zigzag zeros",
            input: [][]int{
                {0, 1, 0},
                {1, 1, 1},
                {0, 1, 0},
            },
            expected: [][]int{
                {0, 1, 0},
                {1, 2, 1},
                {0, 1, 0},
            },
        },
    }

    for _, tc := range tests {
        assert.Equal(t, tc.expected, updateMatrix(tc.input), "failed test: %s", tc.name)
    }
}

// Helper: Build a binary tree from level-order array (nil = missing node)
func buildTree(vals []any) *TreeNode {
    if len(vals) == 0 || vals[0] == nil {
        return nil
    }

    nodes := make([]*TreeNode, len(vals))
    for i, v := range vals {
        if v != nil {
            nodes[i] = &TreeNode{Val: v.(int)}
        }
    }

    for i := 0; i < len(vals); i++ {
        if nodes[i] == nil {
            continue
        }
        leftIdx := 2*i + 1
        rightIdx := 2*i + 2
        if leftIdx < len(vals) {
            nodes[i].Left = nodes[leftIdx]
        }
        if rightIdx < len(vals) {
            nodes[i].Right = nodes[rightIdx]
        }
    }
    return nodes[0]
}

func Test_rightSideView(t *testing.T) {
    tests := []struct {
        name     string
        input    []any
        expected []int
    }{
        {
            name:     "Example tree",
            input:    []any{1, 2, 3, nil, 5, nil, 4},
            expected: []int{1, 3, 4},
        },
        {
            name:     "Single node",
            input:    []any{1},
            expected: []int{1},
        },
        {
            name:     "Left-skewed tree",
            input:    []any{1, 2, nil, 3, nil, nil, nil, 4},
            expected: []int{1, 2, 3, 4},
        },
        {
            name:     "Right-skewed tree",
            input:    []any{1, nil, 2, nil, nil, nil, 3},
            expected: []int{1, 2, 3},
        },
        {
            name:     "Complete binary tree",
            input:    []any{1, 2, 3, 4, 5, 6, 7},
            expected: []int{1, 3, 7},
        },
        {
            name:     "Sparse tree",
            input:    []any{1, 2, 3, nil, 5, nil, 4},
            expected: []int{1, 3, 4},
        },
        {
            name:     "Empty tree",
            input:    []any{},
            expected: []int{},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            root := buildTree(tc.input)
            got := rightSideView(root)
            assert.Equal(t, tc.expected, got)
        })
    }
}
