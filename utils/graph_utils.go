package utils

import . "gocode/types"

// GraphsEqual checks structural equality of two graphs
func GraphsEqual(a, b *Node, visited map[*Node]*Node) bool {
    if a == nil || b == nil {
        return a == b
    }
    if a.Val != b.Val {
        return false
    }
    if v, ok := visited[a]; ok {
        return v == b
    }

    visited[a] = b
    if len(a.Neighbors) != len(b.Neighbors) {
        return false
    }
    for i := range a.Neighbors {
        if !GraphsEqual(a.Neighbors[i], b.Neighbors[i], visited) {
            return false
        }
    }

    return true
}

// insert inserts a value into the BST
func insert(root *TreeNode, val int) *TreeNode {
	if root == nil {
		return &TreeNode{Val: val}
	}
	if val < root.Val {
		root.Left = insert(root.Left, val)
	} else {
		root.Right = insert(root.Right, val)
	}
	return root
}

// buildBST builds a BST from an array of integers
func BuildBST(nums []int) *TreeNode {
	var root *TreeNode
	for _, v := range nums {
		root = insert(root, v)
	}
	return root
}