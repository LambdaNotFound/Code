package tree

import . "gocode/types"

/**
 * 285. Inorder Successor in BST
 *
 * BST inorder traversal, sorted
 *
 *        p.val,  _
 *              greater
 */
func inorderSuccessor(root, p *TreeNode) *TreeNode {
	res := root
	for root != nil {
		if p.Val <= root.Val {
			res = root
			root = root.Left
		} else {
			root = root.Right
		}
	}

	return res
}
