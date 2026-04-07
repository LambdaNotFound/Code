package utils

import (
	. "gocode/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

// inorder returns the in-order traversal of a BST as a slice of values.
func inorder(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	result := inorder(root.Left)
	result = append(result, root.Val)
	result = append(result, inorder(root.Right)...)
	return result
}

func Test_GraphsEqual(t *testing.T) {
	tests := []struct {
		name     string
		a, b     *Node
		expected bool
	}{
		{"both_nil", nil, nil, true},
		{"a_nil_b_not", nil, &Node{Val: 1}, false},
		{"b_nil_a_not", &Node{Val: 1}, nil, false},
		{"single_node_equal", &Node{Val: 5}, &Node{Val: 5}, true},
		{"single_node_different_val", &Node{Val: 1}, &Node{Val: 2}, false},
		{
			name:     "neighbor_count_mismatch",
			a:        &Node{Val: 1, Neighbors: []*Node{{Val: 2}}},
			b:        &Node{Val: 1},
			expected: false,
		},
		{
			name:     "two_nodes_equal",
			a:        &Node{Val: 1, Neighbors: []*Node{{Val: 2}}},
			b:        &Node{Val: 1, Neighbors: []*Node{{Val: 2}}},
			expected: true,
		},
		{
			name:     "two_nodes_neighbor_val_mismatch",
			a:        &Node{Val: 1, Neighbors: []*Node{{Val: 2}}},
			b:        &Node{Val: 1, Neighbors: []*Node{{Val: 99}}},
			expected: false,
		},
		{
			name: "cycle_two_nodes_equal",
			a: func() *Node {
				n1 := &Node{Val: 1}
				n2 := &Node{Val: 2}
				n1.Neighbors = []*Node{n2}
				n2.Neighbors = []*Node{n1}
				return n1
			}(),
			b: func() *Node {
				n1 := &Node{Val: 1}
				n2 := &Node{Val: 2}
				n1.Neighbors = []*Node{n2}
				n2.Neighbors = []*Node{n1}
				return n1
			}(),
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// fresh visited map per subtest to avoid state bleed between cases
			assert.Equal(t, tc.expected, GraphsEqual(tc.a, tc.b, map[*Node]*Node{}))
		})
	}
}

func Test_BuildBST(t *testing.T) {
	tests := []struct {
		name        string
		nums        []int
		wantNil     bool
		wantRoot    int
		wantInorder []int
	}{
		{"empty_input", []int{}, true, 0, nil},
		{"single_element", []int{5}, false, 5, []int{5}},
		{"ascending_right_spine", []int{1, 2, 3, 4}, false, 1, []int{1, 2, 3, 4}},
		{"descending_left_spine", []int{4, 3, 2, 1}, false, 4, []int{1, 2, 3, 4}},
		{"balanced", []int{6, 2, 8, 0, 4, 7, 9, 3, 5}, false, 6, []int{0, 2, 3, 4, 5, 6, 7, 8, 9}},
		{"duplicate_goes_right", []int{5, 5}, false, 5, []int{5, 5}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := BuildBST(tc.nums)
			if tc.wantNil {
				assert.Nil(t, root)
				return
			}
			assert.NotNil(t, root)
			assert.Equal(t, tc.wantRoot, root.Val)
			assert.Equal(t, tc.wantInorder, inorder(root))
		})
	}
}
