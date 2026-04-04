package tree

import (
	. "gocode/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_preorderTraversal(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected []int
	}{
		{"nil", nil, []int{}},
		{"single", &TreeNode{Val: 1}, []int{1}},
		{
			"three_nodes",
			node(1, &TreeNode{Val: 2}, &TreeNode{Val: 3}),
			[]int{1, 2, 3},
		},
		{
			"full_tree",
			node(1,
				node(2, &TreeNode{Val: 4}, &TreeNode{Val: 5}),
				node(3, &TreeNode{Val: 6}, &TreeNode{Val: 7}),
			),
			[]int{1, 2, 4, 5, 3, 6, 7},
		},
		{
			"skewed_right",
			&TreeNode{Val: 1, Right: &TreeNode{Val: 2, Right: &TreeNode{Val: 3}}},
			[]int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, preorderTraversal(tt.root))
		})
	}
}

func Test_inorderTraversal(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected []int
	}{
		{"nil", nil, []int{}},
		{"single", &TreeNode{Val: 1}, []int{1}},
		{
			"three_nodes",
			node(2, &TreeNode{Val: 1}, &TreeNode{Val: 3}),
			[]int{1, 2, 3},
		},
		{
			"bst_order",
			node(4, node(2, &TreeNode{Val: 1}, &TreeNode{Val: 3}), node(6, &TreeNode{Val: 5}, &TreeNode{Val: 7})),
			[]int{1, 2, 3, 4, 5, 6, 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, inorderTraversal(tt.root))
		})
	}
}

func Test_buildTree(t *testing.T) {
	tests := []struct {
		name     string
		preorder []int
		inorder  []int
		// verify via inorder traversal of result
		expectedInorder []int
	}{
		{"single", []int{1}, []int{1}, []int{1}},
		{
			"three_nodes",
			[]int{3, 9, 20},
			[]int{9, 3, 20},
			[]int{9, 3, 20},
		},
		{
			"full_tree",
			[]int{1, 2, 4, 5, 3, 6, 7},
			[]int{4, 2, 5, 1, 6, 3, 7},
			[]int{4, 2, 5, 1, 6, 3, 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := buildTree(tt.preorder, tt.inorder)
			assert.Equal(t, tt.expectedInorder, inorderTraversal(root))
		})
	}
}

func Test_invertTree(t *testing.T) {
	tests := []struct {
		name            string
		root            *TreeNode
		expectedInorder []int
	}{
		{"nil", nil, []int{}},
		{"single", &TreeNode{Val: 1}, []int{1}},
		{
			"three_nodes",
			node(2, &TreeNode{Val: 1}, &TreeNode{Val: 3}),
			[]int{3, 2, 1}, // inverted inorder is reversed
		},
		{
			"full_tree",
			node(4, node(2, &TreeNode{Val: 1}, &TreeNode{Val: 3}), node(7, &TreeNode{Val: 6}, &TreeNode{Val: 9})),
			[]int{9, 7, 6, 4, 3, 2, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedInorder, inorderTraversal(invertTree(tt.root)))
		})
	}
}

func Test_widthOfBinaryTree(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected int
	}{
		{
			"leetcode_1",
			node(1, node(3, &TreeNode{Val: 5}, &TreeNode{Val: 3}), node(2, nil, &TreeNode{Val: 9})),
			4,
		},
		{
			"leetcode_2",
			node(1, node(3, &TreeNode{Val: 5}, nil), &TreeNode{Val: 2}),
			2,
		},
		{"single", &TreeNode{Val: 1}, 1},
		{
			"nil",
			nil, 0,
		},
		{
			"two_children",
			node(1, &TreeNode{Val: 2}, &TreeNode{Val: 3}),
			2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, widthOfBinaryTree(tt.root))
		})
	}
}

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
