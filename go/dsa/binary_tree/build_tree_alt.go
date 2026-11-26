package binarytree

import (
	"slices"

	. "gocode/dsa/types"
)

/**
 * 105. Construct Binary Tree from Preorder and Inorder Traversal
 *
 * 1. Recursive Approach
 *
 */
func buildTreeAlt(preorder []int, inorder []int) *TreeNode {
	if len(preorder) == 0 {
		return nil
	}

	i := slices.Index(inorder, preorder[0])
	lpre, rpre := preorder[1:i+1], preorder[i+1:]
	lin, rin := inorder[:i], inorder[i+1:]

	return &TreeNode{
		Val:   preorder[0],
		Left:  buildTreeAlt(lpre, lin),
		Right: buildTreeAlt(rpre, rin),
	}
}
